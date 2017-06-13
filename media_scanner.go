package media

import (
	"os"
	"path/filepath"
	"strings"
	"sort"
	"fmt"
)

func scanMedia(root string, artists Artists) {
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
				artists = append(artists, artist)
			}

			// does the album exist?
			foundAlbum, ok := artist.Albums.Contains(album)
			if ok {
				album = foundAlbum
			}

			// append the track
			album.Tracks = append(album.Tracks, track)
		} else {
			scanMedia(filepath.Join(root, file.Name()), artists)
		}
	}
}

// sortMedia() sorts the media into output order.
func sortMedia(artists Artists) {
	sort.Sort(artists)

	for _, artist := range artists {
		sort.Sort(artist.Albums)
		for _, album := range artist.Albums {
			sort.Sort(album.Tracks)
		}
	}

fmt.Printf("outContent: %v\n", artists)

	return
}
