package handler

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"github.com/yourname/user-api/internal/logger"
	"github.com/yourname/user-api/internal/models"
	"github.com/yourname/user-api/internal/repository"
	"github.com/yourname/user-api/internal/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	repo      *repository.UserRepository
	validator *validator.Validate
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		repo:      repo,
		validator: validator.New(),
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	requestID := c.Locals("requestId").(string)

	var req models.UserRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Logger.Warn("failed to parse request body", zap.String("request_id", requestID), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		logger.Logger.Warn("validation failed", zap.String("request_id", requestID), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "validation failed: " + err.Error(),
		})
	}

	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		logger.Logger.Warn("invalid date format", zap.String("request_id", requestID), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid date format, use YYYY-MM-DD",
		})
	}

	userID, err := h.repo.CreateUser(c.Context(), req.Name, dob)
	if err != nil {
		logger.Logger.Error("failed to create user in database", zap.String("request_id", requestID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "failed to create user",
		})
	}

	logger.Logger.Info("user created successfully", zap.String("request_id", requestID), zap.Int32("user_id", userID))
	return c.Status(fiber.StatusCreated).JSON(models.UserCreateResponse{
		ID:   userID,
		Name: req.Name,
		DOB:  req.DOB,
	})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	requestID := c.Locals("requestId").(string)

	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		logger.Logger.Warn("invalid user id", zap.String("request_id", requestID), zap.String("id", idStr))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid user id",
		})
	}

	user, err := h.repo.GetUser(c.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logger.Info("user not found", zap.String("request_id", requestID), zap.Int32("user_id", int32(id)))
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		}
		logger.Logger.Error("failed to retrieve user", zap.String("request_id", requestID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "failed to retrieve user",
		})
	}

	age := service.CalculateAge(user.DOB)
	logger.Logger.Info("user retrieved", zap.String("request_id", requestID), zap.Int32("user_id", user.ID))
	return c.Status(fiber.StatusOK).JSON(models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.DOB.Format("2006-01-02"),
		Age:  age,
	})
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	requestID := c.Locals("requestId").(string)

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := int64((page - 1) * limit)
	users, err := h.repo.GetAllUsers(c.Context(), int64(limit), offset)
	if err != nil {
		logger.Logger.Error("failed to retrieve users", zap.String("request_id", requestID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "failed to retrieve users",
		})
	}

	count, err := h.repo.GetUserCount(c.Context())
	if err != nil {
		logger.Logger.Error("failed to get user count", zap.String("request_id", requestID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "failed to retrieve user count",
		})
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		age := service.CalculateAge(user.DOB)
		userResponses[i] = models.UserResponse{
			ID:   user.ID,
			Name: user.Name,
			DOB:  user.DOB.Format("2006-01-02"),
			Age:  age,
		}
	}

	logger.Logger.Info("users list retrieved", zap.String("request_id", requestID), zap.Int("count", len(users)))
	return c.Status(fiber.StatusOK).JSON(models.PaginatedUsersResponse{
		Users: userResponses,
		Total: count,
		Page:  page,
		Limit: limit,
	})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	requestID := c.Locals("requestId").(string)

	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		logger.Logger.Warn("invalid user id", zap.String("request_id", requestID), zap.String("id", idStr))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid user id",
		})
	}

	var req models.UserRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Logger.Warn("failed to parse request body", zap.String("request_id", requestID), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		logger.Logger.Warn("validation failed", zap.String("request_id", requestID), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "validation failed: " + err.Error(),
		})
	}

	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		logger.Logger.Warn("invalid date format", zap.String("request_id", requestID), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid date format, use YYYY-MM-DD",
		})
	}

	err = h.repo.UpdateUser(c.Context(), int32(id), req.Name, dob)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logger.Info("user not found for update", zap.String("request_id", requestID), zap.Int32("user_id", int32(id)))
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		}
		logger.Logger.Error("failed to update user", zap.String("request_id", requestID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "failed to update user",
		})
	}

	logger.Logger.Info("user updated", zap.String("request_id", requestID), zap.Int32("user_id", int32(id)))
	return c.Status(fiber.StatusOK).JSON(models.UserUpdateResponse{
		ID:   int32(id),
		Name: req.Name,
		DOB:  req.DOB,
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	requestID := c.Locals("requestId").(string)

	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		logger.Logger.Warn("invalid user id", zap.String("request_id", requestID), zap.String("id", idStr))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid user id",
		})
	}

	err = h.repo.DeleteUser(c.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logger.Info("user not found for delete", zap.String("request_id", requestID), zap.Int32("user_id", int32(id)))
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		}
		logger.Logger.Error("failed to delete user", zap.String("request_id", requestID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "failed to delete user",
		})
	}

	logger.Logger.Info("user deleted", zap.String("request_id", requestID), zap.Int32("user_id", int32(id)))
	return c.Status(fiber.StatusNoContent).Send(nil)
}
