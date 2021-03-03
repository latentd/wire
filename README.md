# wire

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Install

```sh
go get -u github.com/latentd/wire
```

## Usage

```go
func main() {
    r := wire.NewRouter()

    // use middleware
    r.Chain(middleware1, middleware2)

    // register handler
    r.Get("/", handler)
    // or
    r.GetF("/", handlerFunc)

    // create sub router
    sr := r.SubRouter("/sub")
    sr.Get("/articles", handler)

    // use regex
    sr.Get("/articles/(id:[0-9]+)", handler)

    http.Handle("/", r)
}
```

## License

MIT
