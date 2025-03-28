package handlers

import (
    "DBackend/internal/database"
    "DBackend/model"
    "time"

    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateGrant handles the creation of a new grant
func CreateGrant(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var grant model.Grant
        if err := c.BodyParser(&grant); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
        }

        grant.ID = primitive.NewObjectID()
        grant.CreatedAt = time.Now()
        grant.UpdatedAt = time.Now()

        result, err := db.Founder().CreateGrant(c.Context(), grant)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to create grant"})
        }

        return c.Status(201).JSON(fiber.Map{
            "message": "Grant created successfully",
            "id":      result,
        })
    }
}

// GetGrantByID handles retrieving a grant by its ID
func GetGrantByID(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, err := primitive.ObjectIDFromHex(c.Params("id"))
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid grant ID"})
        }

        grant, err := db.Founder().GetGrantByID(c.Context(), id)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "Grant not found"})
        }

        return c.JSON(grant)
    }
}

// UpdateGrant handles updating an existing grant
func UpdateGrant(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, err := primitive.ObjectIDFromHex(c.Params("id"))
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid grant ID"})
        }

        var updates model.Grant
        if err := c.BodyParser(&updates); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
        }

        updates.UpdatedAt = time.Now()

        err = db.Founder().UpdateGrant(c.Context(), id, updates)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to update grant"})
        }

        return c.JSON(fiber.Map{"message": "Grant updated successfully"})
    }
}

// DeleteGrant handles deleting a grant
func DeleteGrant(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, err := primitive.ObjectIDFromHex(c.Params("id"))
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid grant ID"})
        }

        err = db.Founder().DeleteGrant(c.Context(), id)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to delete grant"})
        }

        return c.JSON(fiber.Map{"message": "Grant deleted successfully"})
    }
}

// ApplyForGrant handles grant application submission
func ApplyForGrant(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var application model.GrantApplication
        if err := c.BodyParser(&application); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
        }

        application.ID = primitive.NewObjectID()
        application.CreatedAt = time.Now()
        application.Status = "Pending"

        // Get user ID from JWT token
        userID, ok := c.Locals("user_id").(string)
        if !ok {
            return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
        }

        founderID, err := primitive.ObjectIDFromHex(userID)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
        }
        application.FounderID = founderID

        result, err := db.Founder().SubmitGrantApplication(c.Context(), application)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to submit application"})
        }

        return c.Status(201).JSON(fiber.Map{
            "message":       "Application submitted successfully",
            "applicationId": result,
        })
    }
}

// GetGrantApplications handles retrieving all grant applications
func GetGrantApplications(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Check if user is admin or founder
        userID, ok := c.Locals("user_id").(string)
        if !ok {
            return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
        }

        role, ok := c.Locals("user_role").(string)
        if !ok {
            return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
        }

        var applications []model.GrantApplication
        var err error

        if role == "admin" {
            // Admins can see all applications
            applications, err = db.Founder().GetAllGrantApplications(c.Context())
        } else {
            // Founders can only see their own applications
            founderID, err := primitive.ObjectIDFromHex(userID)
            if err != nil {
                return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
            }
            applications, err = db.Founder().GetFounderGrantApplications(c.Context(), founderID)
        }

        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve applications"})
        }

        return c.JSON(fiber.Map{"applications": applications})
    }
}

// GetGrantApplicationByID handles retrieving a specific grant application
func GetGrantApplicationByID(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, err := primitive.ObjectIDFromHex(c.Params("id"))
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid application ID"})
        }

        application, err := db.Founder().GetGrantApplicationByID(c.Context(), id)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "Application not found"})
        }

        // Check if user is authorized to view this application
        userID, ok := c.Locals("user_id").(string)
        if !ok {
            return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
        }

        role, ok := c.Locals("user_role").(string)
        if !ok {
            return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
        }

        // If not admin, check if application belongs to user
        if role != "admin" {
            founderID, err := primitive.ObjectIDFromHex(userID)
            if err != nil {
                return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
            }

            if application.FounderID != founderID {
                return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
            }
        }

        return c.JSON(application)
    }
}

// UpdateGrantApplication handles updating a grant application
func UpdateGrantApplication(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, err := primitive.ObjectIDFromHex(c.Params("id"))
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid application ID"})
        }

        var updates struct {
            Status  string `json:"status"`
            Remarks string `json:"remarks"`
        }

        if err := c.BodyParser(&updates); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
        }

        // Only admins can update application status
        role, ok := c.Locals("user_role").(string)
        if !ok || role != "admin" {
            return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
        }

        err = db.Founder().UpdateGrantApplication(c.Context(), id, updates.Status, updates.Remarks)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to update application"})
        }

        return c.JSON(fiber.Map{"message": "Application updated successfully"})
    }
}