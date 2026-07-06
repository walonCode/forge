package tasks

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api/internals/modules/auth"

	"github.com/go-chi/chi/v5"
)

func testTaskHandler(r Repository) *Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return newHandler(newService(r), logger)
}

// taskRequest builds a request carrying a chi "id" path param and an
// authenticated user id, the way the router + AuthMiddleware would.
func taskRequest(method, taskID, userID, body string) *http.Request {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, "/", strings.NewReader(body))
	} else {
		req = httptest.NewRequest(method, "/", nil)
	}

	rctx := chi.NewRouteContext()
	if taskID != "" {
		rctx.URLParams.Add("id", taskID)
	}
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, auth.UserIdKey, userID)
	return req.WithContext(ctx)
}

func TestCreateTaskHandler(t *testing.T) {
	repo := &fakeTaskRepo{createID: "task-id"}
	h := testTaskHandler(repo)
	rr := httptest.NewRecorder()
	req := taskRequest(http.MethodPost, "", "user-1", `{"title":"buy milk","description":"two liters"}`)

	h.CreateTask(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rr.Code)
	}
}

func TestCreateTaskHandler_ValidationRejectsShortFields(t *testing.T) {
	h := testTaskHandler(&fakeTaskRepo{createID: "task-id"})
	rr := httptest.NewRecorder()
	req := taskRequest(http.MethodPost, "", "user-1", `{"title":"t","description":"d"}`)

	h.CreateTask(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
}

func TestDeleteTaskHandler_PassesCorrectIDs(t *testing.T) {
	repo := &fakeTaskRepo{}
	h := testTaskHandler(repo)
	rr := httptest.NewRecorder()
	req := taskRequest(http.MethodDelete, "task-9", "user-1", "")

	h.DeleteTask(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rr.Code)
	}
	// guards against the user id / task id being swapped on the way to the repo
	if repo.gotUserID != "user-1" || repo.gotTaskID != "task-9" {
		t.Fatalf("repo got (userID=%q, taskID=%q), want (user-1, task-9)", repo.gotUserID, repo.gotTaskID)
	}
}

func TestUpdateTaskHandler_DecodesObjectBody(t *testing.T) {
	repo := &fakeTaskRepo{task: &Task{ID: "task-9", IsCompleted: true}}
	h := testTaskHandler(repo)
	rr := httptest.NewRecorder()
	req := taskRequest(http.MethodPatch, "task-9", "user-1", `{"isCompleted":true}`)

	h.UpdateTask(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if !repo.gotComplete {
		t.Fatal("expected isCompleted=true to reach the repo")
	}
	if repo.gotUserID != "user-1" || repo.gotTaskID != "task-9" {
		t.Fatalf("repo got (userID=%q, taskID=%q), want (user-1, task-9)", repo.gotUserID, repo.gotTaskID)
	}
}

func TestUpdateTaskHandler_MissingID(t *testing.T) {
	h := testTaskHandler(&fakeTaskRepo{})
	rr := httptest.NewRecorder()
	req := taskRequest(http.MethodPatch, "", "user-1", `{"isCompleted":true}`)

	h.UpdateTask(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
}
