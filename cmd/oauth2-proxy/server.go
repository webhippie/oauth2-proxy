package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/webhippie/oauth2-proxy/config"
	"github.com/webhippie/oauth2-proxy/router"
	"gopkg.in/urfave/cli.v2"
)

// Server provides the sub-command to start the server.
func Server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start the integrated server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "addr",
				Value:       "0.0.0.0:8080",
				Usage:       "Address to bind the server",
				EnvVars:     []string{"OAUTH2_PROXY_ADDR"},
				Destination: &config.Server.Addr,
			},
			&cli.BoolFlag{
				Name:        "pprof",
				Usage:       "Enable pprof debugging server",
				EnvVars:     []string{"OAUTH2_PROXY_PPROF"},
				Destination: &config.Server.Pprof,
			},
			&cli.StringFlag{
				Name:        "cert",
				Value:       "",
				Usage:       "Path to SSL cert",
				EnvVars:     []string{"OAUTH2_PROXY_CERT"},
				Destination: &config.Server.Cert,
			},
			&cli.StringFlag{
				Name:        "key",
				Value:       "",
				Usage:       "Path to SSL key",
				EnvVars:     []string{"OAUTH2_PROXY_KEY"},
				Destination: &config.Server.Key,
			},
			&cli.StringFlag{
				Name:        "templates",
				Value:       "",
				Usage:       "Path to custom templates",
				EnvVars:     []string{"OAUTH2_PROXY_TEMPLATES"},
				Destination: &config.Server.Templates,
			},
			&cli.StringFlag{
				Name:        "assets",
				Value:       "",
				Usage:       "Path to custom assets",
				EnvVars:     []string{"OAUTH2_PROXY_ASSETS"},
				Destination: &config.Server.Assets,
			},
			&cli.StringFlag{
				Name:        "endpoint",
				Value:       "",
				Usage:       "Endpoint to proxy requests to",
				EnvVars:     []string{"OAUTH2_PROXY_ENDPOINT"},
				Destination: &config.Server.Endpoint,
			},
			&cli.BoolFlag{
				Name:        "oauth2-github",
				Value:       false,
				Usage:       "Enable GitHub provider",
				EnvVars:     []string{"OAUTH2_PROXY_GITHUB_ENABLED"},
				Destination: &config.GitHub.Enabled,
			},
			&cli.StringSliceFlag{
				Name:    "oauth2-github-org",
				Value:   &cli.StringSlice{},
				Usage:   "Allowed organizations from GitHub",
				EnvVars: []string{"OAUTH2_PROXY_GITHUB_ORGS"},
			},
			&cli.StringFlag{
				Name:        "oauth2-github-client",
				Value:       "",
				Usage:       "GitHub client ID",
				EnvVars:     []string{"OAUTH2_PROXY_GITHUB_URL"},
				Destination: &config.GitHub.Client,
			},
			&cli.StringFlag{
				Name:        "oauth2-github-secret",
				Value:       "",
				Usage:       "GitHub client secret",
				EnvVars:     []string{"OAUTH2_PROXY_GITHUB_URL"},
				Destination: &config.GitHub.Secret,
			},
			&cli.StringFlag{
				Name:        "oauth2-github-url",
				Value:       "https://github.com",
				Usage:       "GitHub server URL",
				EnvVars:     []string{"OAUTH2_PROXY_GITHUB_URL"},
				Destination: &config.GitHub.URL,
			},
			&cli.BoolFlag{
				Name:        "oauth2-github-skipverify",
				Value:       false,
				Usage:       "Skip SSL verify for GitHub",
				EnvVars:     []string{"OAUTH2_PROXY_GITHUB_SKIPVERIFY"},
				Destination: &config.GitHub.SkipVerify,
			},
			&cli.BoolFlag{
				Name:        "oauth2-gitlab",
				Value:       false,
				Usage:       "Enable Gitlab provider",
				EnvVars:     []string{"OAUTH2_PROXY_GITLAB_ENABLED"},
				Destination: &config.Gitlab.Enabled,
			},
			&cli.StringSliceFlag{
				Name:    "oauth2-gitlab-org",
				Value:   &cli.StringSlice{},
				Usage:   "Allowed organizations from Gitlab",
				EnvVars: []string{"OAUTH2_PROXY_GITLAB_ORGS"},
			},
			&cli.StringFlag{
				Name:        "oauth2-gitlab-client",
				Value:       "",
				Usage:       "Gitlab client ID",
				EnvVars:     []string{"OAUTH2_PROXY_GITLAB_URL"},
				Destination: &config.Gitlab.Client,
			},
			&cli.StringFlag{
				Name:        "oauth2-gitlab-secret",
				Value:       "",
				Usage:       "Gitlab client secret",
				EnvVars:     []string{"OAUTH2_PROXY_GITLAB_URL"},
				Destination: &config.Gitlab.Secret,
			},
			&cli.StringFlag{
				Name:        "oauth2-gitlab-url",
				Value:       "https://gitlab.com",
				Usage:       "Gitlab server URL",
				EnvVars:     []string{"OAUTH2_PROXY_GITLAB_URL"},
				Destination: &config.Gitlab.URL,
			},
			&cli.BoolFlag{
				Name:        "oauth2-gitlab-skipverify",
				Value:       false,
				Usage:       "Skip SSL verify for Gitlab",
				EnvVars:     []string{"OAUTH2_PROXY_GITLAB_SKIPVERIFY"},
				Destination: &config.Gitlab.SkipVerify,
			},
			&cli.BoolFlag{
				Name:        "oauth2-bitbucket",
				Value:       false,
				Usage:       "Enable Bitbucket provider",
				EnvVars:     []string{"OAUTH2_PROXY_BITBUCKET_ENABLED"},
				Destination: &config.Bitbucket.Enabled,
			},
			&cli.StringSliceFlag{
				Name:    "oauth2-bitbucket-org",
				Value:   &cli.StringSlice{},
				Usage:   "Allowed organizations from Bitbucket",
				EnvVars: []string{"OAUTH2_PROXY_BITBUCKET_ORGS"},
			},
			&cli.StringFlag{
				Name:        "oauth2-bitbucket-client",
				Value:       "",
				Usage:       "Bitbucket client ID",
				EnvVars:     []string{"OAUTH2_PROXY_BITBUCKET_URL"},
				Destination: &config.Bitbucket.Client,
			},
			&cli.StringFlag{
				Name:        "oauth2-bitbucket-secret",
				Value:       "",
				Usage:       "Bitbucket client secret",
				EnvVars:     []string{"OAUTH2_PROXY_BITBUCKET_URL"},
				Destination: &config.Bitbucket.Secret,
			},
			&cli.StringFlag{
				Name:        "oauth2-user-header",
				Value:       "X-PROXY-USER",
				Usage:       "Header for username",
				EnvVars:     []string{"OAUTH2_PROXY_USER_HEADER"},
				Destination: &config.OAuth2.UserHeader,
			},
		},
		Before: func(c *cli.Context) error {
			if len(c.StringSlice("oauth2-github-org")) > 0 {
				// StringSliceFlag doesn't support Destination
				config.GitHub.Orgs = c.StringSlice("oauth2-github-org")
			}

			if len(c.StringSlice("oauth2-gitlab-org")) > 0 {
				// StringSliceFlag doesn't support Destination
				config.Gitlab.Orgs = c.StringSlice("oauth2-gitlab-org")
			}

			if len(c.StringSlice("oauth2-bitbucket-org")) > 0 {
				// StringSliceFlag doesn't support Destination
				config.Bitbucket.Orgs = c.StringSlice("oauth2-bitbucket-org")
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			logrus.Infof("Starting on %s", config.Server.Addr)

			cfg, err := ssl()

			if err != nil {
				return err
			}

			server := &http.Server{
				Addr:         config.Server.Addr,
				Handler:      router.Load(),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
				TLSConfig:    cfg,
			}

			if err := startServer(server); err != nil {
				logrus.Fatal(err)
			}

			return nil
		},
	}
}

func curves() []tls.CurveID {
	return []tls.CurveID{
		tls.CurveP521,
		tls.CurveP384,
		tls.CurveP256,
	}
}

func ciphers() []uint16 {
	return []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	}
}

func ssl() (*tls.Config, error) {
	if config.Server.Cert != "" && config.Server.Key != "" {
		cert, err := tls.LoadX509KeyPair(
			config.Server.Cert,
			config.Server.Key,
		)

		if err != nil {
			return nil, fmt.Errorf("Failed to load SSL certificates. %s", err)
		}

		return &tls.Config{
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         curves(),
			CipherSuites:             ciphers(),
			Certificates:             []tls.Certificate{cert},
		}, nil
	}

	return nil, nil
}
