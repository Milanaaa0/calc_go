package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	path "github.com/pashapdev/calc_go/internal/application"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/calculate", path.CalculateHandler).Methods(http.MethodPost)
	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", router)
}
