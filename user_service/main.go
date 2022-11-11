package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"transaction_queue/db"
	"transaction_queue/queue"
	"transaction_queue/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

var (
	dbDriver           string
	dbDSN              string
	redisHost          string
	redisPort          int
	redisDb            int
	transactionService *db.TransactionService
	userService        *db.UserService
	queueService       *queue.QueueService
)

func init() {
	dbDriver, dbDSN = utils.GetDbParams()
	redisHost, redisPort, redisDb = utils.GetRedisParams()
}

func transactionHandler(c *gin.Context) {
	var transactionRequest db.TransactionCreateDTO

	if err := c.ShouldBindJSON(&transactionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	transaction, err := transactionService.Create(ctx, transactionRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	queueService.PushTx(ctx, transaction)

	c.JSON(http.StatusCreated, transaction)
}

func userHandler(c *gin.Context) {
	var userRequest db.UserCreateDTO

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	user, err := userService.Create(ctx, userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func main() {
	dbConn, err := sql.Open(dbDriver, dbDSN)
	if err != nil {
		log.Fatalf("Cannot connect db: %v", err)
	}
	transactionService = db.NewTransactionService(dbConn)
	userService = db.NewUserService(dbConn)

	redisAddr := fmt.Sprintf("%s:%d", redisHost, redisPort)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
		Password: "",
		DB:       redisDb,
	})
	log.Printf("Connect to redis %s/%d\n", redisAddr, redisDb)
	queueService = queue.NewService(redisClient, nil)

	router := gin.Default()
	fmt.Println("Start server")
	router.POST("transactions", transactionHandler)
	router.POST("users", userHandler)
	router.Run(":8080")
}
