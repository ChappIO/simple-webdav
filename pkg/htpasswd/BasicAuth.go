package htpasswd

import (
	"context"
	"net/http"
)

func (file *File) BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// first we check auth
		writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		username, password, ok := request.BasicAuth()
		if !ok {
			http.Error(writer, "Not Authorized", 401)
			return
		}
		username, ok = file.Authenticate(username, []byte(password))
		if !ok {
			http.Error(writer, "Not Authorized", 401)
			return
		}
		next.ServeHTTP(writer, request.WithContext(
			context.WithValue(
				request.Context(),
				"Username",
				username,
			),
		))
	})
}
