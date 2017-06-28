package main

import(
	"os/signal"
	"os"
	"log"
	"./web"
	"./scanner"
	"./config"
)

func main() {
	subnets, params, e := config.ParseConfig("/etc/ipreg.conf")
	if e != nil {
		log.Fatal(e.Error())
		return
	}
	go scanner.StartScanner(subnets, params.TimeBetweenScansInMinutes,
		params.MaxParallelJobs)
	if err := web.Initialize(subnets, params); err != nil {
		log.Fatal(err)
	}
	go web.StartServer()
	waitForSignal()
	web.StopServer()
	scanner.StopScanner()
}

func waitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<- c
	log.Println("Got signal, exiting")
}
