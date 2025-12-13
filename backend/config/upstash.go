package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// UpstashRESTClient cliente para Upstash REST API
type UpstashRESTClient struct {
	URL   string
	Token string
}

var UpstashClient *UpstashRESTClient

// SetupRedis configura el cliente REST de Upstash
func SetupRedis() {
	redisURL := os.Getenv("UPSTASH_REDIS_REST_URL")
	redisToken := os.Getenv("UPSTASH_REDIS_REST_TOKEN")

	if redisURL == "" || redisToken == "" {
		log.Fatal("UPSTASH_REDIS_REST_URL y UPSTASH_REDIS_REST_TOKEN deben estar definidas en .env")
	}

	UpstashClient = &UpstashRESTClient{
		URL:   redisURL,
		Token: redisToken,
	}

	log.Println("Upstash REST client configurado correctamente")
}

// executeCommand ejecuta un comando de Redis usando la REST API
func (c *UpstashRESTClient) executeCommand(ctx context.Context, command []interface{}) (interface{}, error) {
	jsonData, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result["error"] != nil {
		return nil, fmt.Errorf("redis error: %v", result["error"])
	}

	return result["result"], nil
}

// Incr incrementa un contador
func (c *UpstashRESTClient) Incr(ctx context.Context, key string) (int64, error) {
	result, err := c.executeCommand(ctx, []interface{}{"INCR", key})
	if err != nil {
		return 0, err
	}

	if val, ok := result.(float64); ok {
		return int64(val), nil
	}
	return 0, fmt.Errorf("unexpected result type")
}

// Expire establece un tiempo de expiración para una clave
func (c *UpstashRESTClient) Expire(ctx context.Context, key string, seconds int) error {
	_, err := c.executeCommand(ctx, []interface{}{"EXPIRE", key, seconds})
	return err
}

// TTL obtiene el tiempo restante de vida de una clave
func (c *UpstashRESTClient) TTL(ctx context.Context, key string) (int64, error) {
	result, err := c.executeCommand(ctx, []interface{}{"TTL", key})
	if err != nil {
		return 0, err
	}

	if val, ok := result.(float64); ok {
		return int64(val), nil
	}
	return 0, fmt.Errorf("unexpected result type")
}

// CheckRateLimit verifica si un identificador ha excedido el límite de requests
func CheckRateLimit(ctx context.Context, identifier string, limit int, window time.Duration) (bool, int, time.Time, error) {
	key := fmt.Sprintf("ratelimit:%s", identifier)

	// Incrementar contador
	count, err := UpstashClient.Incr(ctx, key)
	if err != nil {
		// Si Redis falla, permitir la request (fail open)
		return true, limit, time.Now().Add(window), err
	}

	// Establecer expiración si es el primer request
	if count == 1 {
		err = UpstashClient.Expire(ctx, key, int(window.Seconds()))
		if err != nil {
			return true, limit, time.Now().Add(window), err
		}
	}

	// Calcular remaining
	remaining := limit - int(count)
	if remaining < 0 {
		remaining = 0
	}

	// Obtener TTL para calcular reset time
	ttl, err := UpstashClient.TTL(ctx, key)
	if err != nil {
		ttl = int64(window.Seconds())
	}
	resetTime := time.Now().Add(time.Duration(ttl) * time.Second)

	// Permitir si no ha excedido el límite
	allowed := int(count) <= limit

	return allowed, remaining, resetTime, nil
}