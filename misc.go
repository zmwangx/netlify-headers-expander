package main

import (
	"log"
)

func impossiblef(format string, args ...interface{}) {
	log.Panicf("the impossible happened: "+format, args...)
}
