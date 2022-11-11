package queue

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"transaction_queue/db"

	"github.com/go-redis/redis/v8"
)

type QueueReadTimeoutError struct{}

func (*QueueReadTimeoutError) Error() string {
	return "XREAD timeout"
}

type QueueService struct {
	redisClient *redis.Client
	holdService *db.HoldService
}

func NewService(redisClient *redis.Client, holdService *db.HoldService) *QueueService {
	return &QueueService{redisClient: redisClient, holdService: holdService}
}

func (s QueueService) PushTx(ctx context.Context, newTx db.Transaction) {
	values := map[string]interface{}{
		"id":      newTx.Id,
		"user_id": newTx.UserId,
		"amount":  newTx.Amount,
		"type":    newTx.Type,
	}
	streamName := fmt.Sprintf("queue-%d", newTx.UserId)
	result, err := s.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream:       streamName,
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values:       values,
	}).Result()
	if err != nil {
		log.Println("Push tx error: %v", err.Error())
	}
	log.Printf("Tx %v pushed to stream %s: %s", newTx, streamName, result)
}

func (s QueueService) handleMessage(ctx context.Context, msg redis.XMessage) error {
	if s.holdService == nil {
		return fmt.Errorf("Need to initialize hold service")
	}

	fmt.Printf("handle message %v\n", msg)
	isTxConfirmed := true
	fmt.Printf("id type: %T\n", msg.Values["id"])
	txIdField, ok := msg.Values["id"].(string)
	if !ok {
		return fmt.Errorf("Cannot decode message value")
	}
	txId, err := strconv.Atoi(txIdField)
	if err != nil {
		return fmt.Errorf("Cannot decode message value")
	}
	if isTxConfirmed {
		s.holdService.Delete(txId)
	} else {
		s.holdService.Reverse(ctx, txId)
	}
	return nil
}

func (s QueueService) GetUsersTxs(ctx context.Context, userIds []int) error {
	streams := make([]string, 0, len(userIds)*2)
	for _, id := range userIds {
		streams = append(streams, fmt.Sprintf("queue-%d", id))
	}
	for i := 0; i < len(userIds); i++ {
		streams = append(streams, "0")
	}
	log.Printf("watch streams %v", streams)
	xreadSlice := s.redisClient.XRead(ctx, &redis.XReadArgs{
		Streams: streams,
		Count:   1,
		Block:   time.Second * 10,
	})
	result, err := xreadSlice.Result()
	if err != nil {

		return &QueueReadTimeoutError{}
	}
	log.Printf("XREAD result: %v\n", result)
	for _, stream := range result {
		log.Printf("stream: %v", stream)
		for _, message := range stream.Messages {
			log.Printf("message: %v", message)
			if err := s.handleMessage(ctx, message); err != nil {
				log.Printf("error in handle message %v: %s\n", message, err.Error())
			} else {
				_, err := s.redisClient.XDel(ctx, stream.Stream, message.ID).Result()
				if err != nil {
					log.Printf("Error delete message: %v\n", message)
				} else {
					log.Printf("Delete message: %v\n", message)
				}
			}
		}
	}
	log.Println("messages ok")

	return nil
}
