package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wats/internal/database"
	"wats/internal/trader"
)

func main() {
	// TODO: cfg 로직 수정 후 연결 필요
	db := database.NewDatabase()
	t := trader.NewTrader(db, "SUIUSDT")

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
