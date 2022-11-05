package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MatsuoTakuro/go_todo_app/entity"
	"github.com/MatsuoTakuro/go_todo_app/testutil"
	"github.com/go-playground/validator/v10"
)

func TestAddTask(t *testing.T) {
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_task/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_rsp.json.golden",
			},
		},
	}
	for n, sub := range tests {
		sub := sub
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, sub.reqFile)),
			)
			moq := &AddTaskServiceMock{}
			moq.AddTaskFunc = func(
				ctx context.Context, title string,
			) (*entity.Task, error) {
				if sub.want.status == http.StatusOK {
					return &entity.Task{ID: 1}, nil
				}
				return nil, errors.New("error from mock")
			}

			handler := AddTask{
				Service:   moq,
				Validator: validator.New(),
			}
			handler.ServeHTTP(w, r)

			rsp := w.Result()
			testutil.AssertResponse(t,
				rsp, sub.want.status, testutil.LoadFile(t, sub.want.rspFile),
			)
		})
	}
}
