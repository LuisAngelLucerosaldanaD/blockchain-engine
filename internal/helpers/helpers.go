package helpers

import (
	"bjungle/blockchain-engine/internal/env"
	"bjungle/blockchain-engine/internal/logger"
	"bjungle/blockchain-engine/internal/models"
	"bytes"
	"crypto/rsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"log"

	"context"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var (
	signKey   *rsa.PublicKey
	publicKey string
)

type UserClaims struct {
	jwt.StandardClaims
	User models.User `json:"user"`
	Role int         `json:"role"`
}

// init lee los archivos de firma y validación RSA
func init() {
	c := env.NewConfiguration()
	publicKey = c.App.RSAPublicKey
	keyBytes, err := ioutil.ReadFile(publicKey)
	if err != nil {
		logger.Error.Printf("leyendo el archivo privado de firma: %s", err)
	}

	signKey, err = jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		logger.Error.Printf("realizando el parse en auth RSA private: %s", err)
	}
}

func GetUserContext(c *fiber.Ctx) *models.User {
	userUrl := c.Locals("user").(*jwt.Token)
	claims := userUrl.Claims.(jwt.MapClaims)
	for i, cl := range claims {
		if i == "user" {
			u := models.User{}
			ub, _ := json.Marshal(cl)
			_ = json.Unmarshal(ub, &u)
			return &u
		}
	}
	return nil
}

func GetUserContextV2(ctx context.Context) (*models.User, error) {
	tokenStr, err := getTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})
	if err != nil {
		return nil, err
	}

	for i, cl := range claims {
		if i == "user" {
			u := models.User{}
			ub, _ := json.Marshal(cl)
			_ = json.Unmarshal(ub, &u)
			return &u, nil
		}
	}

	return nil, nil
}

func getTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("ErrNoMetadataInContext")
	}

	token, ok := md["authorization"]
	if !ok || len(token) == 0 {
		return "", fmt.Errorf("ErrNoAuthorizationInMetadata")
	}

	return token[0], nil
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)

	}
	return buff.Bytes()
}
