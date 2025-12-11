package main

// import (
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// )

// func main() {
// 	secret := "your-secret-key-change-in-production"

// 	userID := int64(12345678999)
// 	if len(os.Args) > 1 {
// 		id, err := strconv.ParseInt(os.Args[1], 10, 64)
// 		if err == nil {
// 			userID = id
// 		}
// 	}

// 	claims := jwt.MapClaims{
// 		"tg_user_id": userID,
// 		"exp":        time.Now().Add(24 * time.Hour).Unix(),
// 		"iat":        time.Now().Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString([]byte(secret))
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
// 		os.Exit(1)
// 	}

// 	fmt.Println(tokenString)
// }
