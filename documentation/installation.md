# DevBot setup  
For proper project installation or update, I would recommend you to run special scripts for install or update.

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
Once you need to update some events or the devbot database schema, you need to use this script for proper installation of updates
Run the script by using next command:
**For MacOS and Linux**
``` 
./scripts/update/update-{YOUR_SYSTEM}
```
For windows
``` 
start scripts\update\update-windows-{TYPE_OF_SYSTEM}.exe
```