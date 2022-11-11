package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"transaction_queue/db"
	"transaction_queue/queue"
	"transaction_queue/utils"

	"github.com/go-redis/redis/v8"

	_ "github.com/lib/pq"
)

var (
	dbDriver           string
	dbDSN              string
	redisHost          string
	redisPort          int
	redisDb            int
	userService        *db.UserService
	transactionService *db.TransactionService
	holdService        *db.HoldService
	queueService       *queue.QueueService
)

func init() {
	dbDriver, dbDSN = utils.GetDbParams()
	redisHost, redisPort, redisDb = utils.GetRedisParams()
}

func transactionProcess() {
	context := context.Background()

	for {
		userIds := make([]int, 0)

		users, err := userService.GetAll()
		if err != nil {
			log.Fatal(err)
		}

		for _, user := range users {
			userIds = append(userIds, user.Id)
		}
		fmt.Printf("Read queues for users %v\n", users)
		err = queueService.GetUsersTxs(context, userIds)
		if err != nil {
			re, ok := err.(*queue.QueueReadTimeoutError)
			if ok {
				fmt.Printf("No new messages: %v\n", re.Error())
			} else {
				log.Fatal(err)
			}
		}
	}
}

func main() {
	dbConn, err := sql.Open(dbDriver, dbDSN)
	if err != nil {
		log.Fatalf("Cannot connect db: %v", err)
	}
	transactionService = db.NewTransactionService(dbConn)
	userService = db.NewUserService(dbConn)
	holdService = db.NewHoldService(dbConn)

	redisAddr := fmt.Sprintf("%s:%d", redisHost, redisPort)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       redisDb,
	})
	log.Printf("Connect to redis %s/%d\n", redisAddr, redisDb)
	queueService = queue.NewService(redisClient, holdService)

	transactionProcess()
}
