# Project build
In this documentation you can find the information about project build

## Build for current system
In these instructions we assume, that you need to build this project for your current system. For build of the project you need to follow the next steps:

``Warning! The following steps will work for MacOs and Linux systems``
1. Clone the latest version of the project to your machine
``` 
git clone git@github.com:sharovik/devbot.git
```
2. Go to the project dir and run next command:
``` 
make build
```
3. If there are no errors, you will see the next binary files
-- `./bin/slack-bot-current-system` - the slack-bot binary file which is ready for run
-- `./scripts/install/run` - installation script
-- `./scripts/update/run` - update script

## Cross platform build
If you want to run cross-platform build, please use the following instructions

### Before cross-platform build
For cross-platform build I use `karalabe/xgo-latest`. So please before project build do the following steps
1. Install `docker` and `go` to your system
2. Run this command `docker pull karalabe/xgo-latest`
3. Every package should have a .go file inside of the directory. `events` folder it is a defined package. There you can configure the list of events which should have the bot. If you will skip this step, the error will appear because the events package should have go files. Please do the following step:
```
cp events/defined-events.go.dist events/defined-events.go
```
This will fix the issue with undefined `events` package, which might happen during project compilation locally.

4. Your project should be in `GOPATH` folder or `GOPATH` should point to the directory where you clone this project

### Build
For build please run this command
``` 
make build-project-cross-platform
```

This command will build the `slack-bot`, `installation-script` and `update-script` for the following versions:
#### MacOS
- darwin-386
- darwin-amd64
#### Linux
- linux-386
- linux-amd64
#### Windows
- windows-386
- windows-amd64

The result of build command you will see in the `project-build` directory.