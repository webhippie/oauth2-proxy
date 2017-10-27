package config

type server struct {
	Host          string
	Addr          string
	Cert          string
	Key           string
	Root          string
	Storage       string
	Templates     string
	Assets        string
	Endpoint      string
	Title         string
	LetsEncrypt   bool
	StrictCurves  bool
	StrictCiphers bool
	Prometheus    bool
	Pprof         bool
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

var (
	// LogLevel defines the log level used by our logging package.
	LogLevel string

	// Server represents the information about the server bindings.
	Server = &server{}

	// OAuth2 represents the general configuration for OAuth2 bindings.
	OAuth2 = &oauth2{}

	// GitHub represents the information about the GitHub OAuth2 bindings.
	GitHub = &github{}

	// Gitlab represents the information about the Gitlab OAuth2 bindings.
	Gitlab = &gitlab{}
)
