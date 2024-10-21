package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func GetDB() *sql.DB {
	return db
}

// DB can go into models and main but not controllers
// models can go into controllers not main

func Connect() {
	var err error
	db, err = sql.Open("postgres", "user=my-pass-go dbname=my-pass-go password=example host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func SetupTablesAndUser() {
	err := createUsersTable()
	if err != nil {
		log.Fatal(err)
	}

	err = createSessionsTable()
	if err != nil {
		log.Fatal(err)
	}

	// check if a user exists
	email := "drocktoo@gmail.com"
	var id string
	if err := db.QueryRow(
		"SELECT id FROM users WHERE email = $1",
		email,
	).Scan(
		&id,
	); err != nil {
		if err == sql.ErrNoRows {
			_, err = db.Exec(
				`INSERT INTO users(email, createdAt) VALUES($1, $2)`,
				email,
				time.Now().Format(time.RFC3339),
			)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		log.Fatal(err)
	}

	err = createPassTable()
	if err != nil {
		log.Fatal(err)
	}

	err = createTagsTable()
	if err != nil {
		log.Fatal(err)
	}

	err = createWebsitesTable()
	if err != nil {
		log.Fatal(err)
	}
}

func createUsersTable() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id uuid DEFAULT gen_random_uuid() UNIQUE,
		email varchar(45) NOT NULL UNIQUE,
		password varchar(450) NULL,
		createdAt varchar(45) NOT NULL,
		updatedAt varchar(45) NULL,
		PRIMARY KEY (id)
	)`)
	if err != nil {
		return err
	}

	return nil
}

func createSessionsTable() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id uuid DEFAULT gen_random_uuid() UNIQUE,
		userId uuid NOT NULL,
		expiresAt varchar(45) NOT NULL,
		createdAt varchar(45) NOT NULL,
		updatedAt varchar(45) NULL,
		PRIMARY KEY (id)
	)`)
	if err != nil {
		return err
	}

	return nil
}

func createPassTable() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS passes (
		id uuid DEFAULT gen_random_uuid() UNIQUE,
		userId uuid NOT NULL,
		name varchar(450) NULL,
		username varchar(450) NULL,
		password varchar(450) NULL,
		createdAt varchar(45) NOT NULL,
		updatedAt varchar(45) NULL,
		PRIMARY KEY (id)
	)`)
	if err != nil {
		return err
	}

	return nil
}

func createTagsTable() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tags (
		id uuid DEFAULT gen_random_uuid() UNIQUE,
		passId uuid NOT NULL,
		value varchar(450) NULL,
		createdAt varchar(45) NOT NULL,
		updatedAt varchar(45) NULL,
		PRIMARY KEY (id),
		CONSTRAINT fk_passes
      		FOREIGN KEY(passid) 
	  			REFERENCES passes(id)
	  			ON DELETE CASCADE
	)`)
	if err != nil {
		return err
	}

	return nil
}

func createWebsitesTable() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS websites (
		id uuid DEFAULT gen_random_uuid() UNIQUE,
		passId uuid NOT NULL,
		value varchar(450) NULL,
		createdAt varchar(45) NOT NULL,
		updatedAt varchar(45) NULL,
		PRIMARY KEY (id),
		CONSTRAINT fk_passes
      		FOREIGN KEY(passid) 
	  			REFERENCES passes(id)
	  			ON DELETE CASCADE
	)`)
	if err != nil {
		return err
	}

	return nil
}
