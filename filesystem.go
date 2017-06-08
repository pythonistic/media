package media

import (
	"os"
	"net/http"
)

type MediaFileSystem struct {
	Root string
}

func (mfs MediaFileSystem) Open(name string) (file http.File, err error) {
	file, err = os.Open(name)
	return
}
