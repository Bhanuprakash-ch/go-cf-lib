Go Cloud Foundry Library
========================

This is a library to communicate with Cloud Foundry API using REST.
It is written in golang and can be used in applications written in that language.

### IDE
We recommend using [IntelliJ IDEA](https://www.jetbrains.com/idea/) as IDE with [golang plugin](https://github.com/go-lang-plugin-org/go-lang-idea-plugin). To apply formatting automatically on every save you may use go-fmt with [File Watcher plugin](http://www.idmworks.com/blog/entry/automatically-calling-go-fmt-from-intellij).


Tips
-----------------------

### Golang tips

Developing golang apps requires you store all dependencies (Godeps) in separate directory. They shall be placed in source control.

```
godep save ./...
```

Command above places all dependencies from `$GOPATH`, your app uses, in Godeps and writes its versions to Godeps/Godeps.json file.
