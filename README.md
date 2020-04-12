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

Create will run you through creating a droplet on your given cloud provider

```bash
cogo create
```

### list

list will list servers created on that provider printing the name and IP

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

Please add an issue if you want to fix a bug or add a feature allowing us to initially talk through the request. Everyone is welcome to contribute to this project if they have a valid feature/bug.

I am still learning Go at the moment so don't feel like you need to be a wizard to contribute
