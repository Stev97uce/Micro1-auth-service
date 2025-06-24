package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"auth-service/config"
	"auth-service/controller"
	"auth-service/middleware"
	"auth-service/model"
	"auth-service/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando .env")
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar con PostgreSQL:", err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Error en migración de usuarios:", err)
	}

	rdb := config.NewRedisClient()
	defer rdb.Close()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Servidor iniciado en puerto:", port)

	authService := service.NewAuthService(db)
	authController := controller.NewAuthController(authService)

	router.POST("/register", authController.Register)
	router.POST("/login", func(c *gin.Context) {
		authController.Login(c, rdb)
	})

	router.GET("/protected", middleware.AuthMiddleware(rdb), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{
			"message": "Acceso autorizado",
			"user_id": userID,
		})
	})

	router.Run(":" + port)
}
