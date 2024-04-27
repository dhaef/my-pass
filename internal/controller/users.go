package controller

import (
	"fmt"
	"net/http"

	"github.com/dhaef/my-pass/internal/model"
)

func getUsers() ([]model.User, error) {
	users, err := model.GetUsers()
	if err != nil {
		return users, nil
	}

	return users, nil
}

func UsersJson(w http.ResponseWriter, r *http.Request) error {
	users, err := getUsers()
	if err != nil {
		return APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return encode(w, r, http.StatusOK, users)
}

func UsersHtml(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("UserId: ", r.Context().Value(userIdKeyName))
	users, err := getUsers()
	if err != nil {
		return APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return renderTemplate(w, base[map[string]any]{
		Data: map[string]any{
			"users": users,
		},
	}, "layout", []string{"users/users.html", "layout.html"})
}
