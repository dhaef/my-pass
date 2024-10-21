package controller

import (
	"log"
	"net/http"

	"github.com/dhaef/my-pass/internal/jobs"
	"github.com/dhaef/my-pass/internal/model"
	"github.com/google/uuid"
)

func LoginHtml(w http.ResponseWriter, r *http.Request) error {
	return render(w, "", []string{"login.html"})
}

func SignUpHtml(w http.ResponseWriter, r *http.Request) error {
	return render(w, "", []string{"sign-up.html"})
}

func NotFoundHtml(w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(w, "", "layout", []string{"layout.html", "not-found.html"})
}

func LogoutRedirect(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		log.Println(err)
	}

	newCookie := http.Cookie{
		Secure:   true,
		HttpOnly: true,
		Name:     sessionCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
	}

	http.SetCookie(w, &newCookie)

	err = model.InvalidateUserSession(cookie.Value)
	if err != nil {
		log.Println(err)
	}

	w.Header().Add("HX-Redirect", "/login")
	return nil
}

type Auth struct {
	Email    string
	Password string
}

func getAuthAndValidate(r *http.Request) (Auth, error) {
	auth, err := decode[Auth](r)
	if err != nil {
		return Auth{}, APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	if auth.Email == "" || auth.Password == "" {
		return Auth{}, APIError{
			Status:       http.StatusBadRequest,
			Message:      "email and password are required",
			ResponseType: "JSON",
		}
	}

	return auth, nil
}

func handleCreateSessionAndSetCookie(w http.ResponseWriter, r *http.Request, userId string) error {
	sessionId, err := model.CreateUserSession(userId)
	if err != nil {
		return APIError{
			Status:       http.StatusBadRequest,
			Message:      err.Error(),
			ResponseType: "JSON",
		}
	}

	setCookie(w, sessionCookieName, sessionId)
	return encode(w, r, http.StatusOK, map[string]string{"sessionId": sessionId})
}

func LoginJson(w http.ResponseWriter, r *http.Request) error {
	auth, err := getAuthAndValidate(r)
	if err != nil {
		return err
	}

	userId, err := model.AuthenticateUser(auth.Email, auth.Password)
	if err != nil {
		return APIError{
			Status:       http.StatusUnauthorized,
			Message:      "unauthorized",
			ResponseType: "JSON",
		}
	}

	return handleCreateSessionAndSetCookie(w, r, userId)
}

func SignUpJson(w http.ResponseWriter, r *http.Request) error {
	auth, err := decode[Auth](r)
	if err != nil {
		return APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	if auth.Email == "" || auth.Password == "" {
		return APIError{
			Status:       http.StatusBadRequest,
			Message:      "email and password are required",
			ResponseType: "JSON",
		}
	}

	userId, err := model.UpdateUserPassword(auth.Email, auth.Password)
	if err != nil {
		return APIError{
			Status:       http.StatusBadRequest,
			Message:      err.Error(),
			ResponseType: "JSON",
		}
	}

	return handleCreateSessionAndSetCookie(w, r, userId)
}

func SubmitJob(w http.ResponseWriter, r *http.Request) error {
	jobsChan := r.Context().Value("jobsChan").(chan jobs.Job)

	jobsChan <- jobs.Job{
		Id: uuid.NewString(),
		Handler: func(data map[string]string) (string, error) {
			return data["test"], nil
		},
		Data: map[string]string{"test": "this is a test map"},
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
