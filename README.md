A library to enable random access to multi-block xz files.

Using the algorithm described here
https://blastedbio.blogspot.com/2013/04/random-access-to-blocked-xz-format-bxzf.html

## How to compress

try something like

``` shell
pv file | xz --block-size 1KiB -T 0 > file.xz
```

## How to check if compressed file will work

``` shell
xz --list test.xz 
Strms  Blocks   Compressed Uncompressed  Ratio  Check   Filename
    1   9,766      9.8 MiB  9,765.6 KiB  1.031  CRC64   test.xz
```

to get good random access performance, block size should be relatively small
