package graalsystems

import (
	"log"
)

// logger is the implementation of the SDK Logger interface for this terraform plugin.
//
// cf. https://godoc.org/github.com/graalsystems/graalsystems-sdk-go/logger#Logger
type logger struct{}

// l is the global logger singleton
var l = logger{}

// Debugf logs to the DEBUG log. Arguments are handled in the manner of fmt.Printf.
func (l logger) Debugf(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}

// Infof logs to the INFO log. Arguments are handled in the manner of fmt.Printf.
func (l logger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

// Warningf logs to the WARNING log. Arguments are handled in the manner of fmt.Printf.
func (l logger) Warningf(format string, args ...interface{}) {
	log.Printf("[WARN] "+format, args...)
}

// Errorf logs to the ERROR log. Arguments are handled in the manner of fmt.Printf.
func (l logger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

// Printf logs to the DEBUG log. Arguments are handled in the manner of fmt.Printf.
func (l logger) Printf(format string, args ...interface{}) {
	l.Debugf(format, args...)
}
