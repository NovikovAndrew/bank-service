package gapi

import (
	db "bank-service/db/sqlc"
	"bank-service/pb"
	"bank-service/token"
	"bank-service/util"
)

type Sever struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}
