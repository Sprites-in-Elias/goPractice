package main

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"time"
	"github.com/jackc/pgx/v5"
)

func main() {
	log.Println("서버 시작!")

	// 1. .env 로드 시도
	err := godotenv.Load()
	if err != nil {
		log.Println("경고: .env 파일을 찾을 수 없습니다 (서버 배포 환경이면 정상입니다)")
	}

	// 2. 환경 변수 읽기
	m_uri := os.Getenv("MONGODB_URI")
	p_uri := os.Getenv("POSTGRES_URI")

	// 3. 확인용 로그 추가
	if m_uri == "" {
		log.Fatal("심각: MONGODB_URI 환경 변수를 읽어오지 못했습니다! 설정이 되었는지 확인하세요.")
	} else {
		log.Println("성공: 가져온 URI -> ", m_uri)

		// 3. MongoDB 연결 설정
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientOptions := options.Client().ApplyURI(m_uri)
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal("연결 객체 생성 실패:", err)
		}

		// 4. 진짜 연결됐는지 Ping 날려보기 (가장 확실함)
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal("MongoDB 서버 응답 없음 (연결 실패):", err)
		}

		log.Println("성공: MongoDB와 성공적으로 연결되었습니다!")

		// 프로그램 종료 시 연결 해제
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				log.Fatal(err)
			}
		}()

		collection := client.Database("testDatabaseName").Collection("testCollectionName")

		var result bson.M
		// 조회용 context 생성 (5초 제한)
		findCtx, findCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer findCancel()

		// FindOne: 첫 번째 문서를 가져옴
		err = collection.FindOne(findCtx, bson.D{}).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				log.Println("조회 결과: 데이터가 없습니다.")
			} else {
				log.Println("조회 에러:", err)
			}
		} else {
			// 7. 결과 출력
			log.Println("가져온 데이터:", result)
			log.Println("_id 필드값:", result["_id"])
		}

		// 6. 모든 데이터 가져오기 (필터에 nil을 넣으면 전부 다!)
		cursor, err := collection.Find(context.Background(), bson.D{})
		if err != nil {
			log.Fatal("조회 실패:", err)
		}
		defer cursor.Close(context.Background()) // 작업 끝나면 커서 닫기

		// 7. 커서를 이용해 하나씩 출력
		log.Println("--- 전체 데이터 출력 시작 ---")
		for cursor.Next(context.Background()) {
			var result bson.M
			if err := cursor.Decode(&result); err != nil {
				log.Fatal("데이터 변환 실패:", err)
			}
			log.Println("데이터:", result)
		}

		// 커서 순회 중 에러가 발생했는지 확인
		if err := cursor.Err(); err != nil {
			log.Fatal("커서 에러:", err)
		}
		log.Println("--- 전체 데이터 출력 끝 ---")
	}
	if p_uri == "" {
		log.Fatal("심각: POSTGRES_URL 환경 변수를 읽어오지 못했습니다! 설정이 되었는지 확인하세요.")
	} else {
		log.Println("성공: 가져온 URI -> ", p_uri)

		conn, err := pgx.Connect(context.Background(), p_uri)
		
		if err != nil {
			log.Fatal("연결 실패:", err)
		}
		defer conn.Close(context.Background())
		log.Println("연결 성공!!")
		// 1. 쿼리 실행
		rows, err := conn.Query(context.Background(), "SELECT id, name, age, role FROM test")
		if err != nil {
			log.Fatal("조회 실패:", err)
		}
		defer rows.Close()

		// 2. 데이터 순회하며 출력
		for rows.Next() {
			var id int
			var name string
			var age int
			var role string
			err := rows.Scan(&id, &name, &age, &role)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("사용자 정보: ID=%d, 이름=%s, 나이=%d, 역할=%s\n", id, name, age, role)
		}
	}
}
