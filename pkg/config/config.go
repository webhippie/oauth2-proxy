package config

// Server defines the server configuration.
type Server struct {
	Health        string
	Secure        string
	Public        string
	Host          string
	Root          string
	Cert          string
	Key           string
	AutoCert      bool
	StrictCurves  bool
	StrictCiphers bool
	Templates     string
	Assets        string
	Storage       string
}

// Logs defines the logging configuration.
type Logs struct {
	Level   string
	Colored bool
	Pretty  bool
}

// Proxy defines the proxy configuration.
type Proxy struct {
	Title      string
	Endpoints  []string
	UserHeader string
}

// Gitlab defines the gitlab configuration.
type Gitlab struct {
	Enabled    bool
	Orgs       []string
	Client     string
	Secret     string
	URL        string
	SkipVerify bool
}

// GitHub defines the github configuration.
type GitHub struct {
	Enabled bool
	Orgs    []string
	Client  string
	Secret  string
}

// Bitbucket defines the bitbucket configuration.
type Bitbucket struct {
	Enabled bool
	Orgs    []string
	Client  string
	Secret  string
}

// Config defines the general configuration.
type Config struct {
	Server    Server
	Logs      Logs
	Proxy     Proxy
	Gitlab    Gitlab
	GitHub    GitHub
	Bitbucket Bitbucket
}

// New prepares a new default configuration.
func New() *Config {
	return &Config{}
}
