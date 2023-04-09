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
	md := metadata.New(map[string]string{"Authorization": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoiMjVjNTlhM2ItN2YzZi00ZGNkLWI3ZmQtYTRhNzA3MWQ1N2YxIiwidXNlcm5hbWUiOiJtYW5hZ2VyIiwibmFtZSI6Im1hbmFnZXIiLCJsYXN0bmFtZSI6Im1hbmFnZXIiLCJlbWFpbF9ub3RpZmljYXRpb25zIjoibG9yZW5hLnZhcmdhc0BianVuZ2xlLm5ldCIsImlkZW50aWZpY2F0aW9uX251bWJlciI6IjgwMTg4Mjg0IiwiaWRlbnRpZmljYXRpb25fdHlwZSI6IkMiLCJzdGF0dXMiOjEsImxhc3RfY2hhbmdlX3Bhc3N3b3JkIjoiMjAyMi0wNy0xMVQwMDowMDowMFoiLCJibG9ja19kYXRlIjoiMjAyMi0wNC0zMFQyMTo0NjoyNC41MjU5OTVaIiwiY2hhbmdlX3Bhc3N3b3JkIjpmYWxzZSwiaXNfYmxvY2siOmZhbHNlLCJpc19kaXNhYmxlZCI6ZmFsc2UsImNsaWVudF9pZCI6OTkyNiwicmVhbF9pcCI6IjEzMi4xNTcuMTI4LjU4Iiwic2Vzc2lvbl9pZCI6ImIzOTk1MTEwLTQ4NDItNGIxMS1iZmI5LTI0ZjNhOTQ3MGFmMyIsImNvbG9ycyI6eyJwcmltYXJ5IjoiIzM1M0E0OCIsInNlY29uZGFyeSI6IiMwMzliZTUiLCJ0ZXJ0aWFyeSI6IiMyNjI5MzMifSwicm9sZXMiOlsiMzg4NTYyNzQtNTMxMS00MTczLTk3ZDEtZmI3ZTJjOGM5ZTE2IiwiOWM3ZjY3ODctZWEzMC00MmZhLTliMDItYmIwZDZjYTI0OTQ1IiwiYWYwYTNmM2YtNmRjZC00ZTQ4LWE1YjMtNjdhYjhkODQyNjNiIiwiODJhOGM2YjYtOGM3Ni00YWE4LWFlNDYtNmViODI4MjhiODFiIiwiZWM2NzAzMTYtOTg5Zi00M2QyLWE3OWYtN2FiOWE4ZmMyYjc5IiwiODVhZDVkYmMtODcwYi00NWI2LWI3MTEtZjIzYTZhNDJjMmNkIiwiOWMzNmI2MTktYzIyZC00MWY4LWFjM2ItZjg5NWExZTEwMTMwIiwiMmI2YjA4YjAtNWFiYS00MTE2LWFlNTAtODhmYTM3ZWU5ZjdhIiwiZGM2Y2MxMmYtMDFlNy00MDFhLTg3MGUtNGI0ZTQ4MTJlNzA5IiwiOTliNGU3YjItMTYxZi00ZDMwLWExOWItM2M2MGI5Y2QxN2MyIiwiYzliMDYzZTMtNzIyYi00NDRjLWE3MTktYThiN2ViNTkwMGY1IiwiNmQ2OThmZGUtMGY0OC00M2ZmLTkxYTMtMjczNzg5MmIzN2JiIiwiYzlmOTcxOWMtNGEzNS00MWIwLWE2MmMtYzMxN2Y0Yzg1NmRiIiwiYjc3OGVkYTktNGRjZi00OTkzLTg0ZTktNjU3MGIzMjlmZThkIiwiYWFjODllYTctMmViYi00ZDgxLTlkOTAtNDQzOWVjNTA1YzIxIiwiOTEwNDgxNjktZWM0Yi00Zjg2LTk4YjMtMjk5ZTUwMWIxYTQwIiwiZTNlZGYyMGItZTE2MS00YTE4LThlYzItZDJjZmUwNGFhZDA1IiwiMDJiY2E5ODItNWQ4NS00YjBhLWI1ZDEtYTc1NGE5ZTc2ZDBkIiwiMDlkZTIyYWUtNDMwNS00MmRkLTkyOTItNTIyOGMyZDFiMWE2IiwiY2IyOTljOGEtYWQxZS00YjVlLTgxZTItZGNhZmVkYzlkMGQyIiwiZTg3ZjA4ZmUtZTE1My00ZWE1LTkzNWMtNDc3MzY3YWNhZTQ2IiwiOTY0NTM4NGMtMjllNy00YjFiLTg5NDgtNWI4MzdhZDQ2ZWVjIiwiZGZkYjZmOGEtOGFmYi00M2NjLWE3Y2QtNmY0NzY2YThiODI3Il0sInByb2plY3RzIjpbIjU0N2ZkZDQyLWQ1YWItNGM3OS04ZWU4LTlmYTVhNjMxYzA2NCIsImRjZjg4Yzc2LTFhYjItNGM5ZS05ZDQ3LTJhZTUyMjkyMTkyYiJdLCJjcmVhdGVkX2F0IjoiMjAyMC0wOC0xN1QwOTowOToxMC45ODdaIiwidXBkYXRlZF9hdCI6IjIwMjItMTEtMTRUMDc6MDQ6NTMuNjEwN1oifSwiaXBfYWRkcmVzcyI6IjEzMi4xNTcuMTI4LjU4IiwiZXhwIjoxNjcyNDI2MjAwLCJpc3MiOiJFY2F0Y2gtQlBNIn0.HxvyzonNfjw_rBL3ntEz4DvVV90uiaubF5tjgFxWLAjGJnWMFc1t5ZnAuiDBmRNu8nX0UMA2C10ACtTDIP5FQX0Xs6besTGgjqX8iCl6xmRcjd9ug115d0ywEAAyDjPBTUsnztcPOsg0pMeV2gzQQdvxXNBfJu2Cav9z1G7vuUQ"})
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
