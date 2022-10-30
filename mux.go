package main

import (
	"context"
	"net/http"

	"github.com/MatsuoTakuro/go_todo_app/clock"
	"github.com/MatsuoTakuro/go_todo_app/config"
	"github.com/MatsuoTakuro/go_todo_app/handler"
	"github.com/MatsuoTakuro/go_todo_app/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()

	// curl -i -XGET localhost:18000/health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	r := store.Repository{Clocker: clock.RealClocker{}}
	at := &handler.AddTask{
		DB:        db,
		Repo:      &r,
		Validator: v,
	}
	// curl -i -XPOST localhost:18000/tasks -d @./handler/testdata/add_task/ok_req.json.golden
	// curl -i -XPOST localhost:18000/tasks -d @./handler/testdata/add_task/bad_req.json.golden
	mux.Post("/tasks", at.ServeHTTP)

	lt := &handler.ListTask{
		DB:   db,
		Repo: &r,
	}
	// curl -i -XGET localhost:18000/tasks
	mux.Get("/tasks", lt.ServeHTTP)

	return mux, cleanup, nil
}
