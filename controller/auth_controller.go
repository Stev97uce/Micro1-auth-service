package controller

import (
	"auth-service/model"
	"auth-service/service"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	AuthService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req model.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := c.AuthService.FindByEmail(req.Email)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar el email"})
		return
	}

	if existingUser != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "El correo ya está registrado"})
		return
	}

	err = c.AuthService.CreateUser(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Usuario registrado con éxito",
		"user_id": uuid.New().String(),
	})
}

func (c *AuthController) Login(ctx *gin.Context, rdb *redis.Client) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := c.AuthService.FindByEmail(req.Email)
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	token := uuid.NewString()

	rdb.Set(ctx, "session:"+token, user.ID, time.Hour)

	ctx.JSON(200, gin.H{
		"token": token,
	})
}
