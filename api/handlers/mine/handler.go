package mine

import (
	"bjungle/blockchain-engine/internal/env"
	"bjungle/blockchain-engine/internal/grpc/accounting_proto"
	"bjungle/blockchain-engine/internal/grpc/mine_proto"
	"bjungle/blockchain-engine/internal/grpc/transactions_proto"
	"bjungle/blockchain-engine/internal/grpc/users_proto"
	"bjungle/blockchain-engine/internal/grpc/wallet_proto"
	"bjungle/blockchain-engine/internal/hash"
	"bjungle/blockchain-engine/internal/helpers"
	"bjungle/blockchain-engine/internal/logger"
	"bjungle/blockchain-engine/internal/msg"
	"bjungle/blockchain-engine/pkg/bc"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcMetadata "google.golang.org/grpc/metadata"
	"strconv"
	"time"
)

type HandlerMine struct {
	DBMg *mongo.Database
	TxID string
}

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
		res.Error = false
		res.Code, res.Type, res.Msg = 29, 1, "No hay bloques disponibles a minar"
		return &res, nil
	}

	hs, err := srvO1.SrvBlocks.GetHashPrevBlock()
	if err != nil {
		logger.Error.Printf("couldn't get hash prev block: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(94, h.DBMg, h.TxID)
		return &res, err
	}

	hashTemp := []byte(hs)
	dataBk := mine_proto.DataBlockMine{
		Id:         bk.ID,
		Timestamp:  bk.Timestamp.String(),
		PrevHash:   hashTemp,
		Difficulty: int32(e.App.Difficulty),
	}

	res.Data = &dataBk
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DBMg, h.TxID)
	res.Error = false
	return &res, nil
}

func (h *HandlerMine) MineBlock(ctx context.Context, block *mine_proto.RequestMineBlock) (*mine_proto.MineBlockResponse, error) {
	res := &mine_proto.MineBlockResponse{Error: true}
	e := env.NewConfiguration()
	u, err := helpers.GetUserContextV2(ctx)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DBMg, h.TxID)
		return res, err
	}

	connTxt, err := grpc.Dial(e.TransactionsService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio de transacciones: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return res, err
	}
	defer connTxt.Close()

	clientTxt := transactions_proto.NewTransactionsServicesClient(connTxt)

	token, err := helpers.GetTokenFromContext(ctx)
	if err != nil {
		logger.Error.Printf("error de authenticación: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return res, err
	}

	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", token)

	srO1 := bc.NewServerBc(h.DBMg, nil, h.TxID)
	bk, code, err := srO1.SrvBlocksTmp.GetBlockTmpByID(block.Id)
	if err != nil {
		logger.Error.Printf("couldn't bind model requestMineBlock: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DBMg, h.TxID)
		return res, err
	}

	hs, err := srO1.SrvBlocks.GetHashPrevBlock()
	if err != nil {
		logger.Error.Printf("couldn't bind model requestMineBlock: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return res, err
	}

	resTxt, err := clientTxt.GetTransactionsByBlockId(ctx, &transactions_proto.RqGetTransactionByBlock{BlockId: block.Id})
	if err != nil {
		logger.Error.Printf("couldn't get transactions by block id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DBMg, h.TxID)
		return res, err
	}

	if resTxt == nil {
		logger.Error.Printf("couldn't get transactions by block id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DBMg, h.TxID)
		return res, fmt.Errorf("couldn't get transactions by block id")
	}

	if resTxt.Error {
		logger.Error.Printf(resTxt.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DBMg, h.TxID)
		return res, fmt.Errorf(resTxt.Msg)
	}

	tsBytes, _ := json.Marshal(resTxt.Data)

	_, code, err = srO1.SrvBlocks.CreateBlock(block.Id, string(tsBytes), block.Nonce, int(block.Difficulty), u.ID, time.Now(), bk.Timestamp, block.Hash, hs)
	if err != nil {
		logger.Error.Printf("couldn't CreateBlock: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DBMg, h.TxID)
		return res, err
	}

	_, code, err = srO1.SrvBlocksTmp.UpdateBlockTmp(block.Id, 3)
	if err != nil {
		logger.Error.Printf("couldn't update status block temp: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DBMg, h.TxID)
		return res, err
	}

	res.Data = true
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DBMg, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerMine) GenerateBlockGenesis(ctx context.Context, request *mine_proto.RequestGenerateGenesis) (*mine_proto.ResponseGenerateGenesis, error) {
	res := &mine_proto.ResponseGenerateGenesis{Error: true}
	e := env.NewConfiguration()
	srvBc := bc.NewServerBc(h.DBMg, nil, h.TxID)

	if request.KeyGenesis != e.App.KeyGenesis {
		logger.Error.Printf("key genesis is invalid")
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DBMg, h.TxID)
		return res, fmt.Errorf("key genesis is invalid")
	}

	if srvBc.SrvBlocks.ExistsBlocks() {
		logger.Error.Printf("block already exists")
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DBMg, h.TxID)
		return res, fmt.Errorf("block already exists")
	}

	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return res, err
	}
	defer connAuth.Close()

	connTxt, err := grpc.Dial(e.TransactionsService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio de transacciones: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return res, err
	}
	defer connTxt.Close()

	clientWallet := wallet_proto.NewWalletServicesWalletClient(connAuth)
	clientUser := users_proto.NewAuthServicesUsersClient(connAuth)
	clientAccount := accounting_proto.NewAccountingServicesAccountingClient(connAuth)
	clientTxt := transactions_proto.NewTransactionsServicesClient(connTxt)

	token, err := helpers.GetTokenFromContext(ctx)
	if err != nil {
		logger.Error.Printf("error de authenticación: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return res, err
	}

	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", token)

	resUSer, err := clientUser.CreateUserBySystem(ctx, &users_proto.RequestCreateUserBySystem{
		Nickname:      request.Nickname,
		Email:         request.Email,
		Password:      request.Password,
		FullPathPhoto: "",
		Name:          "BJungle",
		Lastname:      "BJungle",
		IdType:        6,
		IdNumber:      "75840278",
		Cellphone:     "+57310000000",
		BirthDate:     time.Now().Format("2006-01-02T15:04:05.000Z"),
	})
	if err != nil {
		logger.Error.Printf("error creando usuario: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return res, err
	}

	if resUSer == nil {
		logger.Error.Printf("error creando usuario: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return res, err
	}

	if resUSer.Error {
		logger.Error.Printf(resUSer.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
		return res, fmt.Errorf(resUSer.Msg)
	}

	user := resUSer.Data

	bkTemp, code, err := srvBc.SrvBlocksTmp.CreateBlockTmp(1, time.Now())
	if err != nil {
		logger.Error.Printf("couldn't create block tmp: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DBMg, h.TxID)
		return res, err
	}

	var walletsMains []*mine_proto.WalletMain

	for i := 0; i < int(request.WalletsEmmit); i++ {
		resWallet, err := clientWallet.CreateWalletBySystem(ctx, &wallet_proto.RqCreateWalletBySystem{IdentityNumber: user.IdNumber})
		if err != nil {
			logger.Error.Printf("error creando wallet: %s", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, err
		}

		if resWallet == nil {
			logger.Error.Printf("error creando wallet: %s", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, fmt.Errorf("error creando wallet")
		}

		if resWallet.Error {
			logger.Error.Printf(resWallet.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, fmt.Errorf(resWallet.Msg)
		}

		wallet := resWallet.Data

		resAccount, err := clientAccount.CreateAccounting(ctx, &accounting_proto.RequestCreateAccounting{
			Id:       uuid.New().String(),
			IdWallet: wallet.Id,
			Amount:   0,
			IdUser:   user.Id,
		})
		if err != nil {
			logger.Error.Printf("error creando cuenta: %s", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, err
		}

		if resAccount == nil {
			logger.Error.Printf("error creando cuenta: %s", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, fmt.Errorf("error creando cuenta")
		}

		if resAccount.Error {
			logger.Error.Printf(resAccount.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, fmt.Errorf(resAccount.Msg)
		}

		resUserWallet, err := clientUser.CreateUserWallet(ctx, &users_proto.RqCreateUserWallet{
			UserId:   user.Id,
			WalletId: wallet.Id,
		})
		if err != nil {
			logger.Error.Printf("couldn't create user wallet: %s", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, err
		}

		if resUserWallet == nil {
			logger.Error.Printf("couldn't create user wallet")
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, fmt.Errorf("couldn't create user wallet")
		}

		if resUserWallet.Error {
			logger.Error.Printf(resUserWallet.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, fmt.Errorf(resUserWallet.Msg)
		}

		resTxt, err := clientTxt.CreateTransactionBySystem(ctx, &transactions_proto.RqCreateTransactionBySystem{
			WalletFrom: request.KeyGenesis,
			WalletTo:   wallet.Id,
			Amount:     request.TokensEmmit,
			TypeId:     18,
			Data:       fmt.Sprintf(dataJson(), wallet.Id, request.TokensEmmit),
			BlockId:    bkTemp.ID,
		})
		if err != nil {
			logger.Error.Printf("error creando transaccion: %s", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, err
		}

		if resTxt == nil {
			logger.Error.Printf("error creando transaccion: %s", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, fmt.Errorf("error creando transaccion")
		}

		if resTxt.Error {
			logger.Error.Printf(resTxt.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, fmt.Errorf(resTxt.Msg)
		}

		resAmount, err := clientAccount.SetAmountToAccounting(ctx, &accounting_proto.RequestSetAmountToAccounting{
			WalletId: wallet.Id,
			Amount:   request.TokensEmmit,
			IdUser:   user.Id,
		})
		if err != nil {
			logger.Error.Printf("error asiganando los acais a la cuenta: %s", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, err
		}

		if resAmount == nil {
			logger.Error.Printf("error asiganando los acais a la cuenta: %s", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, fmt.Errorf("error asiganando los acais a la cuenta")
		}

		if resAmount.Error {
			logger.Error.Printf(resAmount.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DBMg, h.TxID)
			return res, fmt.Errorf(resAmount.Msg)
		}

		walletsMains = append(walletsMains, &mine_proto.WalletMain{Id: wallet.Id, Mnemonic: wallet.Mnemonic})
	}

	_, code, err = srvBc.SrvBlocksTmp.UpdateBlockTmp(bkTemp.ID, 2)
	if err != nil {
		logger.Error.Printf("couldn't update status block temp: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DBMg, h.TxID)
		return res, err
	}

	resTxt, err := clientTxt.GetTransactionsByBlockId(ctx, &transactions_proto.RqGetTransactionByBlock{BlockId: bkTemp.ID})
	if err != nil {
		logger.Error.Printf("couldn't get transactions by block id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DBMg, h.TxID)
		return res, err
	}

	if resTxt == nil {
		logger.Error.Printf("couldn't get transactions by block id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DBMg, h.TxID)
		return res, fmt.Errorf("couldn't get transactions by block id")
	}

	if resTxt.Error {
		logger.Error.Printf(resTxt.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DBMg, h.TxID)
		return res, fmt.Errorf(resTxt.Msg)
	}

	tsBytes, _ := json.Marshal(resTxt.Data)

	timeStamp := []byte(strconv.FormatInt(bkTemp.Timestamp.Unix(), 10))
	hs, nonce, err := hash.GenerateHashToMineBlock(timeStamp, tsBytes, nil, e.App.Difficulty)
	if err != nil {
		logger.Error.Printf("couldn't get transactions by block id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DBMg, h.TxID)
		return res, err
	}
	_, code, err = srvBc.SrvBlocks.CreateBlock(bkTemp.ID, string(tsBytes), int64(nonce), e.App.Difficulty, user.Id, time.Now(), bkTemp.Timestamp, hs, "genesis")
	if err != nil {
		logger.Error.Printf("couldn't CreateBlock: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DBMg, h.TxID)
		return res, err
	}

	_, code, err = srvBc.SrvBlocksTmp.UpdateBlockTmp(bkTemp.ID, 3)
	if err != nil {
		logger.Error.Printf("couldn't update status block temp: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DBMg, h.TxID)
		return res, err
	}

	res.Error = false
	res.Data = &mine_proto.Data{
		UserId:      user.Id,
		WalletsMain: walletsMains,
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DBMg, h.TxID)
	return res, nil
}

func dataJson() string {
	return `{
        "files": [],
        "name": "Genesis",
        "description": "Emmit tokes to main wallet",
        "entities": [
            {
                "name":   "Tokens Emmit",
                "attributes": [
                    {
                        "name": "walletID",
                        "value": "%s"
                    },
                    {
                        "name": "tokens",
                        "value": "%f"
                    }
                ]
            }
        ]
    }`
}
