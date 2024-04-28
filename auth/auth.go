package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func BasicAuth(username, password string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization Header")
			}

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || parts[0] != "Basic" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization Header")
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Error Decoding Authorization Header")
			}

			credentials := strings.SplitN(string(decoded), ":", 2)
			if len(credentials) != 2 || credentials[0] != username || credentials[1] != password {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Credentials")
			}

			return next(c)
		}
	}
}
