# claptrap


Monitor files/folders, and trigger whatever action you wish like HTTP notification, file backup and so so forth. With the go plugin you can implement your own magic !

Built easily thanks to [fsnotify](https://github.com/fsnotify/fsnotify) :sunglasses:.

## Install

- Download

```bash
$ go get github.com/TommyStarK/claptrap
```

- Build 


## Usage

```bash
$ ./claptrap --help
Usage of ./claptrap:
  -path string
    	specify the path to the file/directory to watch
  -plugin string
    	path to the plugin to load (.so)
```

Take a look at the [example](https://github.com/TommyStarK/claptrap/blob/master/example) directory to learn how
to write your own plugin :smile: