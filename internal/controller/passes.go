package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/dhaef/my-pass/internal/model"
)

func PassHtml(w http.ResponseWriter, r *http.Request) error {
	userId := r.Context().Value(userIdKeyName).(string)
	id := r.PathValue("id")

	pass, err := model.GetPass(id, userId)
	if err != nil {
		log.Println(err)
		return APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return renderTemplate(w, base[map[string]any]{
		Data: map[string]any{
			"pass": pass,
		},
	}, "layout", []string{"passes/pass.html", "layout.html"})
}

func PassesHtml(w http.ResponseWriter, r *http.Request) error {
	userId := r.Context().Value(userIdKeyName).(string)
	passes, err := model.GetPasses(userId)
	if err != nil {
		log.Println(err)
		return APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	// log.Printf("%+v\n", passes)
	return renderTemplate(w, base[map[string]any]{
		Data: map[string]any{
			"passes": passes,
		},
	}, "layout", []string{"passes/passes.html", "layout.html"})
}

func CreatePassHtml(w http.ResponseWriter, r *http.Request) error {
	userId := r.Context().Value(userIdKeyName).(string)
	id := r.PathValue("id")

	pass, err := model.GetPass(id, userId)
	if err != nil {
		log.Println(err)
		log.Println(pass)
	}

	return renderTemplate(
		w,
		base[map[string]any]{
			Data: map[string]any{
				"pass": pass,
			},
		},
		"layout",
		[]string{"passes/pass-form.html", "layout.html"},
	)
}

func getPassFromForm(r *http.Request) *model.Pass {
	r.ParseForm()

	userId := r.Context().Value(userIdKeyName).(string)

	pass := model.Pass{
		UserId: userId,
		Name: sql.NullString{
			String: r.FormValue("name"),
			Valid:  true,
		},
		Username: sql.NullString{
			String: r.FormValue("username"),
			Valid:  true,
		},
		Password: sql.NullString{
			String: r.FormValue("password"),
			Valid:  true,
		},
		Tags:     []model.PassItem{},
		Websites: []model.PassItem{},
	}

	for key, value := range r.Form {
		if strings.HasPrefix(key, "tag") {
			pass.Tags = append(pass.Tags, model.PassItem{
				Value: sql.NullString{
					String: value[0],
					Valid:  true,
				},
				Id: sql.NullString{
					String: strings.Join(strings.Split(key, "-")[1:], "-"),
					Valid:  true,
				},
			})
			continue
		}

		if strings.HasPrefix(key, "website") {
			pass.Websites = append(pass.Websites, model.PassItem{
				Value: sql.NullString{
					String: value[0],
					Valid:  true,
				},
				Id: sql.NullString{
					String: strings.Join(strings.Split(key, "-")[1:], "-"),
					Valid:  true,
				},
			})
			continue
		}
	}

	return &pass
}

func CreatePass(w http.ResponseWriter, r *http.Request) error {
	pass := getPassFromForm(r)

	pass, err := model.CreatePass(pass)
	if err != nil {
		log.Println(err)
	}

	w.Header().Add("HX-Redirect", "/passes/"+pass.Id)
	return nil
}

func UpdatePass(w http.ResponseWriter, r *http.Request) error {
	pass := getPassFromForm(r)
	id := r.PathValue("id")
	pass.Id = id

	pass, err := model.UpdatePass(pass)
	if err != nil {
		log.Println(err)
	}

	w.Header().Add("HX-Redirect", "/passes/"+pass.Id)
	return nil
}

func DeletePass(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	err := model.DeletePass(id)
	if err != nil {
		log.Println(err)
		return APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	w.Header().Add("HX-Redirect", "/passes")
	return nil
}
