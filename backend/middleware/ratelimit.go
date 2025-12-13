package middleware

import (
	"backend/config"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware aplica rate limiting a las rutas
// Permite 10 requests por 20 segundos por defecto
func RateLimitMiddleware() gin.HandlerFunc {
	return RateLimitWithConfig(10, 20*time.Second)
}

// RateLimitWithConfig permite personalizar el límite y la ventana de tiempo
func RateLimitWithConfig(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Usar la IP del cliente como identificador
		identifier := c.ClientIP()

		// Crear contexto con timeout
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Verificar el rate limit
		allowed, remaining, resetTime, err := config.CheckRateLimit(ctx, identifier, limit, window)
		
		// Agregar headers informativos sobre el rate limit
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

		if err != nil {
			// Si hay error con Redis, loggear pero permitir la petición
			println("Warning: Rate limit check failed:", err.Error())
			c.Next()
			return
		}

		// Si se excedió el límite, retornar error 429 (Too Many Requests)
		if !allowed {
			retryAfter := int(time.Until(resetTime).Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}

			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":      "Too many requests",
				"message":    "Rate limit exceeded. Please try again later.",
				"retryAfter": retryAfter,
			})
			c.Abort()
			return
		}

		// Continuar con la siguiente función
		c.Next()
	}
}