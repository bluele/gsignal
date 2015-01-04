package main

import (
	"github.com/bluele/gsignal"
	"log"
	"os"
	"syscall"
	"time"
)

func main() {
	gsg := gsignal.NewWatcher()
	gsg.Watch(func(sig os.Signal) {
		log.Println("Catch signal: ", sig)
		os.Exit(0)
	}, syscall.SIGTERM, os.Interrupt)

	gsg.Watch(func(sig os.Signal) {
		log.Println("Ignore signal: ", sig)
	}, syscall.SIGHUP)

	gsg.Run()

	log.Println("PID: ", os.Getpid())
	log.Println("Sleep...")
	time.Sleep(time.Minute)
}
