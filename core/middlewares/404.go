package middlewares

import (
	"log/slog"
	"net/http"
)

func NotFound() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("{\"status\": 404, \"message\": \"Not Found\"}"))
		if err != nil {
			slog.Error("Unable to write not found response", err)
		}

		attributes := []slog.Attr{
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.Int("status_code", http.StatusNotFound),
		}

		slog.LogAttrs(r.Context(), slog.LevelInfo, "Not Found", attributes...)
	}
}
