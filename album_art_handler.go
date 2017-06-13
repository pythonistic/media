package media

import (
	"net/http"
	"os"
	"strings"
	"path/filepath"
	"io/ioutil"
	"strconv"
)

const PATH_ALBUM_ART = "albumArt"
const HEADER_CONTENT_TYPE = "Content-Type"
const HEADER_CONTENT_LENGTH = "Content-Length"

const NO_SONG_MIME = "image/png"
var noSongImage []byte

func init() {
	noSongPath := filepath.Join(PATH_STATIC, "images", "no-song.png")
	fi, err := os.Stat(noSongPath)
	if err != nil {
		os.Stderr.WriteString("ERROR: failed to stat() no-song.png: " + err.Error())
	} else {
		noSongImage = make([]byte, fi.Size())
		f, err := os.Open(noSongPath)
		if err != nil {
			os.Stderr.WriteString("ERROR: failed to open no-song.png: " + err.Error())
		} else {
			noSongImage, err = ioutil.ReadAll(f)
			if err != nil {
				os.Stderr.WriteString("ERROR: failed to read no-song.png: " + err.Error())
			}
		}
		f.Close()
	}
}

func AlbumArtHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		os.Stderr.WriteString("ERROR: Failed to parse album art request: " + err.Error())
	}

	// parse the path parts
	albumPath := r.URL.Path

	// strip the path up to albumArt
	if strings.Contains(albumPath, PATH_ALBUM_ART) {
		albumPath = albumPath[strings.Index(albumPath, PATH_ALBUM_ART) + len(PATH_ALBUM_ART) + 1:]
		// split the parts
		albumPathParts := strings.Split(albumPath, "/")
		if len(albumPathParts) == 2 {
			artistName := albumPathParts[0]
			albumName := albumPathParts[1]

			artist, ok := artists.GetByName(artistName)
			if ok {
				album, ok := artist.Albums.GetByName(albumName)

				if ok && (album.CoverArtMimeType != "" && len(album.CoverArtBytes) > 0) {
						serveImageBytes(w, album.CoverArtMimeType, album.CoverArtBytes)
						return
				}
			}

		}
		// any other size will be discarded
	}

	// fall back on no song image
	serveImageBytes(w, NO_SONG_MIME, noSongImage)
}

func serveImageBytes(w http.ResponseWriter, mime string, image []byte) {
	w.Header().Set(HEADER_CONTENT_TYPE, mime)
	w.Header().Set(HEADER_CONTENT_LENGTH, strconv.Itoa(len(image)))
	w.WriteHeader(http.StatusOK)
	w.Write(image)
}