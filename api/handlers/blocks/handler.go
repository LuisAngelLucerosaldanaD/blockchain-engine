package blocks

import (
	"bjungle/blockchain-engine/internal/grpc/blocks_proto"
	"bjungle/blockchain-engine/internal/logger"
	"bjungle/blockchain-engine/internal/msg"
	"bjungle/blockchain-engine/pkg/bc"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type HandlerBlocks struct {
	DBMg *mongo.Database
	TxID string
}

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

func (h *HandlerBlocks) GetBlockUnCommit(ctx context.Context, commit *blocks_proto.RequestGetBlockUnCommit) (*blocks_proto.ResponseGetBlockUnCommit, error) {
	res := &blocks_proto.ResponseGetBlockUnCommit{Error: true, Data: nil}
	srvBc := bc.NewServerBc(h.DBMg, nil, h.TxID)

	bks, err := srvBc.SrvBlocksTmp.GetBlockUnCommit()
	if err != nil {
		logger.Error.Printf("couldn't get block un commit: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DBMg, h.TxID)
		return res, err
	}

	if bks != nil {
		res.Data = &blocks_proto.BlockTemp{
			Id:        bks.ID,
			Status:    int32(bks.Status),
			Timestamp: bks.Timestamp.String(),
			CreatedAt: bks.CreatedAt.String(),
			UpdatedAt: bks.UpdatedAt.String(),
		}
	}

	res.Error = false
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DBMg, h.TxID)
	return res, nil
}

func (h *HandlerBlocks) CreateBlockTemp(ctx context.Context, blockTemp *blocks_proto.RequestCreateBlockTemp) (*blocks_proto.ResponseCreateBlockTemp, error) {
	res := &blocks_proto.ResponseCreateBlockTemp{Error: true}
	srvBc := bc.NewServerBc(h.DBMg, nil, h.TxID)

	newBk, code, err := srvBc.SrvBlocksTmp.CreateBlockTmp(1, time.Now())
	if err != nil {
		logger.Error.Printf("couldn't created block: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DBMg, h.TxID)
		return res, err
	}

	res.Error = false
	res.Data = &blocks_proto.BlockTemp{
		Id:        newBk.ID,
		Status:    int32(newBk.Status),
		Timestamp: newBk.Timestamp.String(),
		CreatedAt: newBk.CreatedAt.String(),
		UpdatedAt: newBk.UpdatedAt.String(),
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DBMg, h.TxID)
	return res, nil
}

func (h *HandlerBlocks) UpdateBlockTemp(ctx context.Context, blockTemp *blocks_proto.RequestUpdateBlockTemp) (*blocks_proto.ResponseUpdateBlockTemp, error) {
	res := &blocks_proto.ResponseUpdateBlockTemp{Error: true}
	srvBc := bc.NewServerBc(h.DBMg, nil, h.TxID)

	bkTemp, code, err := srvBc.SrvBlocksTmp.UpdateBlockTmp(blockTemp.Id, int(blockTemp.Status))
	if err != nil {
		logger.Error.Printf("couldn't close block: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DBMg, h.TxID)
		return res, err
	}

	res.Error = false
	res.Data = &blocks_proto.BlockTemp{
		Id:        bkTemp.ID,
		Status:    int32(bkTemp.Status),
		Timestamp: bkTemp.Timestamp.String(),
		CreatedAt: bkTemp.CreatedAt.String(),
		UpdatedAt: bkTemp.UpdatedAt.String(),
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DBMg, h.TxID)
	return res, nil
}
