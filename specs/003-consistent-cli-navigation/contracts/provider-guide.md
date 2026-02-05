# Guide: Adding New Cloud Providers to cogo

**Last Updated**: 2026-01-18

## Overview

This guide explains how to add support for a new cloud provider while maintaining navigation consistency with existing providers.

## Provider Structure

Each provider should have its own package with a consistent structure:

```
cogo/
├── digitalocean/          # Existing provider (reference)
│   ├── digitalocean.go    # Legacy API interactions
│   ├── create_flow.go     # Create flow with navigation
│   ├── destroy_flow.go    # Destroy flow with navigation
│   └── *_test.go          # Tests
├── aws/                   # New provider example
│   ├── aws.go             # API client setup
│   ├── ec2_create_flow.go # EC2 instance creation
│   ├── ec2_destroy_flow.go# EC2 instance destruction
│   └── *_test.go          # Tests
└── azure/                 # Another provider example
    ├── azure.go
    ├── vm_create_flow.go
    └── *_test.go
```

## Step 1: Create Provider Package

**File**: `{provider}/{provider}.go`

```go
package aws

import (
    "context"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
)

// Client wraps the AWS SDK client
type Client struct {
    ec2 *ec2.EC2
}

// NewClient creates a new AWS client
func NewClient(region, accessKey, secretKey string) (*Client, error) {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(region),
        // Credentials configuration
    })
    if err != nil {
        return nil, err
    }
    
    return &Client{
        ec2: ec2.New(sess),
    }, nil
}

// Helper functions for API interactions
func (c *Client) ListInstances(ctx context.Context) ([]*ec2.Instance, error) {
    // Implementation
}

func (c *Client) CreateInstance(ctx context.Context, config *InstanceConfig) (*ec2.Instance, error) {
    // Implementation
}

func (c *Client) DeleteInstance(ctx context.Context, instanceID string) error {
    // Implementation
}
```

## Step 2: Implement Create Flow

**File**: `{provider}/create_flow.go`

**Must follow same pattern as DigitalOcean**:
1. Use navigation framework
2. Multiple steps with back navigation
3. Empty state handling
4. Summary before confirmation
5. Standard error messages

```go
package aws

import (
    "context"
    "fmt"
    
    "github.com/Joel-Valentine/cogo/navigation"
    "github.com/fatih/color"
)

// CreateInstanceFlow orchestrates EC2 instance creation
type CreateInstanceFlow struct {
    client *Client
    state  navigation.State
    steps  []navigation.Step
}

func NewCreateInstanceFlow(client *Client) *CreateInstanceFlow {
    flow := &CreateInstanceFlow{
        client: client,
        state:  navigation.NewState(),
    }
    
    // Steps should mirror DigitalOcean for consistency
    flow.steps = []navigation.Step{
        &InstanceNameStep{},
        &AMISelectionStep{client: client},
        &InstanceTypeStep{client: client},
        &RegionStep{client: client},
        &KeyPairStep{client: client},
        &SecurityGroupStep{client: client},
        &ConfirmationStep{},
    }
    
    return flow
}

func (f *CreateInstanceFlow) Name() string {
    return "Create EC2 Instance"
}

func (f *CreateInstanceFlow) Steps() []navigation.Step {
    return f.steps
}

func (f *CreateInstanceFlow) State() navigation.State {
    return f.state
}

// ExecuteCreateInstanceFlow runs the flow
func ExecuteCreateInstanceFlow(client *Client) (*ec2.Instance, error) {
    ctx := context.Background()
    navigator := navigation.NewNavigator()
    flow := NewCreateInstanceFlow(client)
    
    result, err := navigator.Run(ctx, flow)
    
    // REQUIRED: Standard error handling pattern
    if err == navigation.ErrEmptyState {
        return nil, nil
    }
    if err == navigation.ErrCancel {
        color.Cyan("\nInstance creation canceled.")
        return nil, nil
    }
    if err == context.Canceled {
        fmt.Println()
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("create flow failed: %w", err)
    }
    
    confirmed := result.Value.(bool)
    if !confirmed {
        color.Cyan("Instance creation canceled.")
        return nil, nil
    }
    
    // Extract selections and create instance
    state := flow.State()
    // ... create instance using selections
    
    return instance, nil
}
```

## Step 3: Implement Steps

Each step **MUST**:
1. Check for empty state
2. Support back navigation
3. Use standard prompts
4. Return results with metadata
5. Handle errors gracefully

**Example Step**:

```go
type AMISelectionStep struct {
    client *Client
}

func (s *AMISelectionStep) Name() string {
    return "ami_selection"  // Unique, snake_case
}

func (s *AMISelectionStep) Prompt() string {
    return "Select AMI"
}

func (s *AMISelectionStep) Execute(ctx context.Context, state navigation.State) (navigation.Result, error) {
    // 1. Fetch AMIs
    amis, err := s.client.ListAMIs(ctx)
    if err != nil {
        return navigation.Result{}, fmt.Errorf("failed to list AMIs: %w", err)
    }
    
    // 2. REQUIRED: Check empty state
    if len(amis) == 0 {
        handler := &navigation.EmptyStateHandler{
            ResourceName:    "AMIs",
            Context:         "in your AWS account",
            SuggestedAction: "Create an AMI or check your region selection.",
        }
        handler.Display()
        return navigation.Result{}, navigation.ErrEmptyState
    }
    
    // 3. Convert to string slice
    amiNames := make([]string, len(amis))
    for i, ami := range amis {
        amiNames[i] = *ami.Name
    }
    
    // 4. Create prompt with back option
    prompt := navigation.NewSelectPrompt(s.Prompt(), amiNames)
    if state.CanGoBack() {
        prompt = prompt.AddBackOption()
    }
    
    // 5. Run prompt
    index, _, err := prompt.RunWithContext(ctx)
    if err != nil {
        return navigation.Result{}, err
    }
    
    // 6. Return result with metadata
    selected := amis[index]
    return navigation.NewResultWithMetadata(*selected.ImageId, map[string]interface{}{
        "name":  *selected.Name,
        "index": index,
    }), nil
}

func (s *AMISelectionStep) Validate(result navigation.Result) error {
    return nil
}

func (s *AMISelectionStep) Default() interface{} {
    return nil
}
```

## Step 4: Wire Up Provider Selection

**File**: `utils/utils.go` (or create new provider selection)

```go
func AskForProvider() (string, error) {
    providers := []SelectItem{
        {Name: "DigitalOcean", Value: "DO"},
        {Name: "AWS", Value: "AWS"},        // Add new provider
        {Name: "Azure", Value: "AZURE"},    // Add new provider
    }
    
    selected, err := AskAndAnswerCustomSelect("Select Cloud Provider", providers)
    return selected, err
}
```

## Step 5: Add Provider Commands

**File**: `cmd/root.go`

```go
var create = &cobra.Command{
    Use:   "create",
    Short: "Creates a server in selected provider",
    Long:  `Will walk you through a wizard to create a server in a selected provider`,
    Run: func(cmd *cobra.Command, args []string) {
        selectedProvider, err := utils.AskForProvider()
        if err != nil {
            color.Yellow("Something went wrong asking for selected provider\n")
            return
        }
        
        // DigitalOcean
        if selectedProvider == "DO" {
            ctx := context.Background()
            credManager := credentials.NewManager()
            token, _, err := credManager.GetToken(ctx)
            if err != nil {
                color.Red("✗ Error: Unable to get DigitalOcean API token\n")
                fmt.Println()
                fmt.Println("Run 'cogo config set-token' to configure your token.")
                return
            }
            
            client := godo.NewFromToken(token)
            createdDroplet, err := do.ExecuteCreateFlow(client)
            
            if err != nil {
                color.Red("✗ Error: %v\n", err)
                return
            }
            if createdDroplet == nil {
                return
            }
            
            color.Green("✓ Droplet [%s] was created!", createdDroplet.Name)
            color.Cyan("List your droplets in a couple of minutes to see the IP\n")
        }
        
        // AWS - NEW PROVIDER
        if selectedProvider == "AWS" {
            // Get AWS credentials
            accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
            secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
            region := os.Getenv("AWS_REGION")
            
            if accessKey == "" || secretKey == "" {
                color.Red("✗ Error: AWS credentials not configured\n")
                fmt.Println()
                fmt.Println("Set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables.")
                return
            }
            
            client, err := aws.NewClient(region, accessKey, secretKey)
            if err != nil {
                color.Red("✗ Error: Failed to create AWS client: %v\n", err)
                return
            }
            
            instance, err := aws.ExecuteCreateInstanceFlow(client)
            
            if err != nil {
                color.Red("✗ Error: %v\n", err)
                return
            }
            if instance == nil {
                return
            }
            
            color.Green("✓ EC2 Instance [%s] was created!", *instance.InstanceId)
            color.Cyan("List your instances in a couple of minutes to see the IP\n")
        }
    },
}
```

## Provider Consistency Checklist

Ensure your provider implementation has:

### Navigation Patterns
- [ ] Uses navigation framework for all multi-step operations
- [ ] All steps implement `Step` interface
- [ ] Back navigation support in all flows
- [ ] Empty state handling in all steps
- [ ] Standard error handling (EmptyState, Cancel, Canceled)

### User Experience
- [ ] Same keyboard shortcuts as DigitalOcean (Ctrl+C, Esc, b, q)
- [ ] Same message formats (✓ success, ✗ error, ⚠️ warning)
- [ ] Same colored output patterns
- [ ] Summary before confirmation
- [ ] Default to "No" for destructive operations

### Error Handling
- [ ] Exit code 0 for empty states and cancellation
- [ ] Exit code 1 for real errors
- [ ] Clear, actionable error messages
- [ ] No panics or crashes

### Documentation
- [ ] Provider added to README
- [ ] Configuration instructions provided
- [ ] Example commands documented
- [ ] Help text follows standard format

### Testing
- [ ] Unit tests for each step
- [ ] Integration tests for full flows
- [ ] Empty state tests
- [ ] Cancellation tests
- [ ] Back navigation tests

## Credential Management

For new providers, integrate with the credentials system:

**File**: `credentials/{provider}.go`

```go
package credentials

type AWSProvider struct {
    accessKey string
    secretKey string
}

func NewAWSProvider() *AWSProvider {
    return &AWSProvider{}
}

func (p *AWSProvider) GetCredentials(ctx context.Context) (accessKey, secretKey string, err error) {
    // Check environment variables
    accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
    secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
    
    if accessKey != "" && secretKey != "" {
        return accessKey, secretKey, nil
    }
    
    // Check AWS credentials file
    // ...
    
    // Prompt user
    // ...
    
    return "", "", ErrTokenNotFound
}
```

## Testing Your Provider

Create comprehensive tests:

```go
func TestAWSCreateFlow_EmptyAMIs(t *testing.T) {
    client := &mockAWSClient{amis: []*ec2.Image{}}
    flow := NewCreateInstanceFlow(client)
    
    navigator := navigation.NewNavigator()
    _, err := navigator.Run(context.Background(), flow)
    
    assert.ErrorIs(t, err, navigation.ErrEmptyState)
}

func TestAWSCreateFlow_BackNavigation(t *testing.T) {
    // Test going back from each step
}

func TestAWSCreateFlow_FullFlow(t *testing.T) {
    // Test complete instance creation
}
```

## Example Providers to Study

1. **DigitalOcean** (`digitalocean/`) - Complete reference implementation
   - `create_flow.go` - 7-step creation flow
   - `destroy_flow.go` - 4-step destruction flow

2. **Future AWS** (`aws/`) - Your implementation!

3. **Future Azure** (`azure/`) - Follow same patterns

## Common Pitfalls

❌ **Don't**:
- Skip empty state checking
- Use different keyboard shortcuts
- Implement custom navigation (use framework)
- Use different message formats
- Default to "Yes" for destructive operations

✅ **Do**:
- Follow DigitalOcean patterns exactly
- Use navigation framework
- Handle all error cases
- Write comprehensive tests
- Document configuration requirements

## Getting Help

- Study: `digitalocean/create_flow.go` (reference implementation)
- Read: `specs/003-consistent-cli-navigation/navigation-patterns.md`
- Check: `specs/003-consistent-cli-navigation/contracts/examples.md`
- Ask: Project maintainers

## Summary

Adding a new provider to cogo requires:
1. Package structure matching existing providers
2. Navigation framework integration
3. Consistent UX patterns
4. Comprehensive testing
5. Proper documentation

**Key principle**: Users should not be able to tell which cloud provider they're using based on the CLI navigation experience. All providers must feel identical.

