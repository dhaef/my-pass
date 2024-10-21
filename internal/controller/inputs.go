package controller

import (
	"net/http"

	"github.com/google/uuid"
)

func AddMultipleInputHtml(w http.ResponseWriter, r *http.Request) error {
	id := uuid.NewString()
	inputType := r.PathValue("type")
	return render(
		w,
		Base{
			Data: map[string]string{
				"id":   id,
				"type": inputType,
			},
		},
		[]string{"inputs/add-multiple.html"},
	)
}
