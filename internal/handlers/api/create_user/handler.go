package create_user

import (
	"errors"
	"net/http"

	"github.com/m04kA/SMC-UserService/internal/handlers/api"
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

// Handle POST /users
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var input models.CreateUserInputDTO
	if err := api.DecodeJSON(r, &input); err != nil {
		h.log.Warn("POST /users - Invalid request body: %v", err)
		api.RespondBadRequest(w, "Invalid request body")
		return
	}

	user, err := h.service.CreateUser(r.Context(), input)
	if err != nil {
		if errors.Is(err, userservice.ErrUserAlreadyExists) {
			h.log.Warn("POST /users - User already exists: tg_user_id=%d", input.TGUserID)
			api.RespondUserAlreadyExists(w)
			return
		}
		h.log.Error("POST /users - Failed to create user: tg_user_id=%d, error=%v", input.TGUserID, err)
		api.RespondInternalError(w)
		return
	}

	h.log.Info("POST /users - User created successfully: tg_user_id=%d", user.TGUserID)
	api.RespondJSON(w, http.StatusCreated, user)
}
