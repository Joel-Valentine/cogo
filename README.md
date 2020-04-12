# Cogo

Cogo is an easy way to interact with cloud providers. It currently allows you to create a droplet on DO and list your droplets

![Cogo example gif](https://i.imgur.com/jrxccl7.gif)

## Supported providers

- DigitalOcean

## Installing

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
