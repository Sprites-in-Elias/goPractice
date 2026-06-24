package main

import (
    "log"
	"os"
    // "github.com/joho/godotenv"
)


func main() {
	log.Println("서버 시작!")
    
    // 1. .env 로드 시도
    // err := godotenv.Load()
    // if err != nil {
    //     log.Println("경고: .env 파일을 찾을 수 없습니다 (서버 배포 환경이면 정상입니다)")
    // }

    // 2. 환경 변수 읽기
    mongoURI := os.Getenv("MONGODB_URI")
    
    // 3. 확인용 로그 추가
    if mongoURI == "" {
        log.Fatal("심각: MONGODB_URI 환경 변수를 읽어오지 못했습니다! 설정이 되었는지 확인하세요.")
    } else {
        log.Println("성공: 가져온 URI -> ", mongoURI)
    }
}