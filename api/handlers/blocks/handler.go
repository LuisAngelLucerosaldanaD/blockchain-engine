package blocks

import (
	"bjungle/blockchain-engine/internal/grpc/blocks_proto"
	"bjungle/blockchain-engine/internal/logger"
	"bjungle/blockchain-engine/internal/msg"
	"bjungle/blockchain-engine/pkg/bc"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type HandlerBlocks struct {
	DBMg *mongo.Database
	TxID string
}

// GetBlock godoc
// @Summary GetBlock BJungle
// @Description GetBlock BJungle
// @Accept  json
// @Produce  json
// @Success 200 {object} responseGetAllBlock
// @Success 202 {object} responseGetAllBlock
// @Router /api/v1/blocks [get]
func (h *HandlerBlocks) GetBlock(ctx context.Context, request *blocks_proto.GetAllBlockRequest) (*blocks_proto.GetAllBlockResponse, error) {
	res := blocks_proto.GetAllBlockResponse{Error: true}
	srvO1 := bc.NewServerBc(h.DBMg, nil, h.TxID)

	bks, err := srvO1.SrvBlocks.GetAllBlock(&request.Limit, &request.Offset)
	if err != nil {
		logger.Error.Printf("couldn't get all blocks: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return &res, err
	}

	var bksPb []*blocks_proto.DataBlock

	for _, bk := range bks {
		bkTemp := blocks_proto.DataBlock{
			Id:                 bk.ID,
			Data:               bk.Data,
			Nonce:              bk.Nonce,
			Difficulty:         int32(bk.Difficulty),
			MinedBy:            bk.MinedBy,
			MinedAt:            bk.MinedAt.String(),
			Timestamp:          bk.Timestamp.String(),
			Hash:               bk.Hash,
			PrevHash:           bk.PrevHash,
			StatusId:           int32(bk.StatusId),
			IdUser:             bk.IdUser,
			LastValidationDate: bk.LastValidationDate,
			CreatedAt:          bk.CreatedAt.String(),
			UpdatedAt:          bk.UpdatedAt.String(),
		}
		bksPb = append(bksPb, &bkTemp)
	}

	res.Data = bksPb
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DBMg, h.TxID)
	res.Error = false
	return &res, nil
}

// GetBlockByID godoc
// @Summary GetBlockByID BJungle
// @Description GetBlockByID BJungle
// @Accept  json
// @Produce  json
// @Success 200 {object} responseGetBlock
// @Success 202 {object} responseGetBlock
// @Router /api/v1/block [get]
func (h *HandlerBlocks) GetBlockByID(ctx context.Context, request *blocks_proto.GetByIdRequest) (*blocks_proto.GetBlockByIDResponse, error) {
	res := blocks_proto.GetBlockByIDResponse{Error: true}

	srvO1 := bc.NewServerBc(h.DBMg, nil, h.TxID)

	bks, err := srvO1.SrvBlocks.GetBlocksById(request.Id)
	if err != nil {
		logger.Error.Printf("couldn't get blocks by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DBMg, h.TxID)
		return &res, err
	}

	blockRes := blocks_proto.DataBlock{
		Id:                 bks.ID,
		Data:               bks.Data,
		Nonce:              bks.Nonce,
		Difficulty:         int32(bks.Difficulty),
		MinedBy:            bks.MinedBy,
		MinedAt:            bks.MinedAt.String(),
		Timestamp:          bks.Timestamp.String(),
		Hash:               bks.Hash,
		PrevHash:           bks.PrevHash,
		StatusId:           int32(bks.StatusId),
		IdUser:             bks.IdUser,
		LastValidationDate: bks.LastValidationDate,
		CreatedAt:          bks.CreatedAt.String(),
		UpdatedAt:          bks.UpdatedAt.String(),
	}

	res.Data = &blockRes
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DBMg, h.TxID)
	res.Error = false
	return &res, nil
}
