package handlers

import (
	"DBackend/internal/database"
	"DBackend/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db database.Service
}

func NewAuthHandler(db database.Service) *AuthHandler {
	return &AuthHandler{db: db}
}

// LoginHandler handles user login
func (h *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	data := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	if err := c.BodyParser(data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	user, err := h.db.User().FindByEmail(c.Context(), data.Email)
	if err != nil || user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}
  id := user.ID.Hex()
	token, err := utils.GenerateJWT(id, user.Roles)
  c.Locals("token", token)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}

// LogoutHandler logs out the user by expiring the JWT cookie and invalidating the token
func (h *AuthHandler) LogoutHandler(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token != "" {
		token = utils.ExtractBearerToken(token)
	}
	if token == "" {
		return c.JSON(fiber.Map{"message": "Logged out successfully"})
	}
	// Optional: Store the token in a blacklist (if implementing token revocation)
	err := h.db.User().BlacklistToken(c.Context(), token)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to blacklist token"})
	}
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

// creating MeHandler that will return the credit from the jwt token passed after loginin
func (h *AuthHandler) MeHandler(c *fiber.Ctx) error {
	user_id, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	roles, ok := c.Locals("roles").([]string)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	return c.JSON(fiber.Map{"user_id": user_id, "roles": roles})
}


