package controller

import (
	"fmt"
	"net/http"

	"github.com/dhaef/my-pass/internal/model"
)

func getUsers(r *http.Request) ([]model.User, error) {
	db := model.GetDBFromCtx(r)
	users, err := db.GetUsers()
	if err != nil {
		return users, nil
	}

	return users, nil
}

func UsersJson(w http.ResponseWriter, r *http.Request) error {
	users, err := getUsers(r)
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
	users, err := getUsers(r)
	if err != nil {
		return APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return renderTemplate(w, Base{
		Data: map[string]any{
			"users": users,
		},
	}, "layout", []string{"users/users.html", "layout.html"})
}
