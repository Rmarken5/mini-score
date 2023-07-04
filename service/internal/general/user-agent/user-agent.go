package user_agent

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mileusna/useragent"
)

type (
	isMobileUserAgent string
)

var (
	userAgentKey isMobileUserAgent = "myKey"
)

// HandleUserAgent is a middleware function that returns if a user agent is mobile or not and sets it on request context.
func HandleUserAgent(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		agent := useragent.Parse(c.Request().UserAgent())

		ctx := context.WithValue(c.Request().Context(), userAgentKey, agent.Mobile)

		r := c.Request().WithContext(ctx)

		c.SetRequest(r)

		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}

func IsMobile(ctx context.Context) bool {
	val := ctx.Value(userAgentKey)
	var isMobile, ok bool
	if isMobile, ok = val.(bool); !ok {
		return true // if not able to find user agent on context, default to mobile for mobile first dev.
	}
	return isMobile
}
