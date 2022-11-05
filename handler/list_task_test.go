package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MatsuoTakuro/go_todo_app/entity"
	"github.com/MatsuoTakuro/go_todo_app/testutil"
)

func TestListTask(t *testing.T) {
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		tasks []*entity.Task
		want  want
	}{
		"ok": {
			tasks: []*entity.Task{
				{
					ID:     1,
					Title:  "test1",
					Status: entity.TaskStatusTodo,
				},
				{
					ID:     2,
					Title:  "test2",
					Status: entity.TaskStatusDone,
				},
			},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_task/ok_rsp.json.golden",
			},
		},
		"empty": {
			tasks: []*entity.Task{},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_task/empty_rsp.json.golden",
			},
		},
	}
	for n, sub := range tests {
		sub := sub
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/tasks", nil)

			moq := &ListTasksServiceMock{}
			moq.ListTasksFunc = func(ctx context.Context) (entity.Tasks, error) {
				if sub.tasks != nil {
					return sub.tasks, nil
				}
				return nil, errors.New("error from mock")
			}
			handler := ListTask{Service: moq}
			handler.ServeHTTP(w, r)

			rsp := w.Result()
			testutil.AssertResponse(t,
				rsp, sub.want.status, testutil.LoadFile(t, sub.want.rspFile),
			)
		})
	}
}
