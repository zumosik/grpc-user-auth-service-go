package storage

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrNothingReturned = errors.New("nothing is returned")
)

type User struct {
	ID        uint      `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

func (u User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {

		return err
	}
	u.Password = string(hashed)
	return nil
}
