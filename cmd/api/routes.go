package main

import (
	"net/http"

	c "github.com/dhaef/my-pass/internal/controller"
)

func buildRoutes() http.Handler {
	r := http.NewServeMux()

	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../static"))))

	r.HandleFunc("GET /api/users", c.WithAuth(c.MakeHandler(c.UsersJson)))
	r.HandleFunc("GET /users", c.WithAuth(c.MakeHandler(c.UsersHtml)))
	r.HandleFunc("POST /logout", c.WithAuth(c.MakeHandler(c.LogoutRedirect)))
	r.HandleFunc("GET /passes", c.WithAuth(c.MakeHandler(c.PassesHtml)))
	r.HandleFunc("GET /passes/{id}", c.WithAuth(c.MakeHandler(c.PassHtml)))
	r.HandleFunc("GET /passes/create", c.WithAuth(c.MakeHandler(c.CreatePassHtml)))
	r.HandleFunc("GET /passes/create/{id}", c.WithAuth(c.MakeHandler(c.CreatePassHtml)))
	r.HandleFunc("POST /passes/create", c.WithAuth(c.MakeHandler(c.CreatePass)))
	r.HandleFunc("POST /passes/update/{id}", c.WithAuth(c.MakeHandler(c.UpdatePass)))
	r.HandleFunc("POST /passes/delete/{id}", c.WithAuth(c.MakeHandler(c.DeletePass)))
	r.HandleFunc("GET /inputs/{type}", c.WithAuth(c.MakeHandler(c.AddMultipleInputHtml)))

	r.HandleFunc("GET /login", c.MakeHandler(c.LoginHtml))
	r.HandleFunc("POST /api/login", c.MakeHandler(c.LoginJson))
	r.HandleFunc("GET /sign-in", c.MakeHandler(c.LoginHtml))
	r.HandleFunc("GET /sign-up", c.MakeHandler(c.SignUpHtml))
	r.HandleFunc("POST /api/sign-up", c.MakeHandler(c.SignUpJson))
	r.HandleFunc("GET /not-found", c.MakeHandler(c.NotFoundHtml))
	r.HandleFunc("GET /job", c.MakeHandler(c.SubmitJob))

	return r
}
