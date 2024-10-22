package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/dhaef/my-pass/internal/db"
	"github.com/dhaef/my-pass/internal/jobs"
)

func main() {
	conn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	db.SetupTablesAndUser(conn)

	jobsChan, resultsChan := jobs.Start(3)
	s := http.Server{
		Addr:    ":3000",
		Handler: buildRoutes(),
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.WithValue(context.Background(), "dbConn", conn)
			return context.WithValue(ctx, "jobsChan", jobsChan)
		},
	}
	go func() {
		for r := range resultsChan {
			fmt.Println(r)
		}
	}()

	err = s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		jobs.Stop(jobsChan, resultsChan)
	}
}
