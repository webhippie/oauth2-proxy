package handler

import (
	"net/http"

	"github.com/codehack/fail"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/webhippie/oauth2-proxy/pkg/config"
	"github.com/webhippie/oauth2-proxy/pkg/templates"
)

// Auth handles the callback from the OAuth2 provider.
func Auth(logger log.Logger) http.HandlerFunc {
	logger = log.WithPrefix(logger, "handler", "auth")

	return func(w http.ResponseWriter, req *http.Request) {
		err := templates.Load(logger).ExecuteTemplate(
			w,
			"login.tmpl",
			map[string]string{
				"Title": config.Server.Title,
				"Root":  config.Server.Root,
			},
		)

		if err != nil {
			level.Warn(logger).Log(
				"msg", "failed to process index template",
				"err", err,
			)

			fail.Error(w, fail.Cause(err).Unexpected())
			return
		}
	}
}
