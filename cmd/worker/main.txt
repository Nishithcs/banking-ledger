package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"

	"fmt"

	"github.com/Nishithcs/banking-ledger/internal/processor"
	internal "github.com/Nishithcs/banking-ledger/pkg"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5"
)

func main() {
	// Create RabbitMQ connection
	var queue internal.MessageQueue = &internal.RabbitMQ{}

	err := queue.Connect("amqp://" + 
	os.Getenv("RABBITMQ_USER") + ":" +
	os.Getenv("RABBITMQ_PASSWORD") + "@" +
	os.Getenv("RABBITMQ_HOST") + ":" +
	os.Getenv("RABBITMQ_PORT") + "/")
	if err != nil {
		log.Fatalf("Failed to connect to message queue: %v", err)
	}
	defer queue.Close()
	

	queueName := os.Getenv("RABBITMQ_QUEUE_NAME")
	log.Printf(queueName)
	msgs, err := queue.Consume(queueName)
	if err != nil {
		log.Fatalf("Failed to consume: %v", err)
	}

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME")))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	// Initialize Elasticsearch client
	esConfig := elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTICSEARCH_URL")},
	}

	esClient, err := internal.NewElasticsearchClient(esConfig)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Elasticsearch client: %v\n", err)
		os.Exit(1)
	}

	// Test the connection
	res, err := esClient.Info()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to Elasticsearch: %v\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	log.Println("Successfully connected to Elasticsearch")

	// Elasticsearch completes

	// MongoDB starts	
	ctx := context.Background()
	mongoDbClient, err := internal.NewMongoDBClient(ctx, "mongodb://myuser:mypassword@mongodb:27017", "bank")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating MongoDB client: %v\n", err)
		os.Exit(1)
	}
	// MongoDb completes

	wg := sync.WaitGroup{}

	wg.Add(1)
	
	// Start 4 worker goroutines
	numWorkers := 4 // Default value
	if workerCount, err := strconv.Atoi(os.Getenv("NUM_WORKERS")); err == nil && workerCount > 0 {
		numWorkers = workerCount
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		
		go func(workerID int, waitGroup *sync.WaitGroup) {
			defer waitGroup.Done()
			log.Printf("Worker %d started", workerID)
			log.Printf(queueName)
			
			for msg := range msgs {
				log.Printf("Worker %d received a message: %s", workerID, msg)
				
				switch os.Getenv("RABBITMQ_QUEUE_NAME") {
				case "account_creator":
					var accountInfo processor.AccountData
					
					err := json.Unmarshal(msg, &accountInfo)
					
					if err != nil {
						log.Printf("Error: %s\n", err)
						// msg.Ack(false)
						continue
					}
					
					processWorker := processor.CreateAccountProcessor{
						ProcessWorker: processor.ProcessWorker{
							PgxConn: conn,
							EsConn:  esClient,
							MongoDbConn: mongoDbClient,
						},
						Data: accountInfo,
					}
					
					err = processWorker.CreateAccount(context.Background())
					
					if err != nil {
						log.Println(err)
					}
					// msg.Ack(false)
				case "transaction_processor":
					var transactionInfo processor.TransactionData
					err := json.Unmarshal(msg, &transactionInfo)
					
					if err != nil {
						log.Println(err)
						// msg.Ack(false)
						continue
					}
					
					processWorker := processor.TransactionProcessor{
						ProcessWorker: processor.ProcessWorker{
							PgxConn: conn,
							EsConn:  esClient,
							MongoDbConn: mongoDbClient,
						},
						Data: transactionInfo,
					}
					
					err = processWorker.ProcessTransaction(context.Background())
					if err != nil {
						log.Println(err)
					}
					// msg.Ack(false)
				}
			}
		}(i, &wg)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	wg.Wait()
}
