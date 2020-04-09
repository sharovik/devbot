# Prerequisites

## Enable CGO
Because here we use the CGO package for *sqlite* driver, please enable the environment variable `CGO_ENABLED=1` and have a `gcc` compile present within your path.

## Slack token generation
Before the installation I would recommend to prepare the slack application for your account. 
1. Go to [applications page](https://api.slack.com/apps?new_classic_app=1) of slack and create new application there
2. Once new application was created you will be redirected to the application `Basic Information` page, where you have to click in the `Building Apps for Slack` section to the `Add features and functionality` block. There you need to click to the `Bots` button.
3. Add a Bot user. Specify the `Display name`, `Default username` and his `online status`
4. After you created a bot user, please go back to the `Basic Information` page and install your app to your workspace. You can find the `Install your app to your workspace` button in the `Building Apps for Slack` section.
5. Now you need to get the OAuth tokens for our bot user. For that please go to `OAuth & Permissions`, there you will find the `Bot User OAuth Access Token` which appears only after application installation to your slack account. This token you will need to specify in .env configuration file of your bot
6. Set the value from `Bot User OAuth Access Token` into *SLACK_OAUTH_TOKEN* variable in the `.env` file

## Install sqlite3
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