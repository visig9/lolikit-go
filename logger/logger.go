package logger

import (
	"log"
	"os"
)

// Err is a logger for standard error messages (for user).
var Err = log.New(os.Stderr, os.Args[0]+": ", 0)

// Debug is a logger for debug.
// var Debug = log.New(os.Stderr, os.Args[0]+": ", log.Lshortfile)

// Std is a logger for std output.
var Std = log.New(os.Stdout, "", 0)
