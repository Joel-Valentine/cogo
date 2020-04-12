# cogo

Cogo is an easy way to interact with cloud providers. It currently allows you to create a droplet on DO and list your droplets

## Supported providers

- DigitalOcean

## Installing

Find the latest release in releases for your machines architecture, download and run either from within the directory `./cogo` or by moving it into the `/usr/local/bin` to be accessed from anywhere (not recommended currently)

## Usage

The project will look for a `.cogo_config.json` and will follow the format as per the example.

However it is **not necessary** for it to run, the config will also be created if cogo detects you haven't already made this file. (This is optional)

This project will create the `.cogo_config.json` within the same directory you called cogo from (in the future this needs a rework)

### create

Create will run you through creating a droplet on your given cloud provider. Currently the process is:

1. Chose your provider
2. Enter a name
3. Chose an image
4. Chose a region
5. Chose a size
6. Chose an ssh key
7. Are you sure (y/n)

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

1. There will be an 'are you sure (y/n)' question
2. You will need to enter the name of the server you are deleting
3. You will then have to answer another 'are you really really sure (y\n)' question with details of the server you are about to delete

```bash
cogo list
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

Please read the [contributing](CONTRIBUTING.md) guide on how to get started

I am still learning Go at the moment so don't feel like you need to be a wizard to contribute
