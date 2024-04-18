package logger

import (
	"log"
)

func Error(msg error) {
	PrintStack()
	log.Fatalln(msg)
}
