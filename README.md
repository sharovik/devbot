# devbot
[![Gitter](https://badges.gitter.im/devbot-tool/community.svg)](https://gitter.im/devbot-tool/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

This bot can help to automate multiple processes of development and give the possibility to achieve more goals for less time.

## Table of contents
- [Available features](#generate-wordpress-template)
- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Available events](#available-events)
- [Custom events](#custom-events)
- [Dictionary](#dictionary)
- [Cross platform build](#cross-platform-build)
- [Authors](#authors)
- [License](#license)

## Available features
* [trigger custom events by personal message or triggering bot in the channel](documentation/events.md)
* [create the WordPress template just by uploading the file to the specific channel or to the PM of the bot](#generate-wordpress-template)

## Getting Started

These instructions will help you to install the bot to your server (local, development, production).

## Prerequisites

### Enable CGO
Because here we use the CGO package for *sqlite* driver, please enable the environment variable `CGO_ENABLED=1` and have a `gcc` compile present within your path.

### Slack token generation
Before the installation I would recommend to prepare the slack application for your account. 
1. Go to [applications page](https://api.slack.com/apps?new_classic_app=1) of slack and create new application there
2. Once new application was created you will be redirected to the application `Basic Information` page, where you have to click in the `Building Apps for Slack` section to the `Add features and functionality` block. There you need to click to the `Bots` button.
3. Add a Bot user. Specify the `Display name`, `Default username` and his `online status`
4. After you created a bot user, please go back to the `Basic Information` page and install your app to your workspace. You can find the `Install your app to your workspace` button in the `Building Apps for Slack` section.
5. Now you need to get the OAuth tokens for our bot user. For that please go to `OAuth & Permissions`, there you will find the `Bot User OAuth Access Token` which appears only after application installation to your slack account. This token you will need to specify in .env configuration file of your bot

### Install sqlite3
We are using the sqlite3 as main storage of our questions and answers data. To exclude errors related to unknown library sqlite3, please install it.
You can use this command for ubuntu
```
sudo apt-get install sqlite3 libsqlite3-dev
```
Or by using brew
```
brew install sqlite
```
Or for centos
```
sudo yum install sqlite
```

### PHP installation
You server requires php version of 7.1+ with php-dom module. `It is only required if you will use the wordpress template generation event.`
For ubuntu
```
sudo apt install php php-dom
```
Or for brew
```
brew install php
```
Or for centos
```
yum install php php-xml
```
## Installation
You can easily install the devbot application by using the installation script.

1. Go to [this page](https://github.com/sharovik/devbot) and download/clone the latest version of devbot. Or run this command locally:
2. Run this command to install everything related to the database and configuration
``` 
make build-installation-script && scripts/install/run
```
3. Set the value from [`Bot User OAuth Access Token`](#slack-token-generation) into *SLACK_OAUTH_TOKEN* variable from the `.env` file
4. Run bot by using command `./bin/slack-bot-{YOUR_SYSTEM}` you should see in the logs `hello` message type. It means that the bot successfully connected to your account
![Demo start slack-bo](documentation/images/start-slack-bot.gif)

## How to use
Basically you can do whatever you want with the chat-events. All depends on your imagination and on your daily basis workflow.

[You can find an example of the event here](events/example/README.md).

## Custom events
Please read the [events documentation](documentation/events.md)

## Dictionary
Please read the [dictionary documentation](documentation/dictionary.md)

## Cross platform build

### Before build
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
make build
```
This command will build the following versions:
#### MacOS
- darwin-386
- darwin-amd64
#### Linux
- linux-386
- linux-amd64
#### Windows
- windows-386
- windows-amd64

## Authors

* **Pavel Simzicov** - *Initial work* - [sharovik](https://github.com/sharovik)

### Vendors used
* github.com/joho/godotenv - for env files loading
* github.com/karalabe/xgo - for cross platform build
* github.com/karupanerura/go-mock-http-response - for http responses mocking in tests
* github.com/mattn/go-sqlite3 - for sqlite connection
* github.com/pkg/errors - for errors wrapper and trace extracting in logger
* github.com/rs/zerolog - for logger
* github.com/stretchr/testify - for asserts in tests
* golang.org/x/net - for websocket connection

## License
This project is licensed under the BSD License - see the LICENSE.md file for details
