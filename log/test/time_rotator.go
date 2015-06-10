package main

import (
	"log"
	"time"

	rotator "github.com/ggicci/jungo/log"
)

func main() {
	tr, err := rotator.NewTimeRotator("/tmp/logtest/my.log", rotator.Hourly)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(tr)

	for {
		time.Sleep(time.Second)
		log.Print(time.Now().Format(time.RFC3339Nano))
	}
}
