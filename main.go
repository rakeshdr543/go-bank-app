package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakeshdr543/go-bank-app/gapi"
	pb "github.com/rakeshdr543/go-bank-app/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

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

	go runGatewayServer(config, store)
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

func runGatewayServer(config util.Config, store *db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	jsonOptions :=
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		})

	grpcMux := runtime.NewServeMux(jsonOptions)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPAddress)
	if err != nil {
		log.Fatal("cannot listen to grpc server:", err)
	}

	err = http.Serve(listener, mux)
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
