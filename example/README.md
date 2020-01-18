# Example plugin

Take a look at the file [plugin.go](https://github.com/TommyStarK/claptrap/blob/master/example/plugin.go) if you want
to see an example of the simplest plugin for claptrap.

The `Handle` function is the claptrap API. It expectes to find a symbol named `Handle` when looking up
into the shared object specified in the command arguments. It will also ensure that the exported function
match the following signature:

```go
func Handle(action string, file string,  timestamp string) {
    // Add your magic :)
}
```

> :warning: Keep in mind that the `Handle` function will be executed in a separate goroutine
> for each event detected.


To build your plugin, just run the following command:

```bash
$ go build -buildmode=plugin -o NAME_OF_YOUR_PLUGIN
```
