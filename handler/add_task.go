package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MatsuoTakuro/go_todo_app/entity"
	"github.com/MatsuoTakuro/go_todo_app/store"
	"github.com/go-playground/validator/v10"
)

type AddTask struct {
	Store     *store.TaskStore
	validator *validator.Validate
}

func (at *AddTask) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body struct {
		Title string `json:"title" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	err := at.validator.Struct(body)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	t := &entity.Task{
		Title:   body.Title,
		Status:  entity.TaskStatusTodo,
		Created: time.Now(),
	}
	id, err := store.Tasks.Add(t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
	}

	rsp := struct {
		ID int `json:"id"`
	}{ID: id}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
