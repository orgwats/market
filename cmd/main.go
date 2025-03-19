package main

import (
	"database/sql"
	"log"
	"net"
	"sync"
	"wats/internal/config"
	db "wats/internal/db/sqlc"
	"wats/internal/gapi"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/orgwats/idl/gen/go/market"
	"google.golang.org/grpc"
)

func main() {
	// TODO 임시
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server := gapi.NewServer(cfg, store)
	grpcServer := grpc.NewServer()

	pb.RegisterMarketServer(grpcServer, server)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("cannot listen network address:", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("start market server at %s", listener.Addr().String())

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("market server failed to serve:", err)
		}
	}()
	wg.Wait()
}
