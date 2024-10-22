package model

import (
	"database/sql"
	"time"
)

type Session struct {
	Id        string
	UserId    string
	ExpiresAt string
	CreatedAt string
	UpdatedAt sql.NullString
}

func (db *Database) GetSession(id string) (Session, error) {
	var session Session

	if err := db.conn.QueryRow(
		"SELECT id, userId, expiresAt, createdAt, updatedAt FROM sessions WHERE id = $1",
		id,
	).Scan(
		&session.Id,
		&session.UserId,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	); err != nil {
		return session, err
	}
	return session, nil
}

func (db *Database) GetSessionByUserId(userId string) (Session, error) {
	var session Session

	if err := db.conn.QueryRow(
		"SELECT id, userId, expiresAt, createdAt, updatedAt FROM sessions WHERE userId = $1",
		userId,
	).Scan(
		&session.Id,
		&session.UserId,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	); err != nil {
		return session, err
	}
	return session, nil
}

func getExpiresAtTimeStamp() string {
	return time.Now().Add(time.Minute * 10).Format(time.RFC3339)
}

func (db *Database) CreateUserSession(userId string) (string, error) {
	var id string

	if err := db.conn.QueryRow(
		"INSERT INTO sessions(userId, expiresAt, createdAt) VALUES($1, $2, $3) RETURNING id",
		userId,
		getExpiresAtTimeStamp(),
		getNowTimeStamp(),
	).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (db *Database) UpdateUserSession(id string) error {
	_, err := db.conn.Exec(
		`UPDATE sessions SET expiresAt = $1, updatedAt = $2 WHERE id = $3`,
		getExpiresAtTimeStamp(),
		getNowTimeStamp(),
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) InvalidateUserSession(id string) error {
	now := getNowTimeStamp()
	_, err := db.conn.Exec(
		`UPDATE sessions SET expiresAt = $1, updatedAt = $1 WHERE id = $2`,
		now,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}
