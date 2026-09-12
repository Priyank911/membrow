package observability

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
)

const RequestIDKey = "request_id"

func RequestContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-Id")
		if requestID == "" {
			requestID = newRequestID()
		}
		c.Locals(RequestIDKey, requestID)
		c.Set("X-Request-Id", requestID)
		start := time.Now()
		err := c.Next()
		c.Set("X-Response-Time", time.Since(start).String())
		return err
	}
}

func RequestIDFromCtx(c *fiber.Ctx) string {
	requestID, _ := c.Locals(RequestIDKey).(string)
	return requestID
}

func newRequestID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return time.Now().UTC().Format("20060102150405.000")
	}
	return hex.EncodeToString(buf)
}
