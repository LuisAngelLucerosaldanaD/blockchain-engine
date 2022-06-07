package main

import (
	"bjungle/blockchain-engine/api"
	"bjungle/blockchain-engine/internal/env"
)

func main() {
	c := env.NewConfiguration()
	api.Start(c.App.Port)
}
