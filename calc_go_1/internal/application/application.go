package application

import (
	"encoding/json"
	"net/http"
	"strings"

	calculate "github.com/pashapdev/calc_go/pkg/calculation"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := calculate.Calc(req.Expression)
	if err != nil {
		var resp ErrorResponse
		if strings.Contains(err.Error(), "недопустимый символ") || strings.Contains(err.Error(), "несоответствие скобок") {
			resp.Error = "Expression is not valid"
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			resp.Error = "Internal server error"
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := Response{Result: result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
