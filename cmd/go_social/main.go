package main

import (
	"database/sql"
	"os"

	"github.com/fatihesergg/go_social/internal/controller"
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	engine := gin.Default()
	base := engine.Group("/api/v1")

	dotenv := godotenv.Load()
	if dotenv != nil {
		panic("Error loading .env file")
	}

	DSN := os.Getenv("DATABASE_URI")
	if DSN == "" {
		panic("DATABASE_URI is not set")
	}

	if os.Getenv("JWT_SECRET") == "" {
		panic("JWT_SECRET is not set")
	}

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		panic("Error connecting to the database")
	}

	userStore := database.NewUserStore(db)
	storage := database.NewPostgresStorage(userStore)

	userController := controller.UserController{Storage: *storage}

	base.POST("/signup", userController.Signup)
	base.POST("/login", userController.Login)

	userRouter := base.Group("/users")
	userRouter.GET("/:id", userController.GetUserByID)
	if err := engine.Run(":3000"); err != nil {
		panic("Error starting the server")
	}

}
