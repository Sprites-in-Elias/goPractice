package db 

import (
	"os"
	"context"
	"time"
	"log"
	
	"github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/jackc/pgx/v5/pgxpool"
)

// 외부에서 접근 가능하도록 대문자로 시작
var (
    MongoClient *mongo.Client
    PgPool      *pgxpool.Pool
)

func LoadMongoClient(mongoURI string) {
	// MongoDB 클라이언트 연결
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("연결 객체 생성 실패:", err)
	}
	MongoClient = client
	log.Println("MongoDB 연결 성공!")
}

func LoadPgPool(pgURI string) {
	pool, err := pgxpool.New(context.Background(), pgURI)
	if err != nil {
		log.Fatal("Postgres 연결 실패:", err)
	}
	PgPool = pool
	log.Println("Postgres 연결 성공!")
}

func CloseDB() {
	
	// 맥락 생성
	ctx := context.Background() 
	
	// MongoDB 연결 종료
	if MongoClient != nil {
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Printf("MongoDB 종료 에러: %v", err)
		}
	}

	// PostgreSQL 연결 종료
    if PgPool != nil {
        PgPool.Close()
    }

	// 종료 로그
	log.Println("DB 연결이 안전하게 종료되었습니다.")
}

// DB 연결을 담당하는 함수
func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("경고: .env 파일을 찾을 수 없습니다 (서버 배포 환경이면 정상입니다)")
	}
	// 클라이언트 연결
	LoadMongoClient(os.Getenv("MONGODB_URI"))
	LoadPgPool(os.Getenv("POSTGRES_URI"))
}