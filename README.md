# claptrap

Monitor file/folder, and trigger whatever action you wish. HTTP notification, file backup
or anything that cross your mind. With the go plugin you can implement your own magic !

> Built easily thanks to [fsnotify](https://> github.com/fsnotify/fsnotify) :sunglasses:.

- Download

```bash
$ go get github.com/TommyStarK/claptrap
```

- Demo

First, let's build the `shared object` based on the example plugin. This is a dummy plugin
acting as a simple log function that prints on the standard output:

    - the type of event which can be either CREATE|UPDATE|RENAME|REMOVE
    - the target file responsible of the event
    - the timestamp of when the event has been detected

```bash
$ cd example/
$ go build -buildmode=plugin -o demo
$ cd ../
```

Our plugin is ready, we can now build `claptrap` and run it:

```bash
# build claptrap
$ go build -mod=vendor -o claptrap

# see help
$ ./claptrap --help

# for demo purposes we run claptrap in its own directory
$ ./claptrap -path=. -plugin=example/demo
```

It's ready !! To test it you can edit the README and remove this line :sunglasses:.


- Write your plugin

Take a look at the [example](https://github.com/TommyStarK/claptrap/blob/master/example) directory
to learn how to write your first plugin :smile:.