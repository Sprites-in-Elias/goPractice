package main

import (
    "log"
	"os"
    "github.com/joho/godotenv"
)


func main() {
    log.Println("서버 준비 완료 했음!!!")
    _ = godotenv.Load()
	mongoURI := os.Getenv("MONGODB_URI")
	log.Println(mongoURI)
}