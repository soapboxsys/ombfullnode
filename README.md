ombfullnode
===========

This is a fork of btcsuite's full node [btcd](https://github.com/btcsuite/btcd). 
Operating instructions can be found in their documentation.
Ombfullnode is based off of the 0.10.0 release.
This project will not receive updates for the foreseeable future.

## Requirements

[Go](http://golang.org) 1.3 or newer.

- Install Go according to the installation instructions here:
  http://golang.org/doc/install


## Installation

#### Linux/BSD/MacOSX/POSIX - Build from Source

- Run the following commands to obtain ombfullnode, all dependencies, and install it:

```bash
# Just download the required packages.
> go get -d github.com/soapboxsys/ombfullnode/...
# Move into the workspace's path
> cd $GOPATH/src/github.com/soapboxsys/ombfullnode
# Use godep to checkout the correct dependent library commits
> godep restore
# Move into your $GOPATH binary directory and build the binary
> cd $GOPATH/bin/
> go build github.com/soapboxsys/ombfullnode/...
```


- ombfullnode (and utilities) will now be installed in either ```$GOROOT/bin``` or
  ```$GOPATH/bin``` If you did not already
  add the bin directory to your system path during Go installation, we
  recommend you do so now.

