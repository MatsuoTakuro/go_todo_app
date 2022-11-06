package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MatsuoTakuro/go_todo_app/config"
)

var _ http.ResponseWriter = (*httptest.ResponseRecorder)(nil)

func TestNewMux(t *testing.T) {
	w := httptest.NewRecorder() // http.ResponseWriter implements http.ResponseWriter
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	cfg, _ := config.New()
	handler, _, _ := NewMux(context.Background(), cfg)
	handler.ServeHTTP(w, r)
	resp := w.Result()
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	if resp.StatusCode != http.StatusOK {
		t.Error("want status code 200, but", resp.StatusCode)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := `{"status": "ok"}`
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}
