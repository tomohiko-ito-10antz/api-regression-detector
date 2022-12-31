package log

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	flags :=
		log.Lmsgprefix |
			log.Ldate |
			log.Ltime |
			log.Llongfile
	logger = log.New(os.Stderr, "", flags)
}

func Stderr(format string, args ...any) {
	logger.Printf(format, args...)
}
