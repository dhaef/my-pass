package model

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        string         `json:"id"`
	Email     string         `json:"email"`
	Password  sql.NullString `json:"-"`
	CreatedAt string         `json"createdAt"`
	UpdatedAt sql.NullString `json"updatedAt"`
}

func (db *Database) GetUser(id string) (User, error) {
	var user User

	if err := db.conn.QueryRow(
		"SELECT id, email, createdAt, updatedAt from users where id = $1",
		id,
	).Scan(
		&user.Id,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return user, err
	}
	return user, nil
}

func (db *Database) GetUsers() ([]User, error) {
	rows, err := db.conn.Query("SELECT id, email, createdAt, updatedAt FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

func (db *Database) CreateUser(user User) (User, error) {
	var id string
	err := db.conn.QueryRow(
		`INSERT INTO users(email, createdAt) VALUES($1) RETURNING id`,
		user.Email,
		getNowTimeStamp(),
	).Scan(&id)
	if err != nil {
		return User{}, err
	}

	return User{
		Id:    id,
		Email: user.Email,
	}, nil
}

func (db *Database) UpdateUser(user *User) (*User, error) {
	err := db.conn.QueryRow(
		`UPDATE users SET email = $1, updatedAt = $2 WHERE id = $3`,
		user.Email,
		getNowTimeStamp(),
		user.Id,
	).Scan(user.Email)
	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func (db *Database) UpdateUserPassword(email string, password string) (string, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return "", err
	}

	var id string
	err = db.conn.QueryRow(
		`UPDATE users SET password = $1, updatedAt = $2 WHERE email = $3 RETURNING id`,
		hashedPassword,
		getNowTimeStamp(),
		email,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil

}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hashedPassword), err
}

func (db *Database) getUserByEmail(email string) (User, error) {
	var user User
	err := db.conn.QueryRow(
		`SELECT id, email, password FROM users WHERE email = $1`,
		email,
	).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *Database) AuthenticateUser(email string, password string) (string, error) {
	user, err := db.getUserByEmail(email)
	if err != nil {
		return "", err
	}
	if !user.Password.Valid {
		return "", errors.New("password can not be null")
	}
	return user.Id, bcrypt.CompareHashAndPassword(
		[]byte(user.Password.String),
		[]byte(password),
	)
}
