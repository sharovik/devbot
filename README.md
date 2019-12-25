# devbot

This bot can help to automate multiple processes of development and give the possibility to achieve more goals for less time.

### What this bot can
* create the WordPress template just by uploading the file to the specific channel or to the PM of the bot

## Getting Started

These instructions will help you to install the bot to your server (local, development, production).

### Prerequisites

Before the installation I would recommend to prepare the slack application for your account. 
1. Go to [applications page](https://api.slack.com/apps?new_app=1) of slack and create new application there
2. Once new application was created you will be redirected to the application `Basic Information` page, where you have to click in the `Building Apps for Slack` section to the `Add features and functionality` block. There you need to click to the `Bots` button.
3. Add a Bot user. Specify the `Display name`, `Default username` and his `online status`
4. After you created a bot user, please go back to the `Basic Information` page and install your app to your workspace. You can find the `Install your app to your workspace` button in the `Building Apps for Slack` section.
5. Now you need to get the OAuth tokens for our bot user. For that please go to `OAuth & Permissions`, there you will find the `Bot User OAuth Access Token` which appears only after application installation to your slack account. This token you will need to specify in .env configuration file of your bot

### Devbot installation

1. Go to [bin folder of this project](https://github.com/sharovik/devbot/tree/master/bin) and download latest version of devbot application.
2. Prepare the configuration file for our bot
```
cp .env.example .env
```
3. Set the value from `Bot User OAuth Access Token` into *SLACK_OAUTH_TOKEN* variable in .env file

