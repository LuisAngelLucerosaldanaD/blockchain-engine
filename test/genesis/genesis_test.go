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
	md := metadata.New(map[string]string{"Authorization": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU0MTU5MzEsInJvbCI6ImFkbWluIiwidXNlciI6eyJpZCI6ImQ1ZGIwMGZjLTYxMzktNDljYS1iMmQwLWNhZjUyYzI3ZjZiOCIsIm5pY2tuYW1lIjoiYmxpb24iLCJlbWFpbCI6InJvb3RAYmp1bmdsZS5uZXQiLCJwYXNzd29yZCI6IiQyYSQxMCRHUlVYWTFNMEQxL1FLOTQuQ1pDbnQuT3k1aTdkbk9WY20wNGRlTzhXR3l5UEtYd1k4SjBkYSIsImZpcnN0X25hbWUiOiJCTGlvbiIsInNlY29uZF9uYW1lIjoiQkxpb24iLCJmaXJzdF9zdXJuYW1lIjoiQkxpb24iLCJzZWNvbmRfc3VybmFtZSI6IkJMaW9uIiwiYWdlIjoyLCJ0eXBlX2RvY3VtZW50IjoiTklUIiwiZG9jdW1lbnRfbnVtYmVyIjoiOTg3NjU0MzIxIiwiY2VsbHBob25lIjoiOTIzMDYyNzQ5IiwiZ2VuZGVyIjoiTWFzY3VsaW5vIiwibmF0aW9uYWxpdHkiOiJDb2xvbWJpYSIsImNvdW50cnkiOiJDb2xvbWJpYSIsImRlcGFydG1lbnQiOiJCb2dvdGEiLCJjaXR5IjoiQm9nb3RhIERDIiwicmVhbF9pcCI6IjE5Mi4xNjguMTIuMTciLCJzdGF0dXNfaWQiOjEsImZhaWxlZF9hdHRlbXB0cyI6MCwiYmxvY2tfZGF0ZSI6bnVsbCwiZGlzYWJsZWRfZGF0ZSI6bnVsbCwibGFzdF9sb2dpbiI6bnVsbCwibGFzdF9jaGFuZ2VfcGFzc3dvcmQiOiIyMDIzLTA5LTE5VDE1OjExOjE5WiIsImJpcnRoX2RhdGUiOiIyMDAwLTAyLTEwVDE1OjExOjI4WiIsInZlcmlmaWVkX2NvZGUiOm51bGwsImlzX2RlbGV0ZWQiOmZhbHNlLCJkZWxldGVkX2F0IjpudWxsLCJjcmVhdGVkX2F0IjoiMjAyMy0wOS0xOVQyMDoxMjoxMS42OTgxNjlaIiwidXBkYXRlZF9hdCI6IjIwMjMtMDktMTlUMjA6MTI6MTEuNjk4MTY5WiJ9fQ.d4WXp7hv5qzCxaLeJF9qR6ILMndKmup5E8PpYWc9p6sLo3xKmA_5LAmpebdQiefn3EYU5TM-PS0Fu3ylzDMmwuHBrikWs8Gt5bDKPlMxIQuMGYxoDTMJ5vPGdeOtotQp23NPm7gVXYYpfOtXmpsNDtIji_g8gOk1TCVpTZPRYjHGjSHl96dHIZgtoW6cHWYkkOATlvgV7i6RM-Rky_od5ggP-LiFloJSeQTibh5ApU1I473Jhd05bAgv2BKF7Cu3rAZnXe24uB4pSHV_2dxOWXtFZBkkSBkVQB4FEbnXutkgiPUOEXH9Ao3OLKtFdx-PSk1FRA-cjoCC5qsZmBUEmg"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewMineBlockServicesBlocksClient(conn)
	resp, err := client.GenerateBlockGenesis(ctx, &pb.RequestGenerateGenesis{
		KeyGenesis:     "f23a2919-3113-450c-b8cd-ea9d779a1379",
		TokensEmmit:    10,
		WalletsEmmit:   1,
		UserId:         "d5db00fc-6139-49ca-b2d0-caf52c27f6b8",
		IdentityNumber: "987654321",
	})
	if err != nil {
		t.Fatalf("Erro cuando se genero el bloque genesis, error: %v", err)
	}
	log.Printf("Response: %+v", resp)
}
