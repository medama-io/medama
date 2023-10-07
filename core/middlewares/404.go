package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/go-faster/jx"
)

func NotFound() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")

		e := jx.GetEncoder()
		e.ObjStart()
		e.FieldStart("error")
		e.StrEscape("not found")
		e.ObjEnd()

		_, _ = w.Write(e.Bytes())

		attributes := []slog.Attr{
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.Int("status_code", http.StatusNotFound),
		}

		slog.LogAttrs(r.Context(), slog.LevelInfo, "Not Found", attributes...)
	}
}
