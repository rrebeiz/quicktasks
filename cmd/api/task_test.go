package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApp_GetTaskByIDHandler(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		expectedStatusCode int
		expectedBody       string
	}{
		{"valid test", "1", http.StatusOK, "{\"task\":{\"id\":1,\"title\":\"Test Title\",\"description\":\"Test Description\",\"completed\":false}}\n"},
		{"invalid test missing id", "", http.StatusNotFound, "{\"error\":\"the requested resource could not be found\"}\n"},
		{"invalid test invalid id", "abc", http.StatusNotFound, "{\"error\":\"the requested resource could not be found\"}\n"},
	}
	for _, e := range tests {
		req, _ := http.NewRequest("GET", "/v1/tasks/1", nil)
		ctxCtx := chi.NewRouteContext()
		ctxCtx.URLParams.Add("id", e.id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctxCtx))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(testApp.getTaskHandler)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if rr.Body.String() != e.expectedBody {
			t.Errorf("%s: expected %s but got %s", e.name, e.expectedBody, rr.Body.String())
		}
	}
}

func TestApp_CreateTaskHandler(t *testing.T) {
	tests := []struct {
		name               string
		body               string
		expectedStatusCode int
		expectedBody       string
	}{
		{"valid test", `{"title":"Test Task", "description":"Test Description", "completed":false}`, http.StatusCreated, "{\"task\":{\"id\":1,\"title\":\"Test Task\",\"description\":\"Test Description\",\"completed\":false}}\n"},
		{"invalid test: body set but empty", `{"title":"", "description":""}`, http.StatusUnprocessableEntity, "{\"error\":{\"description\":\"should not be empty\",\"title\":\"should not be empty\"}}\n"},
		{"invalid test: empty body", `{}`, http.StatusUnprocessableEntity, "{\"error\":{\"description\":\"should not be empty\",\"title\":\"should not be empty\"}}\n"},
	}
	for _, e := range tests {
		req, _ := http.NewRequest("POST", "/v1/tasks", strings.NewReader(e.body))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(testApp.createTaskHandler)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedBody != rr.Body.String() {
			t.Errorf("%s: expected %s but got %s", e.name, e.expectedBody, rr.Body.String())
		}
	}
}

func TestApp_UpdateTaskHandler(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		body               string
		expectedStatusCode int
		expectedBody       string
	}{
		{"valid test", "1", `{"title":"New Test", "description":"New Description", "completed":true}`, http.StatusOK, "{\"task\":{\"id\":1,\"title\":\"New Test\",\"description\":\"New Description\",\"completed\":true}}\n"},
		{"invalid test: body set but empty", "1", `{"title":"", "description":""}`, http.StatusUnprocessableEntity, "{\"error\":{\"description\":\"should not be empty\",\"title\":\"should not be empty\"}}\n"},
	}

	for _, e := range tests {
		req, _ := http.NewRequest("PATCH", "/v1/tasks/1", strings.NewReader(e.body))
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", e.id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(testApp.updateTaskHandler)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
		if e.expectedBody != rr.Body.String() {
			t.Errorf("%s: expected %s but got %s", e.name, e.expectedBody, rr.Body.String())
		}
	}
}

func TestApp_DeleteTaskHandler(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		expectedStatusCode int
		expectedBody       string
	}{
		{"valid test", "1", http.StatusOK, "{\"success\":\"task with ID 1 deleted\"}\n"},
		{"invalid test missing id", "", http.StatusNotFound, "{\"error\":\"the requested resource could not be found\"}\n"},
		{"invalid test invalid id", "abc", http.StatusNotFound, "{\"error\":\"the requested resource could not be found\"}\n"},
	}
	for _, e := range tests {
		req, _ := http.NewRequest("DELETE", "/v1/tasks/1", nil)
		//if e.id != "" {
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", e.id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		//}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(testApp.deleteTaskHandler)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedBody != rr.Body.String() {
			t.Errorf("%s: expected %s but got %s", e.name, e.expectedBody, rr.Body.String())
		}
	}
}

func TestApp_GetAllTasks(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
		expectedStatusBody string
	}{
		{"valid test", http.StatusOK, "{\"tasks\":[{\"id\":1,\"title\":\"Test Task\",\"description\":\"Test Description\",\"completed\":false}]}\n"},
	}
	for _, e := range tests {
		req, _ := http.NewRequest("GET", "/v1/tasks", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(testApp.getAllTasks)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
		if e.expectedStatusBody != rr.Body.String() {
			t.Errorf("%s: expected %s but got %s", e.name, e.expectedStatusBody, rr.Body.String())
		}
	}
}
