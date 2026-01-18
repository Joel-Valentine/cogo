# Cogo

**Interactive cloud provider tool** with modern, intuitive navigation.

Create, list, and destroy cloud resources with an intelligent CLI that supports back navigation, graceful empty state handling, and consistent keyboard shortcuts across all operations.

**Features**:
- 🔙 **Back Navigation** - Go back and change selections (inspired by gcloud)
- 🎯 **Smart Empty State Handling** - No crashes, only helpful messages
- ⌨️ **Universal Keyboard Shortcuts** - Ctrl+C, Esc, 'b' for back, 'q' to quit
- 🎨 **Colored Output** - ✓ success, ✗ error, ⚠️ warning
- 🔐 **Secure Credentials** - OS keychain integration
- 📝 **Multi-step Flows** - Guided wizards with state preservation

[Contribution Guidelines](./.github/CONTRIBUTING.md)

[![CI](https://github.com/Joel-Valentine/cogo/actions/workflows/ci.yml/badge.svg)](https://github.com/Joel-Valentine/cogo/actions/workflows/ci.yml)
[![GitHub release](https://img.shields.io/github/v/tag/Joel-Valentine/cogo.svg?label=latest)](https://github.com/Joel-Valentine/cogo/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/Joel-Valentine/cogo)](https://goreportcard.com/report/github.com/Joel-Valentine/cogo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

![Cogo example gif](https://imgur.com/oLMYjhL.gif)

## Supported providers

- DigitalOcean

## Installing

### From Brew

```bash
brew install Joel-Valentine/tap/cogo
```

### From binary

Find the latest release in releases for your machines architecture, download and run either from within the directory `./cogo` or by moving it into the `/usr/local/bin` to be accessed from anywhere.

### Configuration

Cogo uses a modern, secure credential management system with multiple storage options.

#### 🔐 Secure Storage (Recommended)

Store your DigitalOcean API token securely in your OS keychain:

```bash
cogo config set-token
```

Your token will be stored in:
- **macOS**: Keychain
- **Windows**: Credential Manager
- **Linux**: Secret Service (GNOME Keyring, KWallet, etc.)

#### 🌍 Environment Variables

For CI/CD or automation, use environment variables:

```bash
export DIGITALOCEAN_TOKEN=dop_v1_xxx
cogo create
```

Supported environment variables (in priority order):
- `DIGITALOCEAN_TOKEN` (standard DigitalOcean env var)
- `COGO_DIGITALOCEAN_TOKEN` (cogo-specific)

#### 📁 Legacy File Configuration (Deprecated)

Cogo still supports the legacy `.cogo` JSON file for backward compatibility:

Locations checked (in order):
1. `$HOME/.cogo`
2. `$HOME/.config/.cogo`
3. `./.cogo`

**⚠️ Warning**: File storage is insecure (plain text). Migrate to keychain:

```bash
cogo config migrate
```

#### Priority Order

Cogo checks for credentials in this order:
1. Command-line flag: `--token`
2. Environment variable: `DIGITALOCEAN_TOKEN`
3. Environment variable: `COGO_DIGITALOCEAN_TOKEN`
4. OS Keychain (secure)
5. Config file (legacy)
6. Interactive prompt

#### Configuration Commands

```bash
# Set token (stores in keychain)
cogo config set-token

# View current token (masked)
cogo config get-token

# Check configuration status
cogo config status

# Migrate from file to keychain
cogo config migrate

# Delete stored token
cogo config delete-token
```

## Usage

### Keyboard Shortcuts

**Universal shortcuts work in all commands**:

| Key | Action |
|-----|--------|
| **Ctrl+C** | Cancel immediately |
| **Esc** or **q** | Quit current operation |
| **b** or **←** | Go back to previous step |
| **↑** / **↓** | Navigate lists |
| **Enter** | Confirm / Continue |

### create

Create a droplet with an interactive, multi-step wizard:

```bash
cogo create
```

**The wizard guides you through**:
1. Choose your provider (DigitalOcean)
2. Enter a name (with smart default)
3. Choose image type (Distributions/Applications/Custom)
4. Select specific image
5. Select droplet size
6. Select region
7. Select SSH key
8. Review summary and confirm

**Navigation Features**:
- 🔙 Press **'b'** to go back and change any selection
- 🎯 Empty states show helpful messages (no crashes)
- ✓ Summary displayed before confirmation
- ⚠️ Clear, colored output throughout

**Example**:
```bash
$ cogo create

? Select provider: DigitalOcean

? Droplet Name: (my-droplet-1737238400)
> my-awesome-server

? Select Image Type: (Use arrow keys, 'b' for back, 'q' to quit)
  ← Back
  › Distributions
    Applications
    Custom

[Press 'b' if you want to change your name]

=== Droplet Configuration ===
Name:     my-awesome-server
Image:    Ubuntu 22.04 LTS
Size:     s-1vcpu-1gb
Region:   nyc3
SSH Key:  my-key
=============================

? Create this droplet? (Y/n)

✓ Droplet [my-awesome-server] was created!
```

### list

Lists all droplets in your account with clear, colored output:

```bash
cogo list
```

**Output**:
```bash
Your droplets:

0  Name: blog
   IP: 192.168.1.100

1  Name: backend
   IP: 192.168.1.101

2  Name: frontend
   IP: 192.168.1.102
```

**Empty state handling**:
```bash
No droplets found in your DigitalOcean account.

Run 'cogo create' to create a droplet.
```

### destroy

Delete a droplet with multiple safety confirmations and back navigation:

```bash
cogo destroy
```

**Multi-step safety flow**:
1. Select droplet to delete
2. First confirmation (y/n)
3. Re-enter droplet name (prevents accidents)
4. View full droplet details
5. Final confirmation (y/n)

**Navigation Features**:
- 🔙 Press **'b'** to go back if you change your mind at any step
- ⚠️ Multiple warnings before destructive action
- ✓ Full droplet details shown before deletion
- 🎯 No crashes if no droplets exist

**Example**:
```bash
$ cogo destroy

? Select droplet to delete: (Use arrow keys, 'b' for back, 'q' to quit)
  ← Back
  › my-droplet-1 (192.168.1.1)
    my-droplet-2 (192.168.1.2)

⚠️  WARNING: You are about to delete droplet: my-droplet-1

? Are you sure? (y/N)

? Re-enter droplet name to confirm delete: my-droplet-1

=== Droplet to be DELETED ===
Name:   my-droplet-1
Size:   s-1vcpu-1gb
Region: nyc3
IP:     192.168.1.1
=============================

? Are you really really sure you want to delete this droplet? (y/N)

✓ Droplet [my-droplet-1] has been destroyed
```

## Installing from source

This project requires Go to be installed.

Running it then should be as simple as cloning the repository then:

```console
$ make build
$ ./bin/cogo
```

### Testing

`make test`

## Contributing

If you've read this far you're probably the right person to add to this project

Please read the [contributing](.github/CONTRIBUTING.md) guide on how to get started

I am still learning Go at the moment so don't feel like you need to be a wizard to contribute
