package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/rakeshdr543/go-bank-app/gapi"
	pb "github.com/rakeshdr543/go-bank-app/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/rakeshdr543/go-bank-app/api"
	db "github.com/rakeshdr543/go-bank-app/db/sqlc"
	"github.com/rakeshdr543/go-bank-app/util"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	runGrpcServer(config, store)

}

func runGrpcServer(config util.Config, store *db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCAddress)
	if err != nil {
		log.Fatal("cannot listen to grpc server:", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server:", err)
	}

	log.Printf("start grpc server on %s", listener.Addr().String())

}

func runGinServer(config util.Config, store *db.Store) {
	server, err := api.NewServer(config, store)

	err = server.Start(config.HTTPAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
