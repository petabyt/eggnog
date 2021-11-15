# eggnog
Basic file uploading service in Go.  
Files are XOR encrypted server side, and are only accessible with the key.  
It's not perfect encryption, but it's a whole lot better than storing raw files.


## Installation
Assuming you have `make` and `go` installed,
```
make init # Create "file" folder
make systemd # Create systemd startup file
```
