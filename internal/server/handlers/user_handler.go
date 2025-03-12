package handlers

import (
	"DBackend/internal/database"
	"DBackend/model"
	"DBackend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserHandler struct to handle user-related requests
type UserHandler struct {
	db database.Service
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(db database.Service) *UserHandler {
	return &UserHandler{db: db}
}

// RegisterHandler handles user registration
func (h *UserHandler) RegisterHandler(c *fiber.Ctx) error {
	data := new(struct {
		FirstName  string `json:"first_name"`
		SecondName string `json:"second_name"`
		Role       string `json:"role"`
		Email      string `json:"email"`
		Password   string `json:"password"`
	})

	// Parse request body
	if err := c.BodyParser(data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if data.FirstName == "" || data.SecondName == "" || data.Role == "" || data.Email == "" || data.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required"})
	}

	// Check if user already exists
	existingUser, err := h.db.User().FindByEmail(c.Context(), data.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// If user exists, add role if not already assigned
	if existingUser != nil {
		for _, role := range existingUser.Roles {
			if role == data.Role {
				return c.Status(400).JSON(fiber.Map{"error": "User already has this role"})
			}
		}

		// Append the new role and update user
		existingUser.Roles = append(existingUser.Roles, data.Role)
		_, err := h.db.User().UpdateRoles(c.Context(), existingUser.Email, existingUser.Roles)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update user roles"})
		}

		// Insert role-specific data
		err = h.db.User().CreateRoleData(c.Context(), existingUser.ID, data.Role)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create role data"})
		}

		return c.JSON(fiber.Map{"message": "Role added successfully"})
	}

	// Hash password securely
	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to encrypt password"})
	}

	// Create new user
	user := model.User{
		ID:         primitive.NewObjectID(),
		FirstName:  data.FirstName,
		SecondName: data.SecondName,
		Email:      data.Email,
		Password:   hashedPassword,
		Roles:      []string{data.Role}, // Assign first role
	}

	// Insert user into the database
	insertResult, err := h.db.User().CreateUser(c.Context(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// Insert role-specific data
	err = h.db.User().CreateRoleData(c.Context(), insertResult.InsertedID.(primitive.ObjectID), data.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create role data"})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}


func (h *UserHandler) GetUserDetailsHandler(c *fiber.Ctx) error {
	// Get the token from the context
	userToken := c.Locals("token")
	tokenStr, ok := userToken.(string)
	if !ok || tokenStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid token"})
	}
	// Validate token and extract claims
	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid token"})
	}
	// Convert claims.UserID to primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	// Find the user in the database
	user, err := h.db.User().FindByID(c.Context(), "users", id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	// Type assert user to model.User and remove sensitive fields
	userDetails, ok := user.(*model.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}
	userDetails.Password = ""
	return c.JSON(userDetails)
}
