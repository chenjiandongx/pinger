# logrwap

[![GoDoc](https://godoc.org/github.com/digineo/go-logwrap?status.svg)](https://godoc.org/github.com/digineo/go-logwrap)
[![Go Report Card](https://goreportcard.com/badge/github.com/digineo/go-logwrap)](https://goreportcard.com/report/github.com/digineo/go-logwrap)

A thin layer around Go's standard library logger. It was extracted from
`github.com/digineo/fastd/fastd`.

Consumers of this package should simply instantiate a log instance, like
so:

```go
import "github.com/digineo/go-logwrap"

var (
	log       = &logwrap.Instance{}
	SetLogger = log.SetLogger
)
```

This then allows your consumers to plug in their own favorite logger:

```go
import (
	"github.com/digineo/fastd/fastd"
	"github.com/sirupsen/logrus"
)

func init() {
	fastd.SetLogger(logrus.WithField("component", "pkg"))
}
```

You are not limited to `sirupsen/logrus` -- your logger only needs to
implement this interface:

```go
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
```

where both methods implement `fmt.Printf` semantics.

This should be compatible with:

- [sirupsen/logrus](https://github.com/sirupsen/logrus) - both `*logrus.Entry` and `*logrus.Logger`
- [uber-go/zap](https://github.com/uber-go/zap) - `*zap.SugaredLogger`
- [google/logger](https://github.com/google/logger) - `*logger.Logger`
