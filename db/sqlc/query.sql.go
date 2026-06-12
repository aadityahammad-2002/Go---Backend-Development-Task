package sqlc

import (
	"database/sql"
	"fmt"
	"time"
)

type Queries struct {
	db *sql.DB
}

func New(db *sql.DB) *Queries {
	return &Queries{db: db}
}

func (q *Queries) CreateUser(name string, dob time.Time) (User, error) {
	result, err := q.db.Exec("INSERT INTO users (name, dob) VALUES (?, ?)", name, dob)
	if err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return User{}, fmt.Errorf("retrieve last insert id: %w", err)
	}

	return User{ID: id, Name: name, Dob: dob}, nil
}

func (q *Queries) GetUserByID(id int64) (User, error) {
	row := q.db.QueryRow("SELECT id, name, dob FROM users WHERE id = ?", id)
	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Dob); err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user not found: %w", err)
		}
		return User{}, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (q *Queries) UpdateUser(id int64, name string, dob time.Time) (User, error) {
	result, err := q.db.Exec("UPDATE users SET name = ?, dob = ? WHERE id = ?", name, dob, id)
	if err != nil {
		return User{}, fmt.Errorf("update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return User{}, fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return User{}, sql.ErrNoRows
	}

	return User{ID: id, Name: name, Dob: dob}, nil
}

func (q *Queries) DeleteUser(id int64) error {
	result, err := q.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (q *Queries) ListUsers(limit, offset int64) ([]User, error) {
	rows, err := q.db.Query("SELECT id, name, dob FROM users ORDER BY id ASC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Dob); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}
