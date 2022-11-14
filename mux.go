package main

import (
	"context"
	"net/http"

	"github.com/MatsuoTakuro/go_todo_app/auth"
	"github.com/MatsuoTakuro/go_todo_app/clock"
	"github.com/MatsuoTakuro/go_todo_app/config"
	"github.com/MatsuoTakuro/go_todo_app/handler"
	"github.com/MatsuoTakuro/go_todo_app/service"
	"github.com/MatsuoTakuro/go_todo_app/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()
	// curl -i -X GET localhost:18000/health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	clocker := clock.RealClocker{}
	repo := store.Repository{Clocker: clocker}
	kvs, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	jwter, err := auth.NewJWTer(kvs, clocker)
	if err != nil {
		return nil, cleanup, err
	}

	ru := &handler.RegisterUser{
		Service: &service.RegisterUser{
			DB:   db,
			Repo: &repo,
		},
		Validator: v,
	}
	// curl -X POST localhost:18000/register -d '{"name": "budou2", "password":"test", "role":"admin"}'
	mux.Post("/register", ru.ServeHTTP)

	l := &handler.Login{
		Service: &service.Login{
			DB:             db,
			Repo:           &repo,
			TokenGenerator: jwter,
		},
		Validator: v,
	}
	// curl -X POST localhost:18000/login -d '{"user_name": "john", "password": "test"}' | jq
	// https://jwt.io/
	mux.Post("/login", l.ServeHTTP)

	at := &handler.AddTask{
		Service: &service.AddTask{
			DB:   db,
			Repo: &repo,
		},
		Validator: v,
	}
	lt := &handler.ListTask{
		Service: &service.ListTask{
			DB:   db,
			Repo: &repo,
		},
	}
	mux.Route("/tasks", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter))
		// curl -i -X POST localhost:18000/tasks -d @./handler/testdata/add_task/ok_req.json.golden
		// curl -i -X POST localhost:18000/tasks -d @./handler/testdata/add_task/bad_req.json.golden
		r.Post("/", at.ServeHTTP)
		// curl -i -X GET localhost:18000/tasks
		r.Get("/", lt.ServeHTTP)
	})

	mux.Route("/admin", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter), handler.AdminMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"message": "admin only"}`))
		})
	})

	return mux, cleanup, nil
}
