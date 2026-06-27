package network

import (
	"net/http"
	"encoding/json"
)

// 2. 전체 응답 구조체 정의 (중첩 구조체 사용)
type APIResponse struct {
    Status  string `json:"status"`
    Code    int    `json:"code"`
    Data    any `json:"data"` // 여러 명의 사용자를 담을 배열(슬라이스)
    Message string `json:"message"`
}

func ResponseSuccess(w http.ResponseWriter, status string, code int, data any, message string) {
	resp := APIResponse{
		Status:  status,
		Code:    code,
		Data:    data,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}