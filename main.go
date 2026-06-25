package main

import (
	"log"
	"net/http"
	"goPractice/db"
	"goPractice/test"
	// "fmt"
)


func httpHandler() {
	http.HandleFunc("/hello", test.HelloHandler)
	http.HandleFunc("/mOneTest", test.MongoOneTest)
	http.HandleFunc("/mListTest", test.MongoListTest)
	http.HandleFunc("/pTest", test.PgTest)
	http.ListenAndServe(":8080", nil)
}

func main() {
	log.Println("서버 시작!")
	db.InitDB()
	defer db.CloseDB()
	httpHandler()	
}
