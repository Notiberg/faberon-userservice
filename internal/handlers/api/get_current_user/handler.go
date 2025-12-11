package get_current_user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/m04kA/SMC-UserService/internal/handlers/api"
	"github.com/m04kA/SMC-UserService/internal/handlers/middleware"
	"github.com/m04kA/SMC-UserService/internal/service/user/models"
	userservice "github.com/m04kA/SMC-UserService/internal/service/user"
)

type Handler struct {
	service *userservice.Service
	log     Logger
}

func NewHandler(service *userservice.Service, log Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

// Handle GET /users/me
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.log.Warn("GET /users/me - Unauthorized access attempt")
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	user, err := h.service.GetUserWithCars(r.Context(), userID)
	if err != nil {
		if errors.Is(err, userservice.ErrUserNotFound) {
			h.log.Info("GET /users/me - User not found, creating new user: user_id=%d", userID)
			// Auto-create user if not found
			createInput := models.CreateUserInputDTO{
				TGUserID: userID,
				Name:     fmt.Sprintf("User_%d", userID),
				Role:     "customer",
			}
			_, err := h.service.CreateUser(r.Context(), createInput)
			if err != nil {
				h.log.Error("GET /users/me - Failed to create user: user_id=%d, error=%v", userID, err)
				api.RespondInternalError(w)
				return
			}
			// Get the newly created user
			user, err = h.service.GetUserWithCars(r.Context(), userID)
			if err != nil {
				h.log.Error("GET /users/me - Failed to get newly created user: user_id=%d, error=%v", userID, err)
				api.RespondInternalError(w)
				return
			}
		} else {
			h.log.Error("GET /users/me - Failed to get user: user_id=%d, error=%v", userID, err)
			api.RespondInternalError(w)
			return
		}
	}

	h.log.Info("GET /users/me - User retrieved successfully: user_id=%d", userID)
	api.RespondJSON(w, http.StatusOK, user)
}
