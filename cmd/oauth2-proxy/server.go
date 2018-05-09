package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/go-chi/chi"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/oklog/run"
	"github.com/rs/zerolog/log"
	"github.com/vulcand/oxy/buffer"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/roundrobin"
	"github.com/webhippie/oauth2-proxy/pkg/config"
	"github.com/webhippie/oauth2-proxy/pkg/router"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/urfave/cli.v2"
)

var (
	httpsAddr  = "0.0.0.0:443"
	httpAddr   = "0.0.0.0:80"
	healthAddr = "127.0.0.1:9000"
)

// Server provides the sub-command to start the server.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:   "server",
		Usage:  "start the integrated server",
		Flags:  serverFlags(cfg),
		Before: serverBefore(cfg),
		Action: serverAction(cfg),
	}
}

func serverFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "health-addr",
			Value:       healthAddr,
			Usage:       "address for metrics and health",
			EnvVars:     []string{"OAUTH2_PROXY_PRIVATE_ADDR"},
			Destination: &cfg.Server.Health,
		},
		&cli.StringFlag{
			Name:        "secure-addr",
			Value:       httpsAddr,
			Usage:       "https address to bind the server",
			EnvVars:     []string{"OAUTH2_PROXY_SERVER_HTTPS"},
			Destination: &cfg.Server.Secure,
		},
		&cli.StringFlag{
			Name:        "server-addr",
			Value:       httpAddr,
			Usage:       "http address to bind the server",
			EnvVars:     []string{"OAUTH2_PROXY_SERVER_ADDR"},
			Destination: &cfg.Server.Public,
		},
		&cli.StringFlag{
			Name:        "server-root",
			Value:       "/oauth2-proxy",
			Usage:       "root path of the proxy",
			EnvVars:     []string{"OAUTH2_PROXY_SERVER_ROOT"},
			Destination: &cfg.Server.Root,
		},
		&cli.StringFlag{
			Name:        "server-host",
			Value:       "http://localhost",
			Usage:       "external access to server",
			EnvVars:     []string{"OAUTH2_PROXY_SERVER_HOST"},
			Destination: &cfg.Server.Host,
		},
		&cli.StringFlag{
			Name:        "server-cert",
			Value:       "",
			Usage:       "path to ssl cert",
			EnvVars:     []string{"OAUTH2_PROXY_SERVER_CERT"},
			Destination: &cfg.Server.Cert,
		},
		&cli.StringFlag{
			Name:        "server-key",
			Value:       "",
			Usage:       "path to ssl key",
			EnvVars:     []string{"OAUTH2_PROXY_SERVER_KEY"},
			Destination: &cfg.Server.Key,
		},
		&cli.BoolFlag{
			Name:        "server-autocert",
			Value:       false,
			Usage:       "enable let's encrypt",
			EnvVars:     []string{"OAUTH2_PROXY_AUTO_CERT"},
			Destination: &cfg.Server.AutoCert,
		},
		&cli.BoolFlag{
			Name:        "strict-curves",
			Value:       false,
			Usage:       "use strict ssl curves",
			EnvVars:     []string{"OAUTH2_PROXY_STRICT_CURVES"},
			Destination: &cfg.Server.StrictCurves,
		},
		&cli.BoolFlag{
			Name:        "strict-ciphers",
			Value:       false,
			Usage:       "use strict ssl ciphers",
			EnvVars:     []string{"OAUTH2_PROXY_STRICT_CIPHERS"},
			Destination: &cfg.Server.StrictCiphers,
		},
		&cli.StringFlag{
			Name:        "templates-path",
			Value:       "",
			Usage:       "path to custom templates",
			EnvVars:     []string{"OAUTH2_PROXY_SERVER_TEMPLATES"},
			Destination: &cfg.Server.Templates,
		},
		&cli.StringFlag{
			Name:        "assets-path",
			Value:       "",
			Usage:       "path to custom assets",
			EnvVars:     []string{"OAUTH2_PROXY_SERVER_ASSETS"},
			Destination: &cfg.Server.Assets,
		},
		&cli.StringFlag{
			Name:        "storage-path",
			Value:       "storage/",
			Usage:       "folder for storing certs and misc files",
			EnvVars:     []string{"OAUTH2_PROXY_SERVER_STORAGE"},
			Destination: &cfg.Server.Storage,
		},
		&cli.StringFlag{
			Name:        "proxy-title",
			Value:       "OAuth2 Proxy",
			Usage:       "title displayed on the login",
			EnvVars:     []string{"OAUTH2_PROXY_SERVER_TITLE"},
			Destination: &cfg.Proxy.Title,
		},
		&cli.StringSliceFlag{
			Name:    "proxy-endpoint",
			Value:   cli.NewStringSlice(),
			Usage:   "endpoints to proxy requests to",
			EnvVars: []string{"OAUTH2_PROXY_SERVER_ENDPOINTS"},
		},
		&cli.StringFlag{
			Name:        "user-header",
			Value:       "X-PROXY-USER",
			Usage:       "header for username",
			EnvVars:     []string{"OAUTH2_PROXY_USER_HEADER"},
			Destination: &cfg.Proxy.UserHeader,
		},
		&cli.BoolFlag{
			Name:        "oauth2-gitlab",
			Value:       false,
			Usage:       "enable gitlab provider",
			EnvVars:     []string{"OAUTH2_PROXY_GITLAB"},
			Destination: &cfg.Gitlab.Enabled,
		},
		&cli.StringSliceFlag{
			Name:    "oauth2-gitlab-org",
			Value:   &cli.StringSlice{},
			Usage:   "allowed organizations from gitlab",
			EnvVars: []string{"OAUTH2_PROXY_GITLAB_ORGS"},
		},
		&cli.StringFlag{
			Name:        "oauth2-gitlab-client",
			Value:       "",
			Usage:       "gitlab client id",
			EnvVars:     []string{"OAUTH2_PROXY_GITLAB_CLIENT"},
			Destination: &cfg.Gitlab.Client,
		},
		&cli.StringFlag{
			Name:        "oauth2-gitlab-secret",
			Value:       "",
			Usage:       "gitlab client secret",
			EnvVars:     []string{"OAUTH2_PROXY_GITLAB_SECRET"},
			Destination: &cfg.Gitlab.Secret,
		},
		&cli.StringFlag{
			Name:        "oauth2-gitlab-url",
			Value:       "https://gitlab.com",
			Usage:       "gitlab server url",
			EnvVars:     []string{"OAUTH2_PROXY_GITLAB_URL"},
			Destination: &cfg.Gitlab.URL,
		},
		&cli.BoolFlag{
			Name:        "oauth2-gitlab-skipverify",
			Value:       false,
			Usage:       "skip ssl verify for gitlab",
			EnvVars:     []string{"OAUTH2_PROXY_GITLAB_SKIPVERIFY"},
			Destination: &cfg.Gitlab.SkipVerify,
		},
		&cli.BoolFlag{
			Name:        "oauth2-github",
			Value:       false,
			Usage:       "enable github provider",
			EnvVars:     []string{"OAUTH2_PROXY_GITHUB"},
			Destination: &cfg.GitHub.Enabled,
		},
		&cli.StringSliceFlag{
			Name:    "oauth2-github-org",
			Value:   &cli.StringSlice{},
			Usage:   "allowed organizations from github",
			EnvVars: []string{"OAUTH2_PROXY_GITHUB_ORGS"},
		},
		&cli.StringFlag{
			Name:        "oauth2-github-client",
			Value:       "",
			Usage:       "github client id",
			EnvVars:     []string{"OAUTH2_PROXY_GITHUB_CLIENT"},
			Destination: &cfg.GitHub.Client,
		},
		&cli.StringFlag{
			Name:        "oauth2-github-secret",
			Value:       "",
			Usage:       "github client secret",
			EnvVars:     []string{"OAUTH2_PROXY_GITHUB_SECRET"},
			Destination: &cfg.GitHub.Secret,
		},
		&cli.BoolFlag{
			Name:        "oauth2-bitbucket",
			Value:       false,
			Usage:       "enable bitbucket provider",
			EnvVars:     []string{"OAUTH2_PROXY_BITBUCKET"},
			Destination: &cfg.Bitbucket.Enabled,
		},
		&cli.StringSliceFlag{
			Name:    "oauth2-bitbucket-org",
			Value:   &cli.StringSlice{},
			Usage:   "allowed organizations from bitbucket",
			EnvVars: []string{"OAUTH2_PROXY_BITBUCKET_ORGS"},
		},
		&cli.StringFlag{
			Name:        "oauth2-bitbucket-client",
			Value:       "",
			Usage:       "bitbucket client id",
			EnvVars:     []string{"OAUTH2_PROXY_BITBUCKET_CLIENT"},
			Destination: &cfg.Bitbucket.Client,
		},
		&cli.StringFlag{
			Name:        "oauth2-bitbucket-secret",
			Value:       "",
			Usage:       "bitbucket client secret",
			EnvVars:     []string{"OAUTH2_PROXY_BITBUCKET_SECRET"},
			Destination: &cfg.Bitbucket.Secret,
		},
	}
}

func serverBefore(cfg *config.Config) cli.BeforeFunc {
	return func(c *cli.Context) error {
		if len(c.StringSlice("proxy-endpoint")) > 0 {
			// StringSliceFlag doesn't support Destination
			cfg.Proxy.Endpoints = c.StringSlice("proxy-endpoint")
		}

		gothic.GetProviderName = func(r *http.Request) (string, error) {
			if provider := chi.URLParam(r, "provider"); provider != "" {
				return provider, nil
			}

			return "", fmt.Errorf("you must select a provider")
		}

		if len(c.StringSlice("oauth2-gitlab-org")) > 0 {
			// StringSliceFlag doesn't support Destination
			cfg.Gitlab.Orgs = c.StringSlice("oauth2-gitlab-org")
		}

		if cfg.Gitlab.Enabled {
			goth.UseProviders(
				gitlab.NewCustomisedURL(
					cfg.Gitlab.Client,
					cfg.Gitlab.Secret,
					fmt.Sprintf("%s%s/gitlab", cfg.Server.Host, cfg.Server.Root),
					fmt.Sprintf("%s/oauth/authorize", cfg.Gitlab.URL),
					fmt.Sprintf("%s/oauth/token", cfg.Gitlab.URL),
					fmt.Sprintf("%s/api/v3/user", cfg.Gitlab.URL),
				),
			)
		}

		if len(c.StringSlice("oauth2-github-org")) > 0 {
			// StringSliceFlag doesn't support Destination
			cfg.GitHub.Orgs = c.StringSlice("oauth2-github-org")
		}

		if cfg.GitHub.Enabled {
			goth.UseProviders(
				github.New(
					cfg.GitHub.Client,
					cfg.GitHub.Secret,
					fmt.Sprintf("%s%s/github", cfg.Server.Host, cfg.Server.Root),
				),
			)
		}

		if len(c.StringSlice("oauth2-bitbucket-org")) > 0 {
			// StringSliceFlag doesn't support Destination
			cfg.Bitbucket.Orgs = c.StringSlice("oauth2-bitbucket-org")
		}

		if cfg.Bitbucket.Enabled {
			goth.UseProviders(
				bitbucket.New(
					cfg.Bitbucket.Client,
					cfg.Bitbucket.Secret,
					fmt.Sprintf("%s%s/bitbucket", cfg.Server.Host, cfg.Server.Root),
				),
			)
		}

		return nil
	}
}

func serverAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		fwd, err := forward.New(
			forward.PassHostHeader(true),
		)

		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to initialize forwarder")

			return err
		}

		lb, err := roundrobin.New(fwd)

		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to initialize balancer")

			return err
		}

		proxy, err := buffer.New(
			lb,
			buffer.Retry(`IsNetworkError() && Attempts() < 3`),
		)

		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to initialize buffer")

			return err
		}

		for _, endpoint := range cfg.Proxy.Endpoints {
			parsed, err := url.Parse(endpoint)

			if err != nil {
				log.Warn().
					Err(err).
					Str("endpoint", endpoint).
					Msg("failed to parse endpoint")

				continue
			}

			lb.UpsertServer(parsed)
		}

		var gr run.Group

		{
			stop := make(chan os.Signal, 1)

			gr.Add(func() error {
				signal.Notify(stop, os.Interrupt)

				<-stop

				return nil
			}, func(err error) {
				close(stop)
			})
		}

		{
			server := &http.Server{
				Addr:         cfg.Server.Health,
				Handler:      router.Status(cfg),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			gr.Add(func() error {
				log.Info().
					Str("addr", cfg.Server.Health).
					Msg("starting status server")

				return server.ListenAndServe()
			}, func(reason error) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				if err := server.Shutdown(ctx); err != nil {
					log.Info().
						Err(err).
						Msg("failed to stop status server gracefully")

					return
				}

				log.Info().
					Err(reason).
					Msg("status server stopped gracefully")
			})
		}

		if cfg.Server.AutoCert {
			parsed, err := url.Parse(
				cfg.Server.Host,
			)

			if err != nil {
				log.Info().
					Err(err).
					Msg("failed to parse host")

				return err
			}

			manager := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(parsed.Host),
				Cache:      autocert.DirCache(path.Join(cfg.Server.Storage, "certs")),
			}

			{
				server := &http.Server{
					Addr:         httpAddr,
					Handler:      router.Redirect(cfg),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
				}

				gr.Add(func() error {
					log.Info().
						Str("addr", httpAddr).
						Msg("starting http server")

					return server.ListenAndServe()
				}, func(reason error) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						log.Info().
							Err(err).
							Msg("failed to stop http server gracefully")

						return
					}

					log.Info().
						Err(reason).
						Msg("http server stopped gracefully")
				})
			}

			{
				server := &http.Server{
					Addr:         httpsAddr,
					Handler:      router.Load(cfg, proxy),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					TLSConfig: &tls.Config{
						PreferServerCipherSuites: true,
						MinVersion:               tls.VersionTLS12,
						CurvePreferences:         curves(cfg),
						CipherSuites:             ciphers(cfg),
						GetCertificate:           manager.GetCertificate,
					},
				}

				gr.Add(func() error {
					log.Info().
						Str("addr", httpsAddr).
						Msg("starting https server")

					return server.ListenAndServeTLS("", "")
				}, func(reason error) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						log.Info().
							Err(err).
							Msg("failed to stop https server gracefully")

						return
					}

					log.Info().
						Err(reason).
						Msg("https server stopped gracefully")
				})
			}

			return gr.Run()
		} else if cfg.Server.Cert != "" && cfg.Server.Key != "" {
			cert, err := tls.LoadX509KeyPair(
				cfg.Server.Cert,
				cfg.Server.Key,
			)

			if err != nil {
				log.Info().
					Err(err).
					Msg("failed to load certificates")

				return err
			}

			{
				server := &http.Server{
					Addr:         cfg.Server.Public,
					Handler:      router.Redirect(cfg),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
				}

				gr.Add(func() error {
					log.Info().
						Str("addr", cfg.Server.Public).
						Msg("starting http server")

					return server.ListenAndServe()
				}, func(reason error) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						log.Info().
							Err(err).
							Msg("failed to stop http server gracefully")

						return
					}

					log.Info().
						Err(reason).
						Msg("http server stopped gracefully")
				})
			}

			{
				server := &http.Server{
					Addr:         cfg.Server.Secure,
					Handler:      router.Load(cfg, proxy),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					TLSConfig: &tls.Config{
						PreferServerCipherSuites: true,
						MinVersion:               tls.VersionTLS12,
						CurvePreferences:         curves(cfg),
						CipherSuites:             ciphers(cfg),
						Certificates:             []tls.Certificate{cert},
					},
				}

				gr.Add(func() error {
					log.Info().
						Str("addr", cfg.Server.Secure).
						Msg("starting https server")

					return server.ListenAndServeTLS("", "")
				}, func(reason error) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						log.Info().
							Err(err).
							Msg("failed to stop https server gracefully")

						return
					}

					log.Info().
						Err(reason).
						Msg("https server stopped gracefully")
				})
			}

			return gr.Run()
		}

		{
			server := &http.Server{
				Addr:         cfg.Server.Public,
				Handler:      router.Load(cfg, proxy),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			gr.Add(func() error {
				log.Info().
					Str("addr", cfg.Server.Public).
					Msg("starting http server")

				return server.ListenAndServe()
			}, func(reason error) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				if err := server.Shutdown(ctx); err != nil {
					log.Info().
						Err(err).
						Msg("failed to stop http server gracefully")

					return
				}

				log.Info().
					Err(reason).
					Msg("http server stopped gracefully")
			})
		}

		return gr.Run()
	}
}

func curves(cfg *config.Config) []tls.CurveID {
	if cfg.Server.StrictCurves {
		return []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		}
	}

	return nil
}

func ciphers(cfg *config.Config) []uint16 {
	if cfg.Server.StrictCiphers {
		return []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		}
	}

	return nil
}
