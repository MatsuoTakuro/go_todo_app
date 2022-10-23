package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertJSON(t *testing.T, want, got []byte) {
	t.Helper()

	var jsonWant, jsonGot any
	if err := json.Unmarshal(want, &jsonWant); err != nil {
		t.Fatalf("cannot unmarshal want %q: %v", want, err)
	}
	if err := json.Unmarshal(got, &jsonGot); err != nil {
		t.Fatalf("cannot unmarshal got %q: %v", got, err)
	}
	if diff := cmp.Diff(jsonGot, jsonWant); diff != "" {
		t.Errorf("got differs: (-got +want)\n%s", diff)
	}
}

func AsserrtReponse(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper()

	t.Cleanup(func() {
		_ = got.Body.Close()
	})
	gotBody, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}
	if got.StatusCode != status {
		t.Fatalf("want status %d, but got %d, body: %q", status, got.StatusCode, gotBody)
	}

	// When the response body is empty
	if len(gotBody) == 0 && len(body) == 0 {
		return
	}
	AssertJSON(t, body, gotBody)
}

func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read from %q: %v", path, err)
	}
	return content
}
