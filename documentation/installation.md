# DevBot setup  
You can do it in 2 ways: download the binary files and start the bot or build the project on your machine

## Prerequisites
Before start the installation, [please be aware of prerequisites](prerequisites.md) 

## Install script
This script should be called once you don't have `database.sqlite` file or `.env` file created.
Run the script by using next command:
**For MacOS and Linux**
``` 
./scripts/install/install-{YOUR_SYSTEM}
```
For windows
``` 
start scripts\install\install-windows-{TYPE_OF_SYSTEM}.exe
```

## Update script
Once you need to update some of the events or the devbot database schema, you need to use this script for proper installation of updates
Run the script by using next command:
**For MacOS and Linux**
``` 
./scripts/update/update-{YOUR_SYSTEM}
```
For windows
``` 
start scripts\update\update-windows-{TYPE_OF_SYSTEM}.exe
```