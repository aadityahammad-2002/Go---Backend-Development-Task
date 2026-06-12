package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/yourname/user-api/internal/logger"
	"go.uber.org/zap"
)

type User struct {
	ID   int32
	Name string
	DOB  time.Time
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, name string, dob time.Time) (int32, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO users (name, dob) VALUES (?, ?)", name, dob)
	if err != nil {
		logger.Logger.Error("failed to create user", zap.Error(err))
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Logger.Error("failed to get last insert id", zap.Error(err))
		return 0, err
	}

	logger.Logger.Info("user created", zap.Int64("user_id", id), zap.String("name", name))
	return int32(id), nil
}

func (r *UserRepository) GetUser(ctx context.Context, id int32) (*User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, dob FROM users WHERE id = ? LIMIT 1", id)
	
	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.DOB)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logger.Info("user not found", zap.Int32("user_id", id))
			return nil, sql.ErrNoRows
		}
		logger.Logger.Error("failed to get user", zap.Int32("user_id", id), zap.Error(err))
		return nil, err
	}

	logger.Logger.Info("user retrieved", zap.Int32("user_id", id))
	return user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context, limit, offset int64) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, dob FROM users ORDER BY id LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		logger.Logger.Error("failed to get all users", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.DOB); err != nil {
			logger.Logger.Error("failed to scan user", zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		logger.Logger.Error("error iterating users", zap.Error(err))
		return nil, err
	}

	logger.Logger.Info("users retrieved", zap.Int("count", len(users)))
	return users, nil
}

func (r *UserRepository) GetUserCount(ctx context.Context) (int64, error) {
	row := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users")
	
	var count int64
	err := row.Scan(&count)
	if err != nil {
		logger.Logger.Error("failed to get user count", zap.Error(err))
		return 0, err
	}

	return count, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id int32, name string, dob time.Time) error {
	result, err := r.db.ExecContext(ctx, "UPDATE users SET name = ?, dob = ? WHERE id = ?", name, dob, id)
	if err != nil {
		logger.Logger.Error("failed to update user", zap.Int32("user_id", id), zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Logger.Error("failed to get rows affected", zap.Int32("user_id", id), zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		logger.Logger.Info("user not found for update", zap.Int32("user_id", id))
		return sql.ErrNoRows
	}

	logger.Logger.Info("user updated", zap.Int32("user_id", id), zap.String("name", name))
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int32) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		logger.Logger.Error("failed to delete user", zap.Int32("user_id", id), zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Logger.Error("failed to get rows affected", zap.Int32("user_id", id), zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		logger.Logger.Info("user not found for delete", zap.Int32("user_id", id))
		return sql.ErrNoRows
	}

	logger.Logger.Info("user deleted", zap.Int32("user_id", id))
	return nil
}
