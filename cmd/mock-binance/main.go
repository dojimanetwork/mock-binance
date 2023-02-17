package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.com/thorchain/bepswap/mock-binance/api"
	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

func main() {
	bstore := store.NewInMemory()

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-time.After(300 * time.Millisecond):
				bstore.AddBlock(nil)
			}
		}
	}()

	web := api.API(bstore)

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "26660"
	}

	err := web.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf(err.Error())
	}
}
