package main

import (
	"log"
	"net/http"
	"goPractice/db"
	"goPractice/routes"
)

func main() {
	log.Println("서버 시작!")
	db.InitDB()
	defer db.CloseDB()

	router := routes.NewRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal(err)
    }
}
