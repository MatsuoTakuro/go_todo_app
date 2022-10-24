package main

import (
	"net/http"

	"github.com/MatsuoTakuro/go_todo_app/handler"
	"github.com/MatsuoTakuro/go_todo_app/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func NewMux() http.Handler {
	mux := chi.NewRouter()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()
	at := &handler.AddTask{
		Store:     store.Tasks,
		Validator: v,
	}
	// curl -i -XPOST localhost:18000/tasks -d @./handler/testdata/add_task/ok_req.json.golden
	// curl -i -XPOST localhost:18000/tasks -d @./handler/testdata/add_task/bad_req.json.golden
	mux.Post("/tasks", at.ServerHTTP)

	lt := &handler.ListTask{
		Store: store.Tasks,
	}
	// curl -i -XGET localhost:18000/tasks
	mux.Get("/tasks", lt.ServerHTTP)

	return mux
}
