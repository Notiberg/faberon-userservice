package get_user_by_id

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/m04kA/SMC-UserService/internal/handlers/api"
	"github.com/m04kA/SMC-UserService/internal/service/user"
)

type Handler struct {
	service *user.Service
	log     Logger
}

func NewHandler(service *user.Service, log Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["tg_user_id"]

	// Парсим tg_user_id из URL
	tgUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.log.Warn("GET /internal/users/{tg_user_id} - invalid user ID format: %s", userIDStr)
		api.RespondError(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	// Получаем пользователя с автомобилями
	userWithCars, err := h.service.GetUserWithCars(r.Context(), tgUserID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			h.log.Warn("GET /internal/users/%d - user not found", tgUserID)
			api.RespondError(w, http.StatusNotFound, "user not found")
			return
		}

		h.log.Error("GET /internal/users/%d - failed to get user: %v", tgUserID, err)
		api.RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.log.Info("GET /internal/users/%d - success", tgUserID)
	api.RespondJSON(w, http.StatusOK, userWithCars)
}
