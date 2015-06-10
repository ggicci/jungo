package main

import (
	"log"
	"time"

	"github.com/ggicci/jungo/logutil"
)

func main() {
	tr, err := logutil.NewTimeRotator("/tmp/logtest/my.log", logutil.Hourly)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(tr)

	for {
		time.Sleep(time.Second)
		log.Print(time.Now().Format(time.RFC3339Nano))
	}
}
