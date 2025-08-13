package repository

import (
	"context"
	"database/sql"
	"pisondev/markdown-notes-api/model/domain"

	"github.com/sirupsen/logrus"
)

type UserRepositoryImpl struct {
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) UserRepository {
	return &UserRepositoryImpl{
		Log: log,
	}
}

func (r *UserRepositoryImpl) Register(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {
	SQL := "INSERT INTO users(username, hashed_password, created_at) VALUES (?,?,?)"
	result, err := tx.ExecContext(ctx, SQL, user.Username, user.HashedPassword, user.CreatedAt)
	if err != nil {
		r.Log.Errorf("failed to execute query: %v", err)
		return domain.User{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		r.Log.Errorf("failed to get last insert id: %v", err)
		return domain.User{}, err
	}

	user.ID = int(id)

	return user, nil
}

func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, tx *sql.Tx, username string) (domain.User, error) {
	SQL := "SELECT * FROM users WHERE username = ?"
	rows, err := tx.QueryContext(ctx, SQL, username)
	if err != nil {
		r.Log.Errorf("failed to execute query: %v", err)
		return domain.User{}, err
	}
	defer rows.Close()

	user := domain.User{}
	if rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.HashedPassword,
			&user.CreatedAt)
		if err != nil {
			r.Log.Errorf("failed to scan selected user: %v", err)
			return domain.User{}, err
		}
	}
	return user, nil
}
