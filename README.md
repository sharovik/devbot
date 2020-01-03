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

1. Go to [this page](https://github.com/sharovik/devbot/releases) and download latest version of devbot application.
2. Prepare the configuration file for our bot
```
cp .env.example .env
```
3. Set the value from `Bot User OAuth Access Token` into *SLACK_OAUTH_TOKEN* variable in .env file
4. Run bot by using command `./bin/slack-bot-{YOUR_SYSTEM}` you should see in the logs `hello` message type. It means that the bot successfully connected to your account

## How to use

### Generate WordPress template
![Demo file processing](documentation/images/demo-file-processing.gif)
1. Prepare a `*.twig` template and compress it in a `*.zip` file ([here you can find the example of template](https://github.com/sharovik/themer)) 
2. Open PM of the bot, channel where this bot was added or tag bot-user during the attachment send
3. In attachment popup please write in comment field `process` word and send the attachment
4. After bot received your attachment you will see the answer
 >Please, wait a bit. I have to process this file
 
 Up to 4 seconds will take to process the template
5. After template was processed you will receive in answer new archive where you will find `html preview of your template` and WordPress template

###Custom events
Please read the [events documentation](documentation/events.md)

## Authors

* **Pavel Simzicov** - *Initial work* - [sharovik](https://github.com/sharovik)

## 

License
This project is licensed under the BSD License - see the LICENSE.md file for details
