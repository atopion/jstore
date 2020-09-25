## JStore - a small JSON file server build in Go

JStore is a small JSON file server build in Go.

JStore enables 3 basic functions for storing JSON files on a server:

 - View a file using GET
 - Store a file using POST
 - Modify a file using PUT

Files are stored using UUIDs as identifiers, by which they can be accessed.

### Installation
- Use `go build` and `go install` to build the source code and install the program. Afterwards, the program can be run with `jstore`.
- Start the precompiled program (currently only Windows/amd64) by running `.\jstore.exe`. 
- Build and run the docker image by using the Dockerfile

### Arguments
- `-f`: Path to the folder in which `jstore` should store the files.
- `-p`: Port on which `jstore` should run
- `-u`: Base URL that will be returned together with the identifier when a new file is stored.