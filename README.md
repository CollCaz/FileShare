# FileServer

A simple file hosting server to make storing and sharing files easier.
Shows all files starting at the root directory of the application, and supports the following actions:
1. Navigating Directories
2. Uploading files, Automatically detects duplicate names and deals with them properly.
3. Downloading files
4. Renaming  files, Automatically detects duplicate names and deals with them properly.
5. Deleting files

Built with Echo, TailwindCSS and HTMX.

![Preview Image](https://github.com/CollCaz/FileShare/blob/main/previewImage.png)

## Getting Started
Test the application:
Download one of the releases (Linux, Windows)
Or build yourself.

```bash
git clone https://github.com/CollCaz/FileShare.git
cd FileShare
go mod tidy
go run cmd/app/main.go

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
