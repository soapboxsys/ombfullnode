ombfullnode
===========

This is a fork of btcsuite's full node [btcd](https://github.com/btcsuite/btcd). 
Operating instructions can be found in their documentation.
Ombfullnode is based off of the 0.10.0 release.
This project will not receive updates for the foreseeable future.

## Requirements

[Go](http://golang.org) 1.3 or newer.

## Installation

#### Linux/BSD/MacOSX/POSIX - Build from Source

- Install Go according to the installation instructions here:
  http://golang.org/doc/install

- Run the following command to obtain btcd, all dependencies, and install it:

```bash
$ go get github.com/soapboxsys/ombfullnode/...
```

- ombfullnode (and utilities) will now be installed in either ```$GOROOT/bin``` or
  ```$GOPATH/bin``` depending on your configuration.  If you did not already
  add the bin directory to your system path during Go installation, we
  recommend you do so now.

