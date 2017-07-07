package handler

import (
	"net/http"

	"github.com/webhippie/oauth2-proxy/pkg/assets"
)

// Static handles all requests to static assets.
func Static() http.Handler {
	return http.StripPrefix(
		"/oauth2-proxy/assets",
		http.FileServer(
			assets.Load(),
		),
	)
}
