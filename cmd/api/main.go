package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/dhaef/my-pass/internal/db"
	"github.com/dhaef/my-pass/internal/jobs"
)

func main() {
	db.Connect()
	db.SetupTablesAndUser()

	jobsChan, resultsChan := jobs.Start(3)
	s := http.Server{
		Addr:    ":3000",
		Handler: buildRoutes(),
		BaseContext: func(l net.Listener) context.Context {
			return context.WithValue(context.Background(), "jobsChan", jobsChan)
		},
	}
	go func() {
		for r := range resultsChan {
			fmt.Println(r)
		}
	}()

	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		jobs.Stop(jobsChan, resultsChan)
	}
}
