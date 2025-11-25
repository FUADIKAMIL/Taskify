package handler

import (
    "encoding/json"
    "net/http"

    "github.com/yourname/taskify/internal/service"
)

type AuthHandler struct{ Svc *service.AuthService }

func NewAuthHandler(s *service.AuthService) *AuthHandler { return &AuthHandler{Svc: s} }

func jsonResp(w http.ResponseWriter, code int, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    _ = json.NewEncoder(w).Encode(v)
}

type registerReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req registerReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        jsonResp(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
        return
    }
    u, err := h.Svc.Register(r.Context(), req.Username, req.Password)
    if err != nil {
        jsonResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }
    jsonResp(w, http.StatusCreated, u)
}

type loginReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req loginReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        jsonResp(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
        return
    }
    tok, err := h.Svc.Login(r.Context(), req.Username, req.Password)
    if err != nil {
        jsonResp(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
        return
    }
    jsonResp(w, http.StatusOK, map[string]string{"token": tok})
}
