# godcr

[![Build Status](https://github.com/planetdecred/godcr/workflows/Build/badge.svg)](https://github.com/planetdecred/godcr/actions)
[![Tests Status](https://github.com/planetdecred/godcr/workflows/Tests/badge.svg)](https://github.com/planetdecred/godcr/actions)

A cross-platform desktop [SPV](https://docs.decred.org/wallets/spv/) wallet for [decred](https://decred.org/) built with [gio](https://gioui.org/).

## Building

Note: You need to have [Go 1.16](https://golang.org/dl/) or above to build.

Then `go build`.

### Linux
To build **godcr** on Linux these [gio dependencies](https://gioui.org/doc/install/linux) are required.

Arch Linux:
`pacman -S vulkan-headers libxkbcommon-x11`

## FreeBSD
To build **godcr** on FreeBSD you will need to `pkg install vulkan-headers` as root. This is a gio dependency.

## Running godcr
### General usage
By default, **godcr** runs on Mainnet network type. However, godcr can run on testnet by issuing commands on the terminal in the format:
```bash
godcr [options]
```
- Run `./godcr --network=testnet` to run godcr on the testnet network.
- Run `godcr -h` or `godcr help` to get general information of commands and options that can be issued on the cli.
- Use `godcr <command> -h` or   `godcr help <command>` to get detailed information about a command.

## Profiling 
Godcr uses [pprof](https://github.com/google/pprof) for profiling. It creates a web server which you can use to save your profiles. To setup a profiling web server, run godcr with the --profile flag and pass a server port to it as an argument.

So, after running the build command above, run the command

`./godcr --profile=6060`

You should now have a local web server running on 127.0.0.1:6060.

To save a profile, you can simply use

`curl -O localhost:6060/debug/pprof/profile`


## Contributing

See [CONTRIBUTING.md](https://github.com/planetdecred/godcr/blob/master/.github/CONTRIBUTING.md)

## Other

Earlier experimental work with other user interface toolkits can be found at [godcr-old](https://github.com/raedahgroup/godcr-old).
