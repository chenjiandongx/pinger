package logwrap

import (
	"fmt"
	"log"
)

// Logger defines log methods used by this package.
type Logger interface {
	Infof(format string, a ...interface{})
	Errorf(format string, a ...interface{})
}

// type check
var _ Logger = (*Instance)(nil)

// Instance allows to instantiate multiple loggers. It is recommended
// to use only one logger instance per package.
type Instance struct {
	l Logger
	o func(int, string) error // is used in tests only
}

// SetLogger updates the logger fastd uses. If l is nil, we'll simply
// wrap the standard log package.
func (i *Instance) SetLogger(l Logger) {
	i.l = l
}

// Infof logs a message with INFO log level. format and args implement
// fmt.Printf semantics.
func (i *Instance) Infof(format string, a ...interface{}) {
	if i.l == nil {
		i.out("INFO", format, a...)
	} else {
		i.l.Infof(format, a...)
	}
}

// Errorf logs a message with ERROR log level. format and args implement
// fmt.Printf semantics.
func (i *Instance) Errorf(format string, a ...interface{}) {
	if i.l == nil {
		i.out("ERROR", format, a...)
	} else {
		i.l.Errorf(format, a...)
	}
}

func (i *Instance) out(level, format string, a ...interface{}) {
	msg := fmt.Sprintf(level+" - "+format, a...)
	if i.o == nil {
		log.Output(3, msg)
		return
	}
	i.o(4, msg)
}
