package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// pointer returns a pointer reference
func pointer[T any](value T) *T {
	return &value
}

func main() {
	router := gin.Default()

	router.POST("/create_account", func(ctx *gin.Context) {
		info := new(User)
		if err := ctx.BindJSON(&info); err != nil {
			panic(err)
		}

		v, _ := json.Marshal(&info)
		log.Println(string(v))

		ctx.JSON(http.StatusOK, &User{
			UserId: pointer(int64(0)),
			Name:   pointer("Ismael"),
		})
	})

	router.Run(":8181")
}

// juaz.RegisterProfileServer(juazeiro.NewServer(), &server{})
// type server struct{}

// func (s *server) EditProfile(ctx context.Context, in string) (string, error) {
// 	log.Println("func EditProfile", in)
// 	return "", nil
// }
