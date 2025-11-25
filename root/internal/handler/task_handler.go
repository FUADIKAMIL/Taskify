package handler

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/yourname/taskify/internal/model"
    "github.com/yourname/taskify/internal/service"
    "github.com/yourname/taskify/pkg/middleware"
)

type TaskHandler struct{ Svc *service.TaskService }

func NewTaskHandler(s *service.TaskService) *TaskHandler { return &TaskHandler{Svc: s} }

func jsonResponse(w http.ResponseWriter, code int, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    _ = json.NewEncoder(w).Encode(v)
}

func getUserIDFromCtx(r *http.Request) int64 {
    v := r.Context().Value(middleware.UserIDKey)
    if v == nil {
        return 0
    }
    return v.(int64)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
    var t model.Task
    if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
        jsonResp(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
        return
    }
    t.UserID = getUserIDFromCtx(r)
    created, err := h.Svc.CreateTask(r.Context(), t)
    if err != nil {
        jsonResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }
    jsonResp(w, http.StatusCreated, created)
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
    uid := getUserIDFromCtx(r)
    list, err := h.Svc.GetAllTasks(r.Context(), uid)
    if err != nil {
        jsonResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }
    jsonResp(w, http.StatusOK, list)
}

func (h *TaskHandler) GetOne(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, _ := strconv.ParseInt(idStr, 10, 64)
    uid := getUserIDFromCtx(r)
    t, err := h.Svc.GetTask(r.Context(), id, uid)
    if err != nil {
        jsonResp(w, http.StatusNotFound, map[string]string{"error": "not found"})
        return
    }
    jsonResp(w, http.StatusOK, t)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, _ := strconv.ParseInt(idStr, 10, 64)
    var t model.Task
    if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
        jsonResp(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
        return
    }
    t.ID = id
    t.UserID = getUserIDFromCtx(r)
    updated, err := h.Svc.UpdateTask(r.Context(), t)
    if err != nil {
        jsonResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }
    jsonResp(w, http.StatusOK, updated)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, _ := strconv.ParseInt(idStr, 10, 64)
    uid := getUserIDFromCtx(r)
    if err := h.Svc.DeleteTask(r.Context(), id, uid); err != nil {
        jsonResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) Toggle(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, _ := strconv.ParseInt(idStr, 10, 64)
    uid := getUserIDFromCtx(r)
    t, err := h.Svc.GetTask(r.Context(), id, uid)
    if err != nil {
        jsonResp(w, http.StatusNotFound, map[string]string{"error": "not found"})
        return
    }
    t.Completed = !t.Completed
    updated, err := h.Svc.UpdateTask(r.Context(), t)
    if err != nil {
        jsonResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }
    jsonResp(w, http.StatusOK, updated)
}
