# Experimenting with go-astilectron
Electron + Go sounds interesting, so I am experimenting with what i can do with it.

Using the tools provided by: https://github.com/asticode/go-astilectron

Using the demo here as a base : https://github.com/asticode/go-astilectron-demo

Run the following to build locally
```
astilectron-bundler -v
```




Both go install and go build will compile the package in the current directory when ran without additional arguments.

If the package is package main, go build will place the resulting executable in the current directory. go install will put the executable in $GOPATH/bin (using the first element of $GOPATH, if you have more than one).

Build and install on this package
```
go install <package>
```

Build package (won't produce an output file. Instead it saves the compiled package in the local build cache.)
```
go build <package>
```


Run tests on this package
```
go test <package>
```
Example
```
go test github.com/user/stringutil
```


# Developer notes

the app *astilectron.Astilectron is set in the onWait function