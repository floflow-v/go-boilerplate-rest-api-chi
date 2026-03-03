package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Operation completed successfully"`
}

func JSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload == nil {
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}
