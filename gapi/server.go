package gapi

import (
	"bank-service/db/sqlc"
	"bank-service/pb"
	"bank-service/token"
	"bank-service/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store, tokenMaker token.Maker) *Server {
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server
}
