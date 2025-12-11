package create_car

import (
	"errors"
	"net/http"

	"github.com/m04kA/SMC-UserService/internal/handlers/api"
	"github.com/m04kA/SMC-UserService/internal/handlers/middleware"
	userservice "github.com/m04kA/SMC-UserService/internal/service/user"
	"github.com/m04kA/SMC-UserService/internal/service/user/models"
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

// Handle POST /users/me/cars
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.log.Warn("POST /users/me/cars - Unauthorized access attempt")
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	var input models.CreateCarInputDTO
	if err := api.DecodeJSON(r, &input); err != nil {
		h.log.Warn("POST /users/me/cars - Invalid request body: user_id=%d, error=%v", userID, err)
		api.RespondBadRequest(w, "Invalid request body")
		return
	}

	car, err := h.service.CreateCar(r.Context(), userID, input)
	if err != nil {
		if errors.Is(err, userservice.ErrUserNotFound) {
			h.log.Warn("POST /users/me/cars - User not found: user_id=%d", userID)
			api.RespondUserNotFound(w)
			return
		}
		h.log.Error("POST /users/me/cars - Failed to create car: user_id=%d, error=%v", userID, err)
		api.RespondInternalError(w)
		return
	}

	h.log.Info("POST /users/me/cars - Car created successfully: user_id=%d, car_id=%d", userID, car.ID)
	api.RespondJSON(w, http.StatusCreated, car)
}
