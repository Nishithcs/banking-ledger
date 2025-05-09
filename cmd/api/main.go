package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Nishithcs/banking-ledger/internal/handlers"
	"github.com/Nishithcs/banking-ledger/pkg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Standard Go context for graceful shutdowns and DB calls
	ctx := context.Background()

	// Initialize Gin router
	router := setupRouter()

	// Setup dependencies
	queue := setupRabbitMQ()
	defer queue.Close()

	mongoClient := setupMongoDB(ctx)

	db := setupPostgres(ctx)
	defer db.Close(ctx)

	// Register routes
	router.POST("/createAccount", handlers.CreateAccountHandler(ctx, queue, "account_creator"))
	router.POST("/transact", handlers.TransactionHandler(ctx, queue, "transaction_processor"))
	router.GET("/account/:accountNumber/transactionHistory", handlers.GetTransactionHistoryHandler(ctx, mongoClient))
	router.GET("/account/status/:accountNumber", handlers.GetAccountStatusHandler(ctx, db, mongoClient))

	// Run the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRouter initializes Gin with CORS config
func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	return router
}

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
