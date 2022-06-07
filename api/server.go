package api

import (
	"bjungle/blockchain-engine/api/handlers/blocks"
	"bjungle/blockchain-engine/api/handlers/mine"
	"bjungle/blockchain-engine/internal/grpc/blocks_proto"
	"bjungle/blockchain-engine/internal/grpc/mine_proto"
	"bjungle/blockchain-engine/pkg/auth/interceptor"
	"fmt"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	version     = "0.0.2"
	website     = "https://www.bjungle.net"
	banner      = `Blockchain BJungle Engine`
	description = `Blockchain APi Engine - %s - Port: %s
by BJungle 
Version: %s
%s`
)

type server struct {
	listening string
	DB        *mongo.Database
	TxID      string
}

func newServer(listening int, db *mongo.Database, txID string) *server {
	return &server{fmt.Sprintf(":%d", listening), db, txID}
}

func (srv *server) Start() {
	color.Blue(banner)
	color.Cyan(fmt.Sprintf(description, "BLion Egine", srv.listening, version, website))
	lis, err := net.Listen("tcp", "0.0.0.0"+srv.listening)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("Error faltal listener %v", err)
	}
	itr := interceptor.NewAuthInterceptor()
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(itr.Unary()),
		grpc.StreamInterceptor(itr.Stream()),
	}

	s := grpc.NewServer(serverOptions...)

	blocks_proto.RegisterBlockServicesBlocksServer(s, &blocks.HandlerBlocks{DBMg: srv.DB, TxID: srv.TxID})
	mine_proto.RegisterMineBlockServicesBlocksServer(s, &mine.HandlerMine{DBMg: srv.DB, TxID: srv.TxID})

	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Error fatal server", err)
	}
}
