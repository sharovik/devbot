# Project build
In this documentation you can find the information about project build

## Build for current system
In these instructions we assume, that you need to build this project for your current system. For build of the project you need to follow the next steps:

``Warning! The following steps will work for MacOs and Linux systems``
1. Clone the latest version of the project to your machine
``` 
git clone git@github.com:sharovik/devbot.git
```
2. Run
```
cp events/defined-events.go.dist events/defined-events.go
```
3. Go to the project dir and run next command:
``` 
make build
```
## Cross-platform build
If you want to run cross-platform build, please use the following instructions

### Before cross-platform build
For cross-platform build I use `karalabe/xgo-latest`. So please before project build do the following steps
1. Install `docker` and `go` to your system
2. Run this command `docker pull karalabe/xgo-latest`
3. Execute this command: `cp events/defined-events.go.dist events/defined-events.go`
4. After that, please run this command: `make build-project-cross-platform`

This command will build the `devbot` binary for the following OS versions:
#### MacOS
- darwin-386
- darwin-amd64
#### Linux
- linux-386
- linux-amd64
#### Windows
- windows-386
- windows-amd64

The result of build command you will see in the `bin` directory.