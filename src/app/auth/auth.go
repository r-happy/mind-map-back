package auth

import (
	"net/http"
	"strings"

	"github.com/form3tech-oss/jwt-go"
	"github.com/labstack/echo/v4"
)

// decode jwt token ========================
func decodeToken(tokenString string) (string, error) {
	resolveTokenString := strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(resolveTokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if user_id, ok := claims["user_id"].(string); ok {
			return user_id, nil
		}
		return "", err
	}
	return "", err
}

// auth handler ========================
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.String(http.StatusUnauthorized, "token is empty")
		}
		user_id, err := decodeToken(token)
		if err != nil {
			return c.String(http.StatusUnauthorized, err.Error())
		}
		c.Set("user_id", user_id)
		return next(c)
	}
}
