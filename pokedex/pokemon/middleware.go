package pokemon

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// EnrichedLogger extracts contextual information (trainer header / target id)
// and logs them. It also stores those values into Gin context for downstream handlers.
func EnrichedLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		trainer := c.GetHeader("X-Trainer")
		if trainer == "" {
			trainer = c.Query("trainer")
		}

		c.Set("trainer", trainer)

		// if the route contains a :id parameter we try to attach the targeted Pokemon
		idStr := c.Param("id")
		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err == nil {
				if p, _ := GetByID(id); p != nil {
					c.Set("target_pokemon", p)
				}
			}
		}

		// Log — keep short, enriched with trainer and target when available
		if trainer != "" {
			log.Printf("[ENRICHED-LOG] method=%s path=%s trainer=%s id=%s", c.Request.Method, c.FullPath(), trainer, idStr)
		} else {
			log.Printf("[ENRICHED-LOG] method=%s path=%s id=%s", c.Request.Method, c.FullPath(), idStr)
		}

		c.Next()
	}
}

// SimpleAuth checks for a simple header or cookie. This is intentionally tiny for the TP.
// Expect either header X-Admin-Token or cookie admin_token with value equal to adminSecret.
func SimpleAuth(adminSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check header
		token := c.GetHeader("X-Admin-Token")

		// Optionnal : try cookie
		if token == "" {
			if cookie, err := c.Cookie("admin_token"); err == nil {
				token = cookie
			}
		}

		if token != adminSecret {
			// no auth — abort
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - admin token missing/invalid"})
			c.Abort()
			return
		}

		// attach admin info to context for handlers
		c.Set("admin_authenticated", true)
		c.Next()
	}
}

// RateLimitMiddleware is a tiny global rate limiter per-route used for demo.
// It allows maxRequests every window duration. When limit is hit, it returns 429.
func RateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	type bucket struct {
		count     int
		windowEnd time.Time
	}

	var mu sync.Mutex
	var b bucket

	return func(c *gin.Context) {
		mu.Lock()
		now := time.Now()
		if now.After(b.windowEnd) {
			// new window
			b.count = 0
			b.windowEnd = now.Add(window)
		}

		if b.count >= maxRequests {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		b.count++
		mu.Unlock()

		c.Next()
	}
}

// FatigueMiddleware demonstrates a server-wide artificial delay when a special header
// is present. This helps show impact of adding delays and why middlewares must be used wisely.
func FatigueMiddleware(delay time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Server-Fatigue") == "true" {
			time.Sleep(delay)
		}
		c.Next()
	}
}
