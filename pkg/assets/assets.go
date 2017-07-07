package assets

import (
	"net/http"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/com"
	"github.com/webhippie/oauth2-proxy/pkg/config"
)

//go:generate fileb0x ab0x.yaml

// Load initializes the static files.
func Load() http.FileSystem {
	return ChainedFS{}
}

// ChainedFS is a simple HTTP filesystem including custom path.
type ChainedFS struct {
}

// Open just implements the HTTP filesystem interface.
func (c ChainedFS) Open(origPath string) (http.File, error) {
	if config.Server.Assets != "" {
		if com.IsDir(config.Server.Assets) {
			customPath := path.Join(config.Server.Assets, origPath)

			if com.IsFile(customPath) {
				f, err := os.Open(customPath)

				if err != nil {
					return nil, err
				}

				return f, nil
			}
		} else {
			logrus.Warnf("Custom assets directory doesn't exist")
		}
	}

	f, err := FS.OpenFile(CTX, origPath, os.O_RDONLY, 0644)

	if err != nil {
		return nil, err
	}

	return f, nil
}
