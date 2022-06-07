package mine

import (
	"bjungle/blockchain-engine/internal/env"
	"bjungle/blockchain-engine/internal/grpc/mine_proto"
	"bjungle/blockchain-engine/internal/helpers"
	"bjungle/blockchain-engine/internal/logger"
	"bjungle/blockchain-engine/internal/msg"
	"bjungle/blockchain-engine/pkg/bc"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type HandlerMine struct {
	DBMg *mongo.Database
	TxID string
}

// GetBlockToMine godoc
// @Summary blockchain
// @Description GetBlockToMine
// @Accept  json
// @Produce  json
// @Success 200 {object} responseBlockToMine
// @Success 202 {object} dataBlockToMine
// @Router /api/v1/block-to-mine [get]
// @Authorization Bearer token
func (h *HandlerMine) GetBlockToMine(ctx context.Context, request *mine_proto.GetBlockToMineRequest) (*mine_proto.GetBlockToMineResponse, error) {

	e := env.NewConfiguration()

	res := mine_proto.GetBlockToMineResponse{Error: true}

	srvO1 := bc.NewServerBc(h.DBMg, nil, h.TxID)

	bk, err := srvO1.SrvBlocksTmp.GetBlockTwoCommit()
	if err != nil {
		logger.Error.Printf("couldn't get block: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return &res, err
	}

	if bk == nil {
		logger.Error.Printf("couldn't get block: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return &res, fmt.Errorf("couldn't get block")
	}

	//TODO get transactions by block id of api transactions - block mine

	/*trs, err := srvO1.SrvTransactions.GetTransactionsByBlockID(bk.ID)
	if err != nil {
		logger.Error.Printf("couldn't get transactions by block id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}*/

	/*rs, _ := json.Marshal(trs)*/

	hash, err := srvO1.SrvBlocks.GetHashPrevBlock()
	if err != nil {
		logger.Error.Printf("couldn't get hash prev block: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(94, h.DBMg, h.TxID)
		return &res, err
	}

	hashTemp := []byte(hash)
	dataBk := mine_proto.DataBlockMine{
		Id:        bk.ID,
		Timestamp: bk.Timestamp.String(),
		/*Data:       ,*/
		PrevHash:   hashTemp,
		Difficulty: int32(e.App.Difficulty),
	}

	res.Data = &dataBk
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DBMg, h.TxID)
	res.Error = false
	return &res, nil
}

// MineBlock godoc
// @Summary blockchain
// @Description MineBlock
// @Accept  json
// @Produce  json
// @Success 200 {object} requestMineBlock
// @Success 202 {object} dataBlockToMine
// @Router /api/v1/block-to-mine [get]
// @Authorization Bearer token
func (h *HandlerMine) MineBlock(ctx context.Context, block *mine_proto.RequestMineBlock) (*mine_proto.MineBlockResponse, error) {
	res := mine_proto.MineBlockResponse{Error: true}
	u, err := helpers.GetUserContextV2(ctx)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DBMg, h.TxID)
		return &res, err
	}

	srO1 := bc.NewServerBc(h.DBMg, nil, h.TxID)
	bk, code, err := srO1.SrvBlocksTmp.GetBlockTmpByID(block.Id)
	if err != nil {
		logger.Error.Printf("couldn't bind model requestMineBlock: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DBMg, h.TxID)
		return &res, err
	}

	hs, err := srO1.SrvBlocks.GetHashPrevBlock()
	if err != nil {
		logger.Error.Printf("couldn't bind model requestMineBlock: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return &res, err
	}

	// TODO get transactions by block id of api transactions - mine block
	/*ts, err := srO1.SrvTransactions.GetTransactionsByBlockID(m.ID)
	if err != nil {
		logger.Error.Printf("couldn't get transactions by block id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DBMg, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	tsBytes, _ := json.Marshal(ts)*/

	//TODO review minedBY
	_, code, err = srO1.SrvBlocks.CreateBlock(block.Id, "", block.Nonce, int(block.Difficulty), u.ID, time.Now(), bk.Timestamp, block.Hash, hs)
	if err != nil {
		logger.Error.Printf("couldn't CreateBlock: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DBMg, h.TxID)
		return &res, err
	}

	_, code, err = srO1.SrvBlocksTmp.UpdateBlockTmp(block.Id, 3)
	if err != nil {
		logger.Error.Printf("couldn't update status block temp: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DBMg, h.TxID)
		return &res, err
	}

	res.Data = true
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DBMg, h.TxID)
	res.Error = false
	return &res, nil
}
