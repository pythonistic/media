package media

import (
	"os"
	"path/filepath"
	"strings"
	"sort"
	"fmt"
)

var artists *Artists = new(Artists)

func init() {
	scanMedia(FSPATH_MEDIA)
	sortMedia()
}

func scanMedia(root string) {
	// scan the media
	mediaPath, err := os.Open(root)
	if err != nil {
		os.Stderr.WriteString("ERROR: failed to open media path " + root + ": " + err.Error() + "\n")
	}
	files, err := mediaPath.Readdir(0)
	if err != nil {
		os.Stderr.WriteString("ERROR: listing media path " + root + ": " + err.Error() + "\n")
	}

	for _, file := range files {
		if !file.IsDir() {
			fname := file.Name()
			fpath := filepath.Join(root, fname)

			var artist *Artist
			var album *Album
			var track *Track

			switch {
			case strings.HasSuffix(fname, ".mp3"):
				artist, album, track = parseMp3File(fname, fpath)
			default:
				println("WARN: Unparseable file: " + fpath)
				continue
			}

			foundArtist, ok := artists.Contains(artist)
			if ok {
				artist = foundArtist
			} else {
				*artists = append(*artists, artist)
			}

			// does the album exist?
			foundAlbum, ok := artist.Albums.Contains(album)
			if ok {
				album = foundAlbum
			} else {
				*artist.Albums = append(*artist.Albums, album)
			}

			// append the track
			track.Filename = fname
			track.Path = fpath
			*album.Tracks = append(*album.Tracks, track)
		} else {
			scanMedia(filepath.Join(root, file.Name()))
		}
	}
}

// sortMedia() sorts the media into output order.
func sortMedia() {
	sort.Sort(artists)

	for _, artist := range *artists {
		sort.Sort(artist.Albums)
		for _, album := range *artist.Albums {
			sort.Sort(*album.Tracks)
		}
	}

fmt.Printf("outContent: %v\n", artists)

	return
}
