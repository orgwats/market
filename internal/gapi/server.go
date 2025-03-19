package gapi

import (
	"github.com/orgwats/stream/internal/config"
	db "github.com/orgwats/stream/internal/db/sqlc"
	"github.com/orgwats/stream/internal/hub"
	"github.com/orgwats/stream/internal/stream"

	pb "github.com/orgwats/idl/gen/go/market"
)

type Server struct {
	pb.UnimplementedMarketServer

	// TODO: 임시
	cfg   *config.Config
	store db.Store

	hub *hub.Hub
}

func NewServer(cfg *config.Config, store db.Store) *Server {
	ch := stream.NewStream(cfg)
	hub := hub.NewHub()

	go func() {
		for c := range ch {
			hub.Broadcast(c.Symbol, c)
		}
	}()

	return &Server{
		cfg:   cfg,
		store: store,
		hub:   hub,
	}
}
