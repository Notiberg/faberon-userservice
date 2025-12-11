package delete_car

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/m04kA/SMC-UserService/internal/handlers/api"
	"github.com/m04kA/SMC-UserService/internal/handlers/middleware"
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

// Handle DELETE /users/me/cars/{car_id}
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.log.Warn("DELETE /users/me/cars/{car_id} - Unauthorized access attempt")
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	role, err := middleware.GetRoleFromContext(r.Context())
	if err != nil {
		h.log.Warn("DELETE /users/me/cars/{car_id} - Failed to get role: user_id=%d", userID)
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	carIDStr := vars["car_id"]
	if carIDStr == "" {
		h.log.Warn("DELETE /users/me/cars/{car_id} - Car ID missing: user_id=%d", userID)
		api.RespondBadRequest(w, "Car ID is required")
		return
	}

	carID, err := strconv.ParseInt(carIDStr, 10, 64)
	if err != nil {
		h.log.Warn("DELETE /users/me/cars/{car_id} - Invalid car ID format: user_id=%d, car_id_str=%s", userID, carIDStr)
		api.RespondBadRequest(w, "Invalid car ID")
		return
	}

	err = h.service.DeleteCar(r.Context(), userID, carID, role)
	if err != nil {
		if errors.Is(err, userservice.ErrCarNotFound) {
			h.log.Warn("DELETE /users/me/cars/{car_id} - Car not found: user_id=%d, car_id=%d", userID, carID)
			api.RespondCarNotFound(w)
			return
		}
		if errors.Is(err, userservice.ErrCarAccessDenied) {
			h.log.Warn("DELETE /users/me/cars/{car_id} - Access denied: user_id=%d, car_id=%d, role=%s", userID, carID, role)
			api.RespondCarAccessDenied(w)
			return
		}
		h.log.Error("DELETE /users/me/cars/{car_id} - Failed to delete car: user_id=%d, car_id=%d, error=%v", userID, carID, err)
		api.RespondInternalError(w)
		return
	}

	h.log.Info("DELETE /users/me/cars/{car_id} - Car deleted successfully: user_id=%d, car_id=%d, role=%s", userID, carID, role)
	w.WriteHeader(http.StatusNoContent)
}
