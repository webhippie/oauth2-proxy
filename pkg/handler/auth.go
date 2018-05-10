package handler

import (
	"net/http"
	"path"

	"github.com/webhippie/oauth2-proxy/pkg/config"
)

// Auth handles the callback from the OAuth2 provider.
func Auth(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(
			w,
			r,
			path.Join(
				cfg.Server.Root,
				"login",
			),
			http.StatusMovedPermanently,
		)
	}
}
