package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Nishithcs/banking-ledger/internal/handlers"
	internal "github.com/Nishithcs/banking-ledger/pkg"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// // Create RabbitMQ connection
	// aqmpConn, err := internal.CreateAMQPConnection(
	// 	"amqp://" +
	// 		os.Getenv("RABBITMQ_USER") + ":" +
	// 		os.Getenv("RABBITMQ_PASSWORD") + "@" +
	// 		os.Getenv("RABBITMQ_HOST") + ":" +
	// 		os.Getenv("RABBITMQ_PORT") + "/")
	// if err != nil {
	// 	panic(err)
	// }

	// defer internal.CloseAMQPConnection(aqmpConn)

	// ctx := gin.Context{}

	// amqpChannel, err := aqmpConn.Channel()
	// if err != nil {
	// 	panic(err)
	// }

	// defer internal.CloseAMQPChannel(amqpChannel)

	// // Declare queue
	// createAccountQueue, err := internal.QueueDeclare(amqpChannel, "account_creator", true, false, false, false)
	// if err != nil {
	// 	panic(err)
	// }

	ctx := gin.Context{}

	// Rabbitmq account creation
	var queue internal.MessageQueue = &internal.RabbitMQ{}

	err := queue.Connect(internal.AmqpURL())
	if err != nil {
		panic(err)
	}

	defer queue.Close()
	// Done

	router.POST("/createAccount", handlers.CreateAccountHandler(&ctx, queue, "account_creator"))
	router.POST("/transact", handlers.TransactionHandler(&ctx, queue, "transaction_processor"))

	// Initialize Elasticsearch client
	esConfig := elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTICSEARCH_URL")},
	}

	esClient, err := internal.NewElasticsearchClient(esConfig)
	if err != nil {
		panic(fmt.Sprintf("Error creating Elasticsearch client: %s", err))
	}

	// MongoDB starts	
	// ctx = context.Background()
	mongoDbClient, err := internal.NewMongoDBClient(&ctx, "mongodb://myuser:mypassword@mongodb:27017", "bank")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating MongoDB client: %v\n", err)
		os.Exit(1)
	}
	// MongoDb completes

	// Test the connection
	res, err := esClient.Info()
	if err != nil {
		panic(fmt.Sprintf("Error getting Elasticsearch info: %s", err))
	}
	defer res.Body.Close()


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


	router.GET("/account/:accountNumber/transactionHistory", handlers.GetTransactionHistoryHandler(&ctx, esClient, mongoDbClient))

	router.GET("/account/status/:referenceId", handlers.GetAccountStatusHandler(&ctx, conn, esClient, mongoDbClient))

	router.Run(":8080")
}
