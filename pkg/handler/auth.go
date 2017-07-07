package handler

import (
	"net/http"

	"github.com/unrolled/render"
	"github.com/webhippie/oauth2-proxy/pkg/config"
)

// Auth handles the callback from the OAuth2 provider.
func Auth(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.HTML(
			w,
			http.StatusOK,
			"login",
			map[string]string{
				"Title": config.Server.Title,
			},
		)
	}
}
