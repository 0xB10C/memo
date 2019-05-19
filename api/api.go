package main

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/0xb10c/memo/api/database"
	"github.com/0xb10c/memo/api/config"
	"github.com/gin-contrib/cors"
	
)

func main() {

	if config.GetBool("api.production") {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"127.0.0.1", "localhost", "mempool.observer"}
	router.Use(cors.New(corsConfig))


	err := database.SetupDatabase()
	if err != nil {
		panic(fmt.Errorf("Failed to setup database: %v", err))
	}

	api := router.Group("/api")
	{
		api.GET("/mempool", getMempool)
		api.GET("/recentBlocks", getRecentBlocks)
	}

	portString := ":" + config.GetString("api.port")
	router.Run(portString)
}


func getMempool(c *gin.Context) {

	timestamp, byCount, megabyteMarkersJSON, mempoolSize, err := database.GetMempool()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Possible REFACTOR: write to database as blob not JSON String to 
	// skip the marshalling when writing and unmarshalling when reading
	// from the database
 	var feerateMap map[int]int
	json.Unmarshal([]byte(byCount), &feerateMap)

	var megabyteMarkers []int
	json.Unmarshal([]byte(megabyteMarkersJSON), &megabyteMarkers)

	c.JSON(http.StatusOK, gin.H{
		"timestamp":					timestamp.Unix(),
		"feerateMap": 				feerateMap,
		"megabyteMarkers":		megabyteMarkers,
		"mempoolSize":	 			mempoolSize,
	})
}



func getRecentBlocks(c *gin.Context) {

	blocks, err := database.GetRecentBlocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, blocks)
}

