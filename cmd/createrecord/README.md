CreateRecord
============

Creates a public record sqlite db from a leveldb block database. Works with 
testnet or mainnet. Will eventually be extended to support bitcoind's blocks.

Usage
-----

```
./createrecord --testnet --overwrite -s=364668
```
This cmd will look in the default data dir and overwrite the existing database 
starting from block 364668.
