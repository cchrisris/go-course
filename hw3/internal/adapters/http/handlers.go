package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/cchrisris/go-course/hw3/internal/ports"
	"github.com/cchrisris/go-course/hw3/internal/usecase"
)

func NewHandler(service usecase.Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /deposit", func(w http.ResponseWriter, r *http.Request) {
		handleDeposit(r.Context(), w, r, service)
	})
	mux.HandleFunc("POST /transfer", func(w http.ResponseWriter, r *http.Request) {
		handleTransfer(r.Context(), w, r, service)
	})
	mux.HandleFunc("GET /balance", func(w http.ResponseWriter, r *http.Request) {
		handleGetBalance(r.Context(), w, r, service)
	})
	return withStandardMiddlewares(mux)
}

type depositRequest struct {
	AccountID string `json:"accountId"`
	Amount    int64  `json:"amount"`
}

type transferRequest struct {
	FromID string `json:"fromId"`
	ToID   string `json:"toId"`
	Amount int64  `json:"amount"`
}

type balanceResponse struct {
	AccountID string `json:"accountId"`
	Balance   int64  `json:"balance"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func handleDeposit(ctx context.Context, w http.ResponseWriter, r *http.Request, service usecase.Service) {
	var req depositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	newBalance, err := service.Deposit(ctx, req.AccountID, req.Amount)
	if err != nil {
		writeByError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, balanceResponse{AccountID: req.AccountID, Balance: newBalance})
}

func handleTransfer(ctx context.Context, w http.ResponseWriter, r *http.Request, service usecase.Service) {
	var req transferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := service.Transfer(ctx, req.FromID, req.ToID, req.Amount); err != nil {
		writeByError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleGetBalance(ctx context.Context, w http.ResponseWriter, r *http.Request, service usecase.Service) {
	accountID := r.URL.Query().Get("id")
	if accountID == "" {
		writeError(w, http.StatusBadRequest, "query param 'id' is required")
		return
	}
	balance, err := service.GetBalance(ctx, accountID)
	if err != nil {
		writeByError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, balanceResponse{AccountID: accountID, Balance: balance})
}

func withStandardMiddlewares(next http.Handler) http.Handler {
	return http.TimeoutHandler(withJSONContentType(next), 15*time.Second, `{"error":"timeout"}`)
}

func withJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func writeByError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ports.ErrNotFound):
		writeError(w, http.StatusNotFound, "account not found")
	case errors.Is(err, usecase.ErrInvalidAmount):
		writeError(w, http.StatusBadRequest, "amount must be positive")
	case errors.Is(err, usecase.ErrSameAccount):
		writeError(w, http.StatusConflict, "cannot transfer to the same account")
	case errors.Is(err, usecase.ErrInsufficientFunds):
		writeError(w, http.StatusConflict, "insufficient funds")
	default:
		writeError(w, http.StatusInternalServerError, "internal error: "+strconv.Quote(err.Error()))
	}
}


