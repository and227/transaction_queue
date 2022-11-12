package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"transaction_queue/db"
	"transaction_queue/queue"
	"transaction_queue/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "transaction_queue/docs"
)

// @title       Simple transaction queue
// @version     1.0
// @description Simple transaction queue
// @host        localhost:8080
// @BasePath    /

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

// @Summary     create new transaction
// @Tags        transaction
// @ID          create-transaction
// @Accept      json
// @Produce     json
// @Param       input body     db.TransactionCreateDTO true "transaction info"
// @Success     200   {object} db.Transaction
// @Failure 	400,500 {object} gin.H  "error response"
// @Router      /transactions [post]
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

// @Summary     create new user
// @Tags        user
// @ID          create-user
// @Accept      json
// @Produce     json
// @Param       input body     db.UserCreateDTO true "user info"
// @Success     200   {object} db.User
// @Failure 	400,500 {object} gin.H  "error response"
// @Router      /users [post]
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

// @Summary     get user with balance
// @Tags        user
// @ID          get-user
// @Accept      json
// @Produce     json
// @Param       user_id path int true "user id"
// @Success     200 {object} db.UserWithBalanceOutDTO
// @Failure 	400,500 {object} gin.H  "error response"
// @Router      /users/{user_id} [get]
func userWithBalanceGet(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userWithBalance, err := userService.GetUserWithBalance(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userWithBalance)
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
	router.GET("users/:id", userWithBalanceGet)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":8080")
}
