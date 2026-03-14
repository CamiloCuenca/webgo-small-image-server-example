# WebGo Simple Image Server

A small example HTTP server written in Go.
This project demonstrates how to build a lightweight web server that serves images and web content with embedded files.

All required files (such as images or static resources) are embedded in the executable, so the server can run without needing external folders.

## Purpose

This project is intended as a simple learning example to understand how to:

* Create a basic HTTP server in Go
* Embed files inside the executable
* Run a portable server without external dependencies

## Usage

### Run from source

If you have Go installed, you can run the server with:

```bash
go run main.go
```

### Run the compiled executables

Precompiled executables are provided for convenience.

**Linux**

```bash
./server
```

**Windows**

```bash
server.exe
```

The server will start locally and can be accessed from a web browser.

## Notes

Because the files are embedded, the executable can run independently without needing additional files or directories.
