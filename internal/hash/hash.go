package hash

import (
	"bjungle/blockchain-engine/internal/helpers"
	"bytes"
	"crypto/sha256"
	"fmt"

	"strings"
)

func GenerateHashToMineBlock(timestamp, data, prevBlockHash []byte, difficulty int) (string, int, error) {
	var nonce int
	var Zeros string
	for i := 1; i <= difficulty; i++ {
		Zeros = Zeros + "0"
	}
	for {
		headers := bytes.Join(
			[][]byte{prevBlockHash,
				data,
				timestamp,
				helpers.ToHex(int64(nonce)),
				helpers.ToHex(int64(difficulty))}, []byte{})
		hashB := sha256.Sum256(headers)
		hash := fmt.Sprintf("%x", hashB)
		if strings.HasPrefix(hash, Zeros) {
			return hash, nonce, nil
		} else {
			nonce++
		}
	}
}
