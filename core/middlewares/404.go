package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/go-faster/jx"
)

const errMessage = "api path not found"

func NotFound() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")

		e := jx.GetEncoder()
		e.ObjStart()
		e.FieldStart("error")
		e.ObjStart()
		e.FieldStart("code")
		e.Int(http.StatusNotFound)
		e.FieldStart("message")
		e.StrEscape(errMessage)
		e.ObjEnd()
		e.ObjEnd()

		_, _ = w.Write(e.Bytes())

		attributes := []slog.Attr{
			slog.String("path", req.URL.Path),
			slog.String("method", req.Method),
			slog.Int("status_code", http.StatusNotFound),
			slog.String("message", errMessage),
			slog.String("User-Agent", req.Header.Get("User-Agent")),
		}

		slog.LogAttrs(req.Context(), slog.LevelInfo, "Not Found", attributes...)
	}
}
