package main

import (
	stdctx "context"
	"fmt"

	"github.com/raminderis/lenslocked/context"
	"github.com/raminderis/lenslocked/models"
)

type ctxKey string

const (
	favColor ctxKey = "color"
)

func main() {
	ctx := stdctx.Background()
	ctx = stdctx.WithValue(ctx, favColor, "blue")
	fmt.Println(favColor)

	ctxValue := ctx.Value(favColor)
	ctxValueString, ok := ctxValue.(string)
	if !ok {
		fmt.Println("it isnt int")
		return
	}
	fmt.Printf("type=%T value=%v\n", favColor, favColor)
	fmt.Printf("type=%T value=%v\n", ctxValueString, ctxValueString)

	user := models.User{
		Email: "raminder@live.com",
	}
	// ctx2 := stdctx.Background()
	ctx = context.WithUser(ctx, &user)
	retrievedUser := context.User(ctx)
	fmt.Println(retrievedUser.Email)

	fmt.Println(ctx.Value((context.UserKey)).(*models.User).Email)
}
