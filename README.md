# Cogo

Interactive cloud provider tool. Currently allows you to create a droplet on DigitalOcean, delete a droplet and list your droplets.

I built this as a way to learn Go and have a quick way to make a server without having to go to digitalocean.com

[Contribution Guidelines](./.github/CONTRIBUTING.md)

[![GitHub release](https://img.shields.io/github/v/tag/Midnight-Conqueror/cogo.svg?label=latest)](https://github.com/Midnight-Conqueror/cogo/releases)
[![CircleCI](https://circleci.com/gh/Midnight-Conqueror/cogo.svg?style=svg)](https://circleci.com/gh/Midnight-Conqueror/cogo)
[![Go Report Card](https://goreportcard.com/badge/github.com/Midnight-Conqueror/cogo)](https://goreportcard.com/report/github.com/Midnight-Conqueror/cogo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

![Cogo example gif](https://imgur.com/oLMYjhL.gif)

## Supported providers

- DigitalOcean

## Installing

### From Brew

```bash
brew install Midnight-Conqueror/tap/cogo
```

### From binary

Find the latest release in releases for your machines architecture, download and run either from within the directory `./cogo` or by moving it into the `/usr/local/bin` to be accessed from anywhere.

### Configuration

1. `$HOME/.cogo`
1. `$HOME/.config/.cogo`
1. `./.cogo`

Cogo will look for a file called **`.cogo`**. The file needs to be of **`json`** format.

Current supported config locations are `$HOME/`, `$HOME/.config/` and `./`

See the `sample_config.json` file as a basis.

> It isn't necessary to add the config as cogo will ask you for tokens without a config

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
