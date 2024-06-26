package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var Validate = validator.New()

func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, map[string]string{"error": err.Error()})
}

func ValidatePayload(w http.ResponseWriter, payload any) bool {
	if err := Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return false
	}

	return true
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

type TransactionFunc func(tx *sql.Tx) error

func WithTransaction(ctx context.Context, db *sql.DB, w http.ResponseWriter, fn TransactionFunc) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			WriteError(w, http.StatusInternalServerError, fmt.Errorf("transaction failed: %v", p))
			panic(p)
		} else if err != nil {
			tx.Rollback()
			WriteError(w, http.StatusInternalServerError, err)
		} else {
			err = tx.Commit()
			if err != nil {
				WriteError(w, http.StatusInternalServerError, err)
			}
		}
	}()

	err = fn(tx)
}

func Exec(execable interface{}, query string, args ...interface{}) (sql.Result, error) {
	switch e := execable.(type) {
	case *sql.Tx:
		return e.Exec(query, args...)
	case *sql.DB:
		return e.Exec(query, args...)
	default:
		return nil, fmt.Errorf("unsupported execable type")
	}
}
