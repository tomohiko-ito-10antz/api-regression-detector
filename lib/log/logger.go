package log

import (
	"fmt"
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
	err := logger.Output(3, fmt.Sprintf(format, args...))
	logger.Fatalln(err)
}
