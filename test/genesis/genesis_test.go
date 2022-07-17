package main

import (
	"bjungle/blockchain-engine/api/handlers/mine"
	"bjungle/blockchain-engine/internal/db"
	pb "bjungle/blockchain-engine/internal/grpc/mine_proto"
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
	pb.RegisterMineBlockServicesBlocksServer(s, &mine.HandlerMine{DBMg: dbMg, TxID: uuid.New().String()})
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
	md := metadata.New(map[string]string{"Authorization": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTgzMzU2NDAsInVzZXIiOnsiaWQiOiI4Yzc4NzJlZC0yN2U5LTRiYjUtYmQ0NC1mOGFjOTE1ZTQ0ZjYiLCJuaWNrbmFtZSI6Im5leHVtIiwiZW1haWwiOiJyb290Lm5leHVtQG5leHVtLmNvIiwicGFzc3dvcmQiOiIiLCJuYW1lIjoibmV4dW0iLCJsYXN0bmFtZSI6IkRpZ2l0YWwiLCJpZF90eXBlIjowLCJpZF9udW1iZXIiOiI3NTg0MDIzNiIsImNlbGxwaG9uZSI6Iis1MTkyOTQ4MzYxOCIsInN0YXR1c19pZCI6OCwiYmxvY2tfZGF0ZSI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwiZGlzYWJsZWRfZGF0ZSI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwibGFzdF9sb2dpbiI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwibGFzdF9jaGFuZ2VfcGFzc3dvcmQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiIsImJpcnRoX2RhdGUiOiIyMDAwLTAyLTEwVDAwOjAwOjAwWiIsInZlcmlmaWVkX2F0IjoiMjAyMi0wNC0wNFQyMzozMzoyNi41NjNaIiwiaWRfcm9sZSI6MjEsInJlY292ZXJ5X2FjY291bnRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiIsImRlbGV0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiIsImNyZWF0ZWRfYXQiOiIyMDIyLTA3LTE3VDExOjM4OjUxLjk2NTM5WiIsInVwZGF0ZWRfYXQiOiIyMDIyLTA3LTE3VDExOjM4OjUxLjk2NTM5WiJ9fQ.BvdfzxJPtXi1U4llmOMdCRbj5Ai37Hhg89c0T3wW4xmrTwEijKDD0Kwj5k8ng8yp7aNKyxXeVhv8edONw-XfIjj00crb34c58dxDSQXczfNCbUg5QrCUllu3Mk7AwKKF3vTyg16eWXFmC6t8Oo83vDzNq6SIl81oX7Z4httxfGU"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewMineBlockServicesBlocksClient(conn)
	resp, err := client.GenerateBlockGenesis(ctx, &pb.RequestGenerateGenesis{
		KeyGenesis:   "f23a2919-3113-450c-b8cd-ea9d779a1379",
		Password:     "bjungle123",
		Nickname:     "bjunglemain",
		Email:        "souport@bjungle.net",
		TokensEmmit:  10,
		WalletsEmmit: 1,
	})
	if err != nil {
		t.Fatalf("Erro cuando se genero el bloque genesis, error: %v", err)
	}
	log.Printf("Response: %+v", resp)
}
