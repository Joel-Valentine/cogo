# Cogo

Interactive cloud provider tool. Currently allows you to create a droplet on DigitalOcean, delete a droplet and list your droplets.

I built this as a way to learn Go and have a quick way to make a server without having to go to digitalocean.com

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

### create

Create will run you through creating a droplet on your given cloud provider. Currently the process is:

1. Chose your provider
1. Enter a name
1. Chose an image
1. Chose a region
1. Chose a size
1. Chose an ssh key
1. Are you sure (y/n)

Finally you will be told the droplet has been created. You can then list your servers from that provider once you think its been created / assigned an IP.

```bash
cogo create
```

### list

list will list servers created on that provider printing the name and IP

```bash
cogo list


Your droplets:

0  Name: blog
   IP: xxx.xxx.xxx.xxx

1  Name: backend
   IP: xxx.xxx.xxx.xxx

2  Name: frontend
   IP: xxx.xxx.xxx.xxx
```

### destroy

Destroy will allow you to delete one of your servers **Safely** there will be a total of three checks to make sure you understand what you are deleting.

1. Chose the provider you wish to delete from
1. There will be an 'are you sure (y/n)' question
1. You will need to enter the name of the server you are deleting
1. You will then have to answer another 'are you really really sure (y\n)' question with details of the server you are about to delete

```bash
cogo destroy
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
