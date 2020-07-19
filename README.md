# egkong
[![PkgGoDev](https://pkg.go.dev/badge/github.com/matthewdunsdon/egkong)](https://pkg.go.dev/github.com/matthewdunsdon/egkong)
![Build go](https://github.com/matthewdunsdon/egkong/workflows/Build%20go/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/matthewdunsdon/egkong/badge.svg?branch=master)](https://coveralls.io/github/matthewdunsdon/egkong?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/matthewdunsdon/egkong)](https://goreportcard.com/report/github.com/matthewdunsdon/egkong)

Package `matthewdunsdon/egkong` provides you with a library to seamlessly integrate [`kong`](https://github.com/alecthomas/kong) with [`egcmd`](https://github.com/matthewdunsdon/egcmd).

Currently, this library is not intended for use in production, predominantly as it has been created learning golang.

---

* [Install](#install)
* [Usage](#usage)
* [License](./LICENSE)

---

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get github.com/matthewdunsdon/egkong
```

## Usage

The simplist approach is to create a configured kingpin app is to use the `.New()` function:

```go
import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/matthewdunsdon/egkong"
)

var cli struct {
	// kong cli
}

var (
	parser, appExamples = egkong.New(&cli, kong.Name("myapp"), kong.Description("This is my app."))
	_                   = appExamples.Example("init", "Ius legimus nonumes te, pri dicat nominavi copiosae id, odio rebum facilis ea pro.")

	initCmd   = app.Command("init", "Initialise cashflow data.")
	initCmdEx = appExamples.Command("init")
	_         = initCmdEx.Example("--yes", "At vis primis debitis, ei verear omittantur signiferumque mei, quo esse aperiri an. Dolore vocent consequuntur pro an, nam no iusto tamquam suscipit.")
)
```

## License

MIT licensed. See the LICENSE file for details.
