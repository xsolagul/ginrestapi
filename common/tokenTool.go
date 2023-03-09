package common

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"simplerestapi/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type MyCustomClaims struct {
	jwt.StandardClaims
	UserId int `json:"userId"`
}
type tokenTool struct {
	Db *gorm.DB
}

func NewTokenTool(db *gorm.DB) *tokenTool {
	return &tokenTool{Db: db}
}

func GenerateToken(user models.User) (ss *string, err error) {
	// add role in claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyCustomClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "test",
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.ID,
	})

	result, error := token.SignedString([]byte("MySignature"))
	err = error
	ss = &result
	if err != nil {
		return nil, err
	}
	return ss, nil

}
func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		} //check sign method

		return []byte("MySignature"), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

func (t *tokenTool) AuthorizationMiddlewareNorRole(c *gin.Context) {
	// add role of user in claim
	s := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(s, "Bearer ")

	if Claims, err := validateToken(token); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	} else {
		var user models.User
		userId := Claims["userId"].(float64)
		models.GetUser(t.Db, &user, int(userId))
		fmt.Print(reflect.TypeOf(Claims["userId"]))
	}
}
func (t *tokenTool) AuthorizationMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.Request.Header.Get("Authorization")
		token := strings.TrimPrefix(s, "Bearer ")

		if Claims, err := validateToken(token); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error": err.Error()})
			return
		} else {
			var user models.User
			userId := Claims["userId"].(float64)
			err := models.GetUser(t.Db, &user, int(userId))
			if err!= nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token id not found"})
					return
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}else {
				userRoles := user.Role
				for _, role := range allowedRoles {
					if role == userRoles {
						c.Set("userId",user.ID)
						return
					}
				}
				c.AbortWithStatus(http.StatusForbidden)
			}
		}
	}
}
