package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/0xb10c/memo/config"
	"github.com/0xb10c/memo/database"
	"github.com/0xb10c/memo/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {

	if config.GetBool("api.production") {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	if config.GetBool("api.production") {
		corsConfig.AllowOrigins = []string{"https://mempool.observer/"}
	} else {
		corsConfig.AllowOrigins = []string{"*"}
	}

	router.Use(cors.New(corsConfig))
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	pool, err := database.SetupRedis()
	if err != nil {
		panic(fmt.Errorf("Failed to setup redis: %v", err))
	}

	go database.SetupMempoolEntriesCache(pool)

	router.Use(redisPoolMiddleware(pool))

	api := router.Group("/api")
	{
		api.GET("/mempool", getMempool)
		api.GET("/recentBlocks", getRecentBlocks)
		api.GET("/historicalMempool/:timeframe/:by", getHistoricalMempool)
		api.GET("/transactionStats", getTransactionStats)
		api.GET("/getMempoolEntries", getCachedMempoolEntries)
		api.GET("/getRecentFeerateAPIData", getRecentFeerateAPIEntries)
		api.GET("/getBlockEntries", getBlockEntries)
	}

	portString := ":" + config.GetString("api.port")
	router.Run(portString)
}

// ApiMiddleware will add the db connection to the context
func redisPoolMiddleware(pool *database.RedisPool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redisPool", pool)
		c.Next()
	}
}

func getMempool(c *gin.Context) {
	pool, ok := c.MustGet("redisPool").(*database.RedisPool)
	if !ok {
		logger.Error.Println("Could not get Redis Pool.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	timestamp, byCount, megabyteMarkersJSON, mempoolSize, err := pool.GetMempool()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
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
		"timestamp":       timestamp.Unix(),
		"feerateMap":      feerateMap,
		"megabyteMarkers": megabyteMarkers,
		"mempoolSize":     mempoolSize,
	})
}

func getRecentBlocks(c *gin.Context) {
	pool, ok := c.MustGet("redisPool").(*database.RedisPool)
	if !ok {
		logger.Error.Println("Could not get Redis Pool.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	blocks, err := pool.GetRecentBlocks()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, blocks)
}

func getBlockEntries(c *gin.Context) {
	pool, ok := c.MustGet("redisPool").(*database.RedisPool)
	if !ok {
		logger.Error.Println("Could not get Redis Pool.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	blocks, err := pool.GetBlockEntries()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, blocks)
}

func getHistoricalMempool(c *gin.Context) {
	pool, ok := c.MustGet("redisPool").(*database.RedisPool)
	if !ok {
		logger.Error.Println("Could not get Redis Pool.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	timeframe, err := strconv.ParseInt(c.Param("timeframe"), 10, 0)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid input error",
		})
		return
	}

	by := c.Param("by")
	if by != "byCount" && by != "byFee" && by != "bySize" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid input error",
		})
		return
	}

	mempoolStates, err := pool.GetHistorical(int(timeframe), by)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, mempoolStates)
}

func getTransactionStats(c *gin.Context) {
	pool, ok := c.MustGet("redisPool").(*database.RedisPool)
	if !ok {
		logger.Error.Println("Could not get Redis Pool.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	tss, err := pool.GetTransactionStats()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, tss)
}

func getMempoolEntries(c *gin.Context) {
	pool, ok := c.MustGet("redisPool").(*database.RedisPool)
	if !ok {
		logger.Error.Println("Could not get Redis Pool.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	mes, err := pool.GetMempoolEntries()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, mes)
}

func getRecentFeerateAPIEntries(c *gin.Context) {
	pool, ok := c.MustGet("redisPool").(*database.RedisPool)
	if !ok {
		logger.Error.Println("Could not get Redis Pool.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	entries, err := pool.GetRecentFeerateAPIEntries()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func getCachedMempoolEntries(c *gin.Context) {
	pool, ok := c.MustGet("redisPool").(*database.RedisPool)
	if !ok {
		logger.Error.Println("Could not get Redis Pool.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	entries, err := pool.GetMempoolEntriesCache()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(http.StatusOK, entries)
}
