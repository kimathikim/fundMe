package handlers

import (
	"time"
  "strings"

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
	// Find user by email
	user, err := h.db.User().FindByEmail(c.Context(), data.Email)
	if err != nil || user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	// Generate JWT token with user roles
	token, err := utils.GenerateJWT(user.Email, user.Roles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "roles",
		Value:    strings.Join(user.Roles, ","),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"roles":   user.Roles,
	})
}

// LogoutHandler logs out the user by expiring the JWT cookie and invalidating the token
func (h *AuthHandler) LogoutHandler(c *fiber.Ctx) error {
	// Expire the JWT cookie (for browser users)
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expire immediately
		HTTPOnly: true,
		Secure:   true, // Set to false for local development without HTTPS
		SameSite: "Lax",
	})

	// Get token from Authorization header (for non-browser users)
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

//creating MeHandler that will return the credit from the jwt token passed after loginin
func (h *AuthHandler) MeHandler(c *fiber.Ctx) error {
  email, ok := c.Locals("email").(string)
  if !ok {
    return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error"})
  }
  roles, ok := c.Locals("roles").([]string)
  if !ok {
    return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error"})
  }
  return c.JSON(fiber.Map{"email": email, "roles": roles})
}
