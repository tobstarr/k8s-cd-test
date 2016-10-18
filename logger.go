package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Logger interface {
	Printf(string, ...interface{})
}

type logger struct {
	Prefixes []string
	Timed    bool
	Writer   io.Writer
	started  time.Time
	reader   io.ReadCloser
	writer   io.WriteCloser
	c        chan error
}

func newLogger(opts ...loggerOpt) *logger {
	l := &logger{
		started: time.Now(),
		Writer:  os.Stderr,
	}
	for _, o := range opts {
		o(l)
	}
	return l
}

type loggerOpt func(*logger)

func timed(o *logger) {
	o.Timed = true
}

func (l *logger) Printf(m string, args ...interface{}) {
	fmt.Fprint(l.Writer, l.prefix()+fmt.Sprintf(m, args...)+"\n")
}

func (l *logger) prefix() string {
	p := []string{}
	if l.Timed {
		p = append(p, "["+formatTime(time.Since(l.started).Seconds(), 3)+"]")
	}
	for _, prefix := range l.Prefixes {
		p = append(p, "["+prefix+"]")
	}
	if len(p) == 0 {
		return ""
	}
	return strings.Join(p, "") + " "
}

func formatTime(in float64, digits int) string {
	s := fmt.Sprintf("%.03f", in)
	return fmt.Sprintf("%*s", digits+1+3, s)
}
