package block_id

import (
	"bjungle/blockchain-engine/api/handlers/blocks"
	"bjungle/blockchain-engine/internal/db"
	pb "bjungle/blockchain-engine/internal/grpc/blocks_proto"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
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

func TestGetBlockById(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewBlockServicesBlocksClient(conn)
	resp, err := client.GetBlockByID(ctx, &pb.GetByIdRequest{Id: 12})
	if err != nil {
		t.Fatalf("Get block by id failed, error: %v", err)
	}
	log.Printf("Response: %+v", resp)
}
