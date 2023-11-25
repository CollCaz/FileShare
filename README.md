# FileServer

A simple file hosting server to make storing and sharing files easier.
Shows all files starting at the root directory of the application, and supports the following features:
1. Navigating Directories
2. Uploading files
3. Downloading files
4. Renaming  files
5. Deleting files

## Getting Started
Test the application:

```bash
git clone (repo)
cd (repo name)
go run ./cmd/app/main.go

```
Build the application:

1. For Linux

```bash
make build
```
2. For Windlows
```bash
make buildWindows
```

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application for linux
```bash
make build
```

build the application for windows
```bash
make buildWindows
```
build the application for both
```bash
make build all
```

run the application
```bash
make run
```

clean up binary from the last build
```bash
make clean
```
