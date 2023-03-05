package common

import (
	"errors"
	"fmt"
	"net/http"
	"simplerestapi/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)


type MyCustomClaims struct {
    jwt.StandardClaims
    UserId   int `json:"userId"`
    Username string `json:"username"`
}

func GenerateToken(user models.User)(ss *string,err error){
	// add role in claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyCustomClaims{
        StandardClaims: jwt.StandardClaims{
            Issuer:    "test",
            ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
        UserId:   user.ID,
        Username: user.Name,
    })

	result, error := token.SignedString([]byte("MySignature"))
	err = error
	ss = &result
	if err != nil {
		return nil,err
	}
	return ss,nil

}
func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}//check sign method 

		return []byte("MySignature"), nil
	}) 
	if err != nil{
		return nil,err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

func authorizationMiddleware(c *gin.Context) {
	// add role of user in claim
	s := c.Request.Header.Get("Authorization")

	token := strings.TrimPrefix(s, "Bearer ")

	if Claims,err := validateToken(token); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}else {
		fmt.Print(Claims)
	}
}