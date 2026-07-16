package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go_final_project/internal/nextdate"
	"go_final_project/pkg/db"
)

type apiTaskRequest struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type apiTaskResponse map[string]any

func RegisterHandlers() {
	http.HandleFunc("/api/signin", handleSignIn)
	http.HandleFunc("/api/nextdate", auth(handleNextDate))
	http.HandleFunc("/api/task", auth(handleTask))
	http.HandleFunc("/api/tasks", auth(handleTasks))
	http.HandleFunc("/api/task/done", auth(handleTaskDone))
}

func handleNextDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query()
	now := query.Get("now")
	date := query.Get("date")
	repeat := query.Get("repeat")

	if date == "" || repeat == "" {
		respondError(w, http.StatusBadRequest, "missing required query parameters")
		return
	}

	var parsedNow time.Time
	var err error
	if now == "" {
		parsedNow = time.Now()
	} else {
		parsedNow, err = time.Parse("20060102", now)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	next, err := nextdate.NextDate(parsedNow, date, repeat)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, next)
}

func handleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleTaskCreate(w, r)
	case http.MethodPut:
		handleTaskUpdate(w, r)
	case http.MethodGet:
		handleTaskGet(w, r)
	case http.MethodDelete:
		handleTaskDelete(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	var req apiTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Title == "" {
		respondError(w, http.StatusBadRequest, "title is required")
		return
	}
	date, err := normalizeTaskDate(req.Date, req.Repeat)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.Date = date

	task := db.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}
	id, err := db.AddTask(task)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, apiTaskResponse{"id": fmt.Sprint(id)})
}

func handleTaskGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}
	parsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	task, err := db.GetTask(parsed)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "task not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, apiTaskResponse{
		"id":      fmt.Sprint(task.ID),
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	})
}

func handleTaskUpdate(w http.ResponseWriter, r *http.Request) {
	var req apiTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ID == "" {
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}
	parsedID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if req.Title == "" {
		respondError(w, http.StatusBadRequest, "title is required")
		return
	}
	date, err := normalizeTaskDate(req.Date, req.Repeat)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.Date = date

	task := db.Task{
		ID:      parsedID,
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}
	if err := db.UpdateTask(task); err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "task not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, apiTaskResponse{})
}

func handleTaskDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}
	parsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := db.DeleteTask(parsed); err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "task not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, apiTaskResponse{})
}

func handleTaskDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}
	parsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	task, err := db.GetTask(parsed)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "task not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if task.Repeat == "" {
		if err := db.DeleteTask(parsed); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, apiTaskResponse{})
		return
	}
	now := time.Now()
	next, err := nextdate.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	task.Date = next
	if err := db.UpdateTask(*task); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, apiTaskResponse{})
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	search := r.URL.Query().Get("search")
	tasks, err := db.ListTasks(search)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseTasks := make([]map[string]any, 0, len(tasks))
	for _, task := range tasks {
		responseTasks = append(responseTasks, map[string]any{
			"id":      fmt.Sprint(task.ID),
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		})
	}
	response := map[string]any{"tasks": responseTasks}
	respondJSON(w, response)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(apiTaskResponse{"error": message})
}

func respondJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}

func normalizeTaskDate(dateValue, repeat string) (string, error) {
	if dateValue == "" {
		return time.Now().Format("20060102"), nil
	}

	parsedDate, err := time.Parse("20060102", dateValue)
	if err != nil {
		return "", fmt.Errorf("invalid date format")
	}

	today := time.Now().Truncate(24 * time.Hour)
	parsedToday, _ := time.Parse("20060102", today.Format("20060102"))

	if repeat == "" {
		if parsedDate.Before(parsedToday) {
			return parsedToday.Format("20060102"), nil
		}
		return dateValue, nil
	}

	if parsedDate.Before(parsedToday) {
		reference := parsedToday.AddDate(0, 0, -1)
		next, err := nextdate.NextDate(reference, parsedDate.Format("20060102"), repeat)
		if err != nil {
			return "", err
		}
		return next, nil
	}

	return dateValue, nil
}
