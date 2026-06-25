package test

import (
	"context"
	"log"
	"time"
	"fmt"
	"net/http"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"goPractice/db"
)

// 1. 하위 데이터 구조체 정의
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Age   int32  `json:"age"`
	Role  string `json:"role"`
}

// 2. 전체 응답 구조체 정의 (중첩 구조체 사용)
type APIResponse struct {
    Status  string `json:"status"`
    Code    int    `json:"code"`
    Data    any `json:"data"` // 여러 명의 사용자를 담을 배열(슬라이스)
    Message string `json:"message"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
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

func MongoOneTest(w http.ResponseWriter, r *http.Request) {
	collection := db.MongoClient.Database("testDatabaseName").Collection("testCollectionName")
	var result bson.M

	findCtx, findCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer findCancel()

	err := collection.FindOne(findCtx, bson.D{}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("조회 결과: 데이터가 없습니다.")
			ResponseSuccess(w, "error", 404, nil, "데이터가 없습니다.")
		} else {
			log.Println("조회 에러:", err)
			ResponseSuccess(w, "error", 500, nil, fmt.Sprintf("조회 에러: %v", err))
		}
	} else {
		users := []User{
			{ID: 1, Name: result["name"].(string), Age: result["age"].(int32), Role: result["role"].(string)},
		}
		ResponseSuccess(w, "success", 200, users, "사용자 정보를 성공적으로 불러왔습니다.")
	}
}

func MongoListTest(w http.ResponseWriter, r *http.Request) {
	collection := db.MongoClient.Database("testDatabaseName").Collection("testCollectionName")
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		ResponseSuccess(w, "error", 500, nil, fmt.Sprintf("조회 실패: %v", err))
		return
	}
	defer cursor.Close(context.Background())

	var userList []User

	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			ResponseSuccess(w, "error", 500, nil, err.Error())
			return
		}
		userList = append(userList, User{
			ID:   1,
			Name: result["name"].(string),
			Age:  result["age"].(int32),
			Role: result["role"].(string),
		})
	}

	// 커서 순회 중 에러가 발생했는지 확인
	if err := cursor.Err(); err != nil {
		log.Println("커서 에러:", err)
		ResponseSuccess(w, "error", 500, nil, fmt.Sprintf("커서 에러: %v", err))
	} else {
		ResponseSuccess(w, "success", 200, userList, "사용자 목록을 성공적으로 불러왔습니다.")
	}
}

func PgTest(w http.ResponseWriter, r *http.Request) {
	rows, err := db.PgPool.Query(context.Background(), "SELECT id, name, age, role FROM test")
	if err != nil {
		log.Fatal("조회 실패:", err)
	}
	defer rows.Close()

	var userList []User
	// 2. 데이터 순회하며 출력
	for rows.Next() {
		var id int
		var name string
		var age int32
		var role string
		err := rows.Scan(&id, &name, &age, &role)
		if err != nil {
			log.Fatal(err)
		}
		userList = append(userList, User{
			ID:   id,
			Name: name,
			Age:  age,
			Role: role,
		})
	}
	ResponseSuccess(w, "success", 200, userList, "사용자 목록을 성공적으로 불러왔습니다.")
}