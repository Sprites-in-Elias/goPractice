package liveKit

import (
	"os"
	"log"
	"fmt"
	
	"github.com/joho/godotenv"
	"github.com/livekit/protocol/auth"
    "time"
	"net/http"
	"goPractice/network"
)

func generateToken(roomName string, identity string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("경고: .env 파일을 찾을 수 없습니다 (서버 배포 환경이면 정상입니다)")
	}
	apiKey := os.Getenv("LIVEKIT_API_KEY")
    apiSecret := os.Getenv("LIVEKIT_API_SECRET")

    at := auth.NewAccessToken(apiKey, apiSecret)
    at.SetIdentity(identity)
    at.SetValidFor(time.Hour * 1) // 토큰 유효시간 1시간
    
    grant := &auth.VideoGrant{
        RoomJoin: true,
        Room:     roomName,
    }
    at.SetVideoGrant(grant)

    return at.ToJWT()
}

func TestRoomToken(w http.ResponseWriter, r *http.Request) {
	testToken, err := generateToken("testRoom", "testUser")
	if err != nil {
		network.ResponseSuccess(w, "error", 500, nil, fmt.Sprintf("토큰 생성 실패: %v", err))
		return
	}
	network.ResponseSuccess(w, "success", 200, map[string]string{"token": testToken}, "테스트 방 토큰이 성공적으로 생성되었습니다.")
}