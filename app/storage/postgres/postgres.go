package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zumosik/grpc-user-auth-service-go/storage"
)

type Storage struct {
	db *sqlx.DB
}

type ConnectionData struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"name" env-required:"true"`
}

func New(d ConnectionData) (*Storage, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", d.Host, d.Port, d.User, d.Password, d.DBName))

	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

// Close is function that will just close database connection
func (st *Storage) Close() error {
	return st.db.Close()
}

func (st *Storage) CreateUser(ctx context.Context, u storage.User) (uint, error) {
	insertQuery := `
        INSERT INTO users (username, email, password, created_at)
        VALUES (:username, :email, :password, NOW())
        RETURNING id`

	var userID uint
	err := st.db.QueryRowxContext(ctx, insertQuery, u).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (st *Storage) UpdateUser(ctx context.Context, u storage.User) error {
	updateQuery := `
        UPDATE users
        SET username = :username, email = :email, password = :password
        WHERE id = :id`

	_, err := st.db.NamedExecContext(ctx, updateQuery, u)
	if err != nil {
		return err
	}

	return nil
}

func (st *Storage) GetUsers(ctx context.Context, limit uint) ([]storage.User, error) {
	var users []storage.User
	query := `
        SELECT * FROM users
        LIMIT $1`

	err := st.db.SelectContext(ctx, &users, query, limit)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (st *Storage) GetUserByID(ctx context.Context, id uint) (storage.User, error) {
	var u storage.User
	query := "SELECT * FROM users WHERE id = $1"

	err := st.db.GetContext(ctx, &u, query, id)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			// if nothing is returned
			return storage.User{}, storage.ErrNothingReturned
		}

		return storage.User{}, err
	}

	return u, nil
}
