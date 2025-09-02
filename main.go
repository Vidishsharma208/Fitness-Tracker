package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func connectMongo() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://vidishsharma20:fitness1234@cluster1.hxfzg0k.mongodb.net/?retryWrites=true&w=majority&appName=Cluster1")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ MongoDB connection error:", err)
	}

	collection = client.Database("fitness_tracker").Collection("entries")
	log.Println("✅ Connected to MongoDB!")
}

func saveData(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	data["timestamp"] = time.Now()

	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save entry"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entry saved successfully"})
}

func main() {
	connectMongo()

	r := gin.Default()
	r.Use(CORSMiddleware())

	r.POST("/save", saveData)

	log.Println("🚀 Server running on http://localhost:8080")
	r.Run(":8080")
}
