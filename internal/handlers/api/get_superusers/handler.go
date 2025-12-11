package get_superusers

import (
	"net/http"

	"github.com/m04kA/SMC-UserService/internal/handlers/api"
	"github.com/m04kA/SMC-UserService/internal/service/user"
)

const (
	ErrInternalServer = "внутренняя ошибка сервера"
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

// Response структура для ответа со списком superuser ID
type Response struct {
	SuperUserIDs []int64 `json:"super_user_ids"`
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	// Получаем список всех superuser ID
	superUserIDs, err := h.service.GetSuperUsers(r.Context())
	if err != nil {
		h.log.Error("GET /internal/users/superusers - failed to get superusers: %v", err)
		api.RespondError(w, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	h.log.Info("GET /internal/users/superusers - success, found %d superusers", len(superUserIDs))
	api.RespondJSON(w, http.StatusOK, Response{SuperUserIDs: superUserIDs})
}
