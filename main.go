package main

import (
	"bank-service/api"
	"bank-service/db/sqlc"
	"bank-service/gapi"
	"bank-service/pb"
	"bank-service/token"
	"bank-service/util"
	"database/sql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	pasetoMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		log.Fatal(err)
	}

	startCh := make(chan struct{})

	go func() {
		runGinServer(config, store, pasetoMaker)
		startCh <- struct{}{}
	}()

	go func() {
		runGrpcServer(config, store, pasetoMaker)
		startCh <- struct{}{}
	}()

	<-startCh
}

func runGinServer(config util.Config, store db.Store, tokenMaker token.Maker) {
	server := api.NewServer(store, tokenMaker, config)

	if err := server.Start(config.HTTPServerAddress); err != nil {
		log.Fatal(err)
	}
}

func runGrpcServer(config util.Config, store db.Store, tokenMaker token.Maker) {
	server := gapi.NewServer(config, store, tokenMaker)
	grpcServer := grpc.NewServer()

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	tcpListenr, err := net.Listen("tcp", config.GRPCServerAddress)

	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s address", config.GRPCServerAddress)
	err = grpcServer.Serve(tcpListenr)

	if err != nil {
		log.Fatal("cannot to serve gRpc server")
	}
}
