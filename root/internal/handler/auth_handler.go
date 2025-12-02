package handler

import (
    "encoding/json"
    "net/http"

    "github.com/FUADIKAMIL/taskify/internal/service"
)

func jsonResp(w http.ResponseWriter, code int, v any) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    return json.NewEncoder(w).Encode(v)
}

type registerReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type loginReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func RegisterHandler(auth *service.AuthService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        decode := func() (registerReq, error) {
            var req registerReq
            return req, json.NewDecoder(r.Body).Decode(&req)
        }

        req, err := decode()
        if err != nil {
            _ = jsonResp(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
            return
        }

        u, err := auth.Register(r.Context(), req.Username, req.Password)
        if err != nil {
            _ = jsonResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
            return
        }

        _ = jsonResp(w, http.StatusCreated, u)
    }
}

func LoginHandler(auth *service.AuthService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        decode := func() (loginReq, error) {
            var req loginReq
            return req, json.NewDecoder(r.Body).Decode(&req)
        }

        req, err := decode()
        if err != nil {
            _ = jsonResp(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
            return
        }

        tok, err := auth.Login(r.Context(), req.Username, req.Password)
        if err != nil {
            _ = jsonResp(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
            return
        }

        _ = jsonResp(w, http.StatusOK, map[string]string{"token": tok})
    }
}