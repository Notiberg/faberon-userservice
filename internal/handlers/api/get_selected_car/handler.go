package get_selected_car

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/m04kA/SMC-UserService/internal/handlers/api"
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

// Handle GET /internal/users/{tg_user_id}/cars/selected
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["tg_user_id"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.log.Warn("GET /internal/users/{tg_user_id}/cars/selected - Invalid user_id format: %s", userIDStr)
		api.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid user_id format",
		})
		return
	}

	car, err := h.service.GetSelectedCar(r.Context(), userID)
	if err != nil {
		if errors.Is(err, userservice.ErrCarNotFound) {
			h.log.Warn("GET /internal/users/{tg_user_id}/cars/selected - No selected car: user_id=%d", userID)
			api.RespondJSON(w, http.StatusNotFound, map[string]string{
				"error": "No selected car found",
			})
			return
		}
		h.log.Error("GET /internal/users/{tg_user_id}/cars/selected - Failed to get selected car: user_id=%d, error=%v", userID, err)
		api.RespondInternalError(w)
		return
	}

	h.log.Info("GET /internal/users/{tg_user_id}/cars/selected - Selected car retrieved: user_id=%d, car_id=%d", userID, car.ID)
	api.RespondJSON(w, http.StatusOK, car)
}
