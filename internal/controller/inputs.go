package controller

import (
	"net/http"

	"github.com/google/uuid"
)

func AddMultipleInputHtml(w http.ResponseWriter, r *http.Request) error {
	id := uuid.New()
	inputType := r.PathValue("type")
	return render(
		w,
		base[map[string]any]{
			Data: map[string]any{
				"id":   id,
				"type": inputType,
			},
		},
		[]string{"inputs/add-multiple.html"},
	)
}
