# ANMplugin cli
This tool allows you to safely and comfortably develop plugins for *hidden*

## Usage
This tool is pretty simple you just put its binary that is built for your system, inside your plugin directory and call it like this
```
./anmplugin.exe serve
```
or like this for help
```
./anmplugin.exe help
```
so for linux 
```
./anmplugin-linux serve
```
and for mac
```
./anmplugin-mac serve
```
## Useful features
to not send an unnecessary or conifdential data into your zip file you can just add `.ignore` file which works like any normal `.gitignore` file , example config
```
example-directory/
example-file.txt
anmplugin-linux
anmplugin-mac
anmplugin.exe
anmplugin-mac-arm
```
for now there's no comment support, but in theory you can still comment as it ignores any invalid path
