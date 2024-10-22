package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	return sql.Open("postgres", "user=my-pass-go dbname=my-pass-go password=example host=localhost port=5432 sslmode=disable")
}

func SetupTablesAndUser(db *sql.DB) {
	err := createUsersTable(db)
	if err != nil {
		log.Fatal(err)
	}

	err = createSessionsTable(db)
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

	err = createPassTable(db)
	if err != nil {
		log.Fatal(err)
	}

	err = createTagsTable(db)
	if err != nil {
		log.Fatal(err)
	}

	err = createWebsitesTable(db)
	if err != nil {
		log.Fatal(err)
	}
}

func createUsersTable(db *sql.DB) error {
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

func createSessionsTable(db *sql.DB) error {
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

func createPassTable(db *sql.DB) error {
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

func createTagsTable(db *sql.DB) error {
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

func createWebsitesTable(db *sql.DB) error {
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
