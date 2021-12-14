# Prerequisites

## Enable CGO
Because here we use the CGO package for *sqlite* driver, please enable the environment variable `CGO_ENABLED=1` and have a `gcc` compile present within your path.

## Slack token generation
Please [see details here](slack.md).

## Install sqlite3
If you want to use the sqlite as main database, please install the sqlite extension to your system.
You can use this command for ubuntu
```
sudo apt-get install sqlite3 libsqlite3-dev
```
Or by **using brew**
```
brew install sqlite
```
Or for **Centos**
```
sudo yum install sqlite
```