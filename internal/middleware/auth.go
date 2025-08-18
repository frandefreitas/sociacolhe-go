package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth struct { Secret string }

type Claims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j JWTAuth) Issue(userID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(j.Secret))
}

func (j JWTAuth) Require(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"missing bearer token"})
			return
		}
		token := strings.TrimPrefix(h, "Bearer ")
		claims := &Claims{}
		parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
			return []byte(j.Secret), nil
		})
		if err != nil || !parsed.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"invalid token"}); return
		}
		if len(roles) > 0 {
			ok := false
			for _, r := range roles { if r == claims.Role { ok = true; break } }
			if !ok { c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error":"forbidden"}); return }
		}
		c.Set("uid", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
