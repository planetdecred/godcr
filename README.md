# godcr

[![Build Status](https://github.com/raedahgroup/godcr/workflows/Build/badge.svg)](https://github.com/raedahgroup/godcr/actions)
[![Tests Status](https://github.com/raedahgroup/godcr/workflows/Tests/badge.svg)](https://github.com/raedahgroup/godcr/actions)

A cross-platform desktop [SPV](https://docs.decred.org/wallets/spv/) wallet for [decred](https://decred.org/) built with [gio](https://gioui.org/).

## Building

Note: You need to have [Go 1.13](https://golang.org/dl/) or above to build.

Follow the [installation](https://gioui.org/doc/install) instructions for gio.

Pkger is required to bundle static resources before building. Install Pkger globally by running
`go get github.com/markbates/pkger/cmd/pkger`

In the root directory, run
`pkger`. A pkged.go file should be silently created in your root directory.

Then `go build`.

## Profiling 
Godcr uses [pprof](https://github.com/google/pprof) for profiling. It creates a web server which you can use to save your profiles. To setup a profiling web server, run godcr with the --profile flag and pass a server port to it as an argument.

So, after running the build command above, run the command

`./godcr --profile=6060`

You should now have a local web server running on 127.0.0.1:6060.

To save a profile, you can simply use

`curl -O localhost:6060/debug/pprof/profile`


## Contributing

See [CONTRIBUTING.md](https://github.com/raedahgroup/godcr/blob/master/.github/CONTRIBUTING.md)
