package main

import (
	"bjungle/blockchain-engine/api/handlers/blocks"
	"bjungle/blockchain-engine/internal/db"
	pb "bjungle/blockchain-engine/internal/grpc/blocks_proto"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {

	dbMg := db.GetMgConnection()
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterBlockServicesBlocksServer(s, &blocks.HandlerBlocks{DBMg: dbMg, TxID: uuid.New().String()})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGetAllBlocks(t *testing.T) {
	md := metadata.New(map[string]string{"Authorization": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTQyMjAxMDEsInVzZXIiOnsiaWQiOiI4Yzc4NzJlZC0yN2U5LTRiYjUtYmQ0NC1mOGFjOTE1ZTQ0ZjYiLCJuaWNrbmFtZSI6Im5leHVtIiwiZW1haWwiOiJyb290Lm5leHVtQG5leHVtLmNvIiwicGFzc3dvcmQiOiIiLCJuYW1lIjoibmV4dW0iLCJsYXN0bmFtZSI6IkRpZ2l0YWwiLCJpZF90eXBlIjowLCJpZF9udW1iZXIiOiI3NTg0MDIzNiIsImNlbGxwaG9uZSI6Iis1MTkyOTQ4MzYxOCIsInN0YXR1c19pZCI6OCwiYmxvY2tfZGF0ZSI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwiZGlzYWJsZWRfZGF0ZSI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwibGFzdF9sb2dpbiI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwibGFzdF9jaGFuZ2VfcGFzc3dvcmQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiIsImJpcnRoX2RhdGUiOiIyMDAwLTAyLTEwVDAwOjAwOjAwWiIsInZlcmlmaWVkX2F0IjoiMjAyMi0wNC0wNFQyMzozMzoyNi41NjM4NjFaIiwiaWRfcm9sZSI6MjEsInJlY292ZXJ5X2FjY291bnRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiIsImRlbGV0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiIsImNyZWF0ZWRfYXQiOiIyMDIyLTA0LTA0VDIzOjMzOjI2LjU2Mzg2MVoiLCJ1cGRhdGVkX2F0IjoiMjAyMi0wNS0yNFQxOToyMDo1OS42Njg1MzVaIn19.QI3iGlwYGbfCn_xd1N3HAQueH4C1l_dC1zF-2cxAhnAz-U6i2e0lhEG3bYve6ZVp1vtUe_CFBcel7FPCIrEwpc43jhUFgi_5ShhkkPqOzCeJ6OKCOSZ0xcLBMVCOulUX8o9iV-dO_7FboMN-yrWHWBBsv3I6FbL819vA8KH-mF4"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewBlockServicesBlocksClient(conn)
	resp, err := client.GetBlock(ctx, &pb.GetAllBlockRequest{
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("Get all blocks failed, error: %v", err)
	}

	log.Printf("Response: %v", resp)
}
