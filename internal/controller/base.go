package controller

import (
	"context"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/dhaef/my-pass/internal/model"
)

type key string

const sessionCookieName = "sessionId"
const userIdKeyName key = "userId"

type Base struct {
	Data any
}

func buildTemplatePaths(files []string) []string {
	ex, err := os.Executable()
	if err != nil {
		log.Println(err)
		return []string{}
	}
	exPath := filepath.Dir(ex)
	templatesPath := path.Join(exPath, "..", "..", "templates")

	for index, file := range files {
		files[index] = path.Join(templatesPath, file)
	}

	return files
}

func handleTemplateFiles(files []string) (*template.Template, error) {
	filesWithFullPath := buildTemplatePaths(files)

	return template.ParseFiles(filesWithFullPath...)
}

func render(w http.ResponseWriter, data any, files []string) error {
	t, _ := handleTemplateFiles(files)
	return t.Execute(w, data)

}
func renderTemplate(w http.ResponseWriter, data any, name string, files []string) error {
	t, _ := handleTemplateFiles(files)
	return t.ExecuteTemplate(w, name, data)
}

func setCookie(w http.ResponseWriter, name string, value string) {
	cookie := http.Cookie{
		Secure:   true,
		HttpOnly: true,
		Name:     name,
		Value:    value,
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
}

type APIError struct {
	Status       int
	Message      string
	ResponseType string
}

func (e APIError) Error() string {
	return e.Message
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func MakeHandler(h apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// log some data
		uri := r.RequestURI
		method := r.Method
		log.Printf("Incoming request: uri=%s method=%s", uri, method)

		if err := h(w, r); err != nil {
			if e, ok := err.(APIError); ok {
				log.Println(e.Error())

				if e.ResponseType == "JSON" {
					encode(w, r, e.Status, map[string]string{"message": e.Error()})
					return
				}

				renderTemplate(w, "", "layout", []string{"not-found.html", "layout.html"})
			}
		}

		duration := time.Since(start)
		log.Printf("Handled request: uri=%s method=%s duration=%s", uri, method, duration)
	}
}

func WithAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		db := model.GetDBFromCtx(r)
		session, err := db.GetSession(cookie.Value)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if time.Now().After(expiresAt) {
			log.Println("Session is expired")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		err = db.UpdateUserSession(session.Id)
		if err != nil {
			log.Println("failed to dump session expires time", err)
		}

		r = r.WithContext(context.WithValue(r.Context(), userIdKeyName, session.UserId))
		h(w, r)
	}
}
