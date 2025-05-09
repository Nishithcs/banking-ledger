package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/Nishithcs/banking-ledger/internal/processor"
	"github.com/Nishithcs/banking-ledger/pkg"
)

func main() {
	ctx := context.Background()

	// Setup dependencies
	queue := setupRabbitMQ()
	defer queue.Close()

	mongoClient := setupMongoDB(ctx)

	db := setupPostgres(ctx)
	defer db.Close(ctx)

	queueName := os.Getenv("RABBITMQ_QUEUE_NAME")
	msgs, err := queue.Consume(queueName)
	if err != nil {
		log.Fatalf("Failed to consume messages from queue: %v", err)
	}

	numWorkers := getWorkerCount()
	wg := sync.WaitGroup{}

	log.Printf(" [*] Starting %d workers for queue '%s'...", numWorkers, queueName)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go startWorker(ctx, i, msgs, queueName, db, mongoClient, &wg)
	}

	wg.Wait()
}

// ----------------- INIT HELPERS --------------------

// setupRabbitMQ initializes and returns a RabbitMQ connection
func setupRabbitMQ() pkg.MessageQueue {
	queue := &pkg.RabbitMQ{}
	if err := queue.Connect(pkg.AmqpURL()); err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	return queue
}

// setupMongoDB initializes and returns a MongoDB client
func setupMongoDB(ctx context.Context) pkg.MongoDBClient {
	client, err := pkg.NewMongoDBClient(ctx, "mongodb://myuser:mypassword@mongodb:27017", "bank")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return client
}

// setupPostgres initializes and returns a PostgreSQL connection
func setupPostgres(ctx context.Context) pkg.Database {
	var postgres pkg.Database = &pkg.PostgresDB{}
	err := postgres.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME")))
	if err != nil {
		panic(err)
	}
	return postgres
}


func getWorkerCount() int {
	if val, err := strconv.Atoi(os.Getenv("NUM_WORKERS")); err == nil && val > 0 {
		return val
	}
	return 4
}

// ----------------- WORKER --------------------

func startWorker(ctx context.Context, workerID int, msgs <-chan []byte, 
	queueName string, database pkg.Database, mongoClient pkg.MongoDBClient, wg *sync.WaitGroup) {
		
	defer wg.Done()
	log.Printf("Worker %d started", workerID)

	for msg := range msgs {
		log.Printf("Worker %d received message: %s", workerID, msg)

		switch queueName {
		case "account_creator":
			var account processor.AccountData
			if err := json.Unmarshal(msg, &account); err != nil {
				log.Printf("Worker %d: JSON decode error: %v", workerID, err)
				continue
			}

			processor := processor.CreateAccountProcessor{
				ProcessWorker: processor.ProcessWorker{
					Database:     database,
					MongoDbConn:  mongoClient,
				},
				Data: account,
			}

			if err := processor.CreateAccount(ctx); err != nil {
				log.Printf("Worker %d: Account creation failed: %v", workerID, err)
			}

		case "transaction_processor":
			var txn processor.TransactionData
			if err := json.Unmarshal(msg, &txn); err != nil {
				log.Printf("Worker %d: JSON decode error: %v", workerID, err)
				continue
			}

			processor := processor.TransactionProcessor{
				ProcessWorker: processor.ProcessWorker{
					Database:     database,
					MongoDbConn:  mongoClient,
				},
				Data: txn,
			}

			if err := processor.ProcessTransaction(ctx); err != nil {
				log.Printf("Worker %d: Transaction processing failed: %v", workerID, err)
			}

		default:
			log.Printf("Worker %d: Unknown queue name '%s'", workerID, queueName)
		}
	}
}
