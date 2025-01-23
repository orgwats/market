package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wats/config"
	"wats/internal/trader"
)

func main() {
	c := config.LoadConfig()
	t := trader.NewTrader(c)

	go func() {
		listenForSignals()
		log.Println("[main] Received OS signal, stopping trader...")

		t.Stop()

		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()

	t.Start()

	select {}
}

func listenForSignals() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
}
