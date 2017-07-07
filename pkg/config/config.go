package config

type server struct {
	Addr       string
	Cert       string
	Key        string
	Templates  string
	Assets     string
	Endpoint   string
	Title      string
	Pprof      bool
	Prometheus bool
}

type oauth2 struct {
	UserHeader string
}

type github struct {
	Enabled    bool
	Orgs       []string
	Client     string
	Secret     string
	URL        string
	SkipVerify bool
}

type gitlab struct {
	Enabled    bool
	Orgs       []string
	Client     string
	Secret     string
	URL        string
	SkipVerify bool
}

type bitbucket struct {
	Enabled bool
	Orgs    []string
	Client  string
	Secret  string
}

var (
	// Debug represents the flag to enable or disable debug logging.
	Debug bool

	// Server represents the information about the server bindings.
	Server = &server{}

	// OAuth2 represents the general configuration for OAuth2 bindings.
	OAuth2 = &oauth2{}

	// GitHub represents the information about the GitHub OAuth2 bindings.
	GitHub = &github{}

	// Gitlab represents the information about the Gitlab OAuth2 bindings.
	Gitlab = &gitlab{}

	// Bitbucket represents the information about the Bitbucket OAuth2 bindings.
	Bitbucket = &bitbucket{}
)
