package teleport

import (
	"log"
)

var Debug bool

func debugPrintf(format string, v ...interface{}) {
	if !Debug {
		return
	}
	log.Printf(format, v...)
}

func debugPrintln(v ...interface{}) {
	if !Debug {
		return
	}
	log.Println(v...)
}

func debugFatal(v ...interface{}) {
	if !Debug {
		return
	}
	log.Fatal(v...)
}
