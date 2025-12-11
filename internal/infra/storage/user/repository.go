package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/m04kA/SMC-UserService/internal/domain"
	userservice "github.com/m04kA/SMC-UserService/internal/service/user"
	"github.com/m04kA/SMC-UserService/pkg/psqlbuilder"
)

var (
	ErrCreateUser      = errors.New("failed to create user in database")
	ErrGetUser         = errors.New("failed to get user from database")
	ErrUpdateUser      = errors.New("failed to update user in database")
	ErrDeleteUser      = errors.New("failed to delete user from database")
	ErrGetSuperUsers   = errors.New("failed to get super users from database")
	ErrBuildQuery      = errors.New("failed to build SQL query")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(executor *sqlx.DB) *Repository {
	return &Repository{
		db: executor,
	}
}

// Create сохраняет нового пользователя в базу данных
func (r *Repository) Create(ctx context.Context, user *domain.User) error {
	query, args, err := psqlbuilder.Insert("users").
		Columns("tg_user_id", "name", "phone_number", "tg_link", "role_id", "created_at").
		Values(user.TGUserID, user.Name, user.PhoneNumber, user.TGLink, user.RoleID, user.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBuildQuery, err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCreateUser, err)
	}

	return nil
}

// GetByTGID находит пользователя по Telegram ID
func (r *Repository) GetByTGID(ctx context.Context, tgID int64) (*domain.User, error) {
	query, args, err := psqlbuilder.Select(
		"u.tg_user_id",
		"u.name",
		"u.phone_number",
		"u.tg_link",
		"u.role_id",
		"r.name as role_name",
		"u.created_at",
	).
		From("users u").
		LeftJoin("roles r ON u.role_id = r.id").
		Where(squirrel.Eq{"u.tg_user_id": tgID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBuildQuery, err)
	}

	var user domain.User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userservice.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrGetUser, err)
	}

	return &user, nil
}

// Update обновляет данные пользователя
func (r *Repository) Update(ctx context.Context, user *domain.User) error {
	query, args, err := psqlbuilder.Update("users").
		Set("name", user.Name).
		Set("phone_number", user.PhoneNumber).
		Set("tg_link", user.TGLink).
		Where(squirrel.Eq{"tg_user_id": user.TGUserID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBuildQuery, err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateUser, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: failed to get rows affected: %v", ErrUpdateUser, err)
	}

	if rowsAffected == 0 {
		return userservice.ErrUserNotFound
	}

	return nil
}

// Delete удаляет пользователя
func (r *Repository) Delete(ctx context.Context, tgID int64) error {
	query, args, err := psqlbuilder.Delete("users").
		Where(squirrel.Eq{"tg_user_id": tgID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBuildQuery, err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteUser, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: failed to get rows affected: %v", ErrDeleteUser, err)
	}

	if rowsAffected == 0 {
		return userservice.ErrUserNotFound
	}

	return nil
}

// GetSuperUsers возвращает список tg_user_id всех суперпользователей
func (r *Repository) GetSuperUsers(ctx context.Context) ([]int64, error) {
	query, args, err := psqlbuilder.Select("u.tg_user_id").
		From("users u").
		Where(squirrel.Eq{"u.role_id": domain.RoleIDSuperUser}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBuildQuery, err)
	}

	var userIDs []int64
	err = r.db.SelectContext(ctx, &userIDs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetSuperUsers, err)
	}

	return userIDs, nil
}
