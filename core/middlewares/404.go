package middlewares

import (
	"net/http"

	"github.com/go-faster/jx"
	"github.com/rs/zerolog"
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

		zerolog.Ctx(req.Context()).
			Info().
			Str("path", req.URL.Path).
			Str("method", req.Method).
			Int("status_code", http.StatusNotFound).
			Str("message", errMessage).
			Str("User-Agent", req.Header.Get("User-Agent")).
			Msg("Not Found")
	}
}
