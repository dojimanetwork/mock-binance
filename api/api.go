package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"log"
	"sync"

	"github.com/gin-gonic/gin"

	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

func API(bstore store.Store) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
	}))

	// this mutex is used to ensure we don't have concurrent writes
	// (broadcasts) to the chain at the same time (give bad seq num behavior)
	mutex := &sync.Mutex{}

	// healthcheck
	r.GET("/ping", ping())

	// Get a block
	r.GET("/block", block(bstore))

	// ABCIinfo
	r.GET("/abci_info", ABCIinfo(bstore))
	r.GET("/abci_query", ABCIquery(bstore))

	// Status
	r.GET("/status", status())

	// search for tx
	r.GET("/tx_search", txSearch(bstore))
	r.GET("/blocks", list_blocks(bstore))

	// Broadcast endpoints
	r.POST("/broadcast/easy", broadcastEasy(bstore, mutex)) // easier interface to broadcast a transaction
	r.POST("/broadcast", broadcast(bstore, mutex))
	r.POST("/broadcast_tx_commit", broadcast(bstore, mutex))
	r.POST("/broadcast_tx_async", broadcast(bstore, mutex))

	// Account information
	// Smoke tests only use this endpoint for the acc number and seq number.
	// And since we don't care about that for the mock binance service, we can
	// just return any numbers.
	r.GET("/account/:acc", acc_info(bstore))
	r.GET("/accounts", acc_list(bstore))

	// setup default not found message
	r.NoRoute(func(c *gin.Context) {
		uri := c.Request.URL.String()
		msg := fmt.Sprintf("Page not found: %s", uri)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": msg})
	})

	return r
}

func PrintErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // execute all the handlers

		errorToPrint := c.Errors.Last()
		if errorToPrint != nil {
			log.Printf("API ERROR: %+v", errorToPrint.Error())
		}
	}
}

// health-check to test service is up
func ping() func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
