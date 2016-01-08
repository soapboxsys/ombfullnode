ombfullnode
===========

This is a fork of btcsuite's full node
[btcd](https://github.com/btcsuite/btcd). 
Operating instructions can be found in their documentation.
This fork pulls in releases from btcd.

## Requirements

[Go](http://golang.org) 1.3 or newer.

- Install Go according to the installation instructions here:
  http://golang.org/doc/install


## Installation

#### Linux/BSD/MacOSX/POSIX - Build from Source

- Run the following commands to obtain ombfullnode, all dependencies, and
  launch it:

  ```bash
# Download and configure the dependencies
  > apt-get install gcc git
  > # Install a version of go > 1.3
# Download and build the required packages.
  > go get -v -u github.com/soapboxsys/ombfullnode/...
# Launch it
  > ombfullnode
  ```

- ombfullnode (and utilities) will now be installed in either ```$GOPATH/bin```
If you have not already added the bin directory to your system's path when you
setup go, we recommend you do so now.


