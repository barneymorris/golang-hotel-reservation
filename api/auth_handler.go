package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/betelgeusexru/golang-hotel-reservation/db"
	"github.com/betelgeusexru/golang-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User *types.User `json:"user"`
	Token string `json:"token"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.NewError(404, "invalid credentials")
		}

		return fiber.NewError(400, "invalid credentials")
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return fiber.NewError(400, "invalid credentials")
	}

	resp := AuthResponse{
		User: user,
		Token: CreateTokenFromUser(user),
	}

	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()

	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["expires"] = expires

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	secret := os.Getenv("JWT_SECRET")
	
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}

	return tokenStr 
}