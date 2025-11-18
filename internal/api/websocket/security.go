package websocket

import (
	"backend_go/internal/infrastructure/config"
	"net/http"
	"strings"
)

// checkOrigin создает функцию проверки origin
func checkOrigin(cfg *config.Config) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		origin := r.Header.Get("Origin")

		// Разрешить запросы без Origin
		if origin == "" {
			return true
		}

		// В режиме разработки разрешаем все origin
		if cfg.Environment == "development" {
			return true
		}

		// Проверяем разрешенные origin'ы
		for _, allowedOrigin := range cfg.WebSocket.AllowedOrigins {
			if origin == allowedOrigin {
				return true
			}

			// Поддержка wildcard поддоменов
			if strings.HasPrefix(allowedOrigin, "*.") {
				domain := strings.TrimPrefix(allowedOrigin, "*.")
				if strings.HasSuffix(origin, domain) {
					return true
				}
			}
		}

		return false
	}
}
