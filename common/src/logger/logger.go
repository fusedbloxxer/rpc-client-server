package logger

import (
	"fmt"
	"log"
)

type NamedEntity interface {
	Name() string
}

type Logger struct {
	Entity NamedEntity
}

func (l *Logger) Log(v ...interface{}) {
	fmt.Print("(", l.Entity.Name(), ") ")
	fmt.Print(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	fmt.Print("(", l.Entity.Name(), ") ")
	log.Fatal(v...)
}