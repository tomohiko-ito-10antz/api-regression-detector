package log

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "",
	log.Lmsgprefix|
		log.Ldate|
		log.Ltime|
		log.Llongfile)

func Stderr(format string, args ...any) {
	callPath := 3

	err := logger.Output(callPath, fmt.Sprintf(format, args...))
	if err != nil {
		logger.Fatalln(err)
	}
}
