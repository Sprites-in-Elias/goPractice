package routes

import (
    "goPractice/test"
	"goPractice/liveKit"
    "github.com/go-chi/chi/v5"
)

func NewRouter() *chi.Mux {
    r := chi.NewRouter()

    r.Route("/test", func(r chi.Router) {
		r.Get("/hello", test.HelloHandler)	
		r.Get("/dummyNumber", test.GetDummyNumber)
		r.Get("/livekit/token", liveKit.TestRoomToken)
    })
    
	r.Get("/mOneTest", test.MongoOneTest)
	r.Get("/mListTest", test.MongoListTest)
	r.Get("/pTest", test.PgTest)

    return r
}