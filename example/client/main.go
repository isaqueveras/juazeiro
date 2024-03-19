package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/isaqueveras/juazeiro"
)

var conn *juazeiro.ClientConn

func init() {
	var err error
	if conn, err = juazeiro.NewClient("http://localhost:8181"); err != nil {
		panic(err)
	}
}

// pointer returns a pointer reference
func pointer[T any](value T) *T {
	return &value
}

func main() {
	in := &User{
		UserId: pointer(int64(1)),
		Name:   pointer("Isaque Veras"),
		Parameters: &Parameters{
			Limit:     pointer(int64(10)),
			Offset:    pointer(int64(10)),
			Total:     pointer(false),
			UserName:  pointer("Oliveira"),
			CreatedAt: pointer(time.Now()),
		},
	}

	repo := NewMainClient(conn)
	data, err := repo.GetUser(context.Background(), in)
	if err != nil {
		log.Println(err)
	}

	if data != nil {
		fmt.Println(*data.UserId, *data.Name)
	}
}
