package handler

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/FUADIKAMIL/taskify/internal/model"
    "github.com/FUADIKAMIL/taskify/internal/service"
    "github.com/FUADIKAMIL/taskify/pkg/middleware"
)

func getUserIDFromCtx(r *http.Request) int64 {
    v := r.Context().Value(middleware.UserIDKey)
    if v == nil {
        return 0
    }
    return v.(int64)
}

func CreateTaskHandler(svc *service.TaskService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var t model.Task
        if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
            jsonResp(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
            return
        }

        t.UserID = getUserIDFromCtx(r)

        created, err := svc.CreateTask(r.Context(), t)
        if err != nil {
            jsonResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
            return
        }

        jsonResp(w, http.StatusCreated, created)
    }
}

func ListTaskHandler(svc *service.TaskService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        uid := getUserIDFromCtx(r)

        list, err := svc.GetAllTasks(r.Context(), uid)
        if err != nil {
            jsonResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
            return
        }

        jsonResp(w, http.StatusOK, list)
    }
}

func GetOneTaskHandler(svc *service.TaskService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, _ := strconv.ParseInt(idStr, 10, 64)
        uid := getUserIDFromCtx(r)

        t, err := svc.GetTask(r.Context(), id, uid)
        if err != nil {
            jsonResp(w, http.StatusNotFound, map[string]string{"error": "not found"})
            return
        }

        jsonResp(w, http.StatusOK, t)
    }
}

func UpdateTaskHandler(svc *service.TaskService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, _ := strconv.ParseInt(idStr, 10, 64)

        var t model.Task
        if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
            jsonResp(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
            return
        }

        t.ID = id
        t.UserID = getUserIDFromCtx(r)

        updated, err := svc.UpdateTask(r.Context(), t)
        if err != nil {
            jsonResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
            return
        }

        jsonResp(w, http.StatusOK, updated)
    }
}

func DeleteTaskHandler(svc *service.TaskService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, _ := strconv.ParseInt(idStr, 10, 64)
        uid := getUserIDFromCtx(r)

        if err := svc.DeleteTask(r.Context(), id, uid); err != nil {
            jsonResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func ToggleTaskHandler(svc *service.TaskService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, _ := strconv.ParseInt(idStr, 10, 64)
        uid := getUserIDFromCtx(r)

        t, err := svc.GetTask(r.Context(), id, uid)
        if err != nil {
            jsonResp(w, http.StatusNotFound, map[string]string{"error": "not found"})
            return
        }

        t.Completed = !t.Completed

        updated, err := svc.UpdateTask(r.Context(), t)
        if err != nil {
            jsonResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
            return
        }

        jsonResp(w, http.StatusOK, updated)
    }
}