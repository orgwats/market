package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/orgwats/market/internal/config"
	db "github.com/orgwats/market/internal/db/sqlc"
	"github.com/orgwats/market/internal/gapi"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/orgwats/idl/gen/go/market"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbCfg := cfg.Service.Market.Database
	conn, err := sql.Open(dbCfg.Type, dbCfg.DSN())
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server := gapi.NewServer(cfg, store)
	grpcServer := grpc.NewServer()

	pb.RegisterMarketServer(grpcServer, server)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Market.Port))
	if err != nil {
		log.Fatal("cannot listen network address:", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("start market server at %s", listener.Addr().String())

		server.Run()
		defer server.Stop()

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("market server failed to serve:", err)
		}
	}()
	wg.Wait()
}
