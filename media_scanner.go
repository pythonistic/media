package media

import (
	"os"
	"path/filepath"
	"strings"
	"sort"
	"strconv"
	"regexp"
)

const RE_TRACK_NUMBER = `(\d+)(\D.*)?`

var id int
var reTrack *regexp.Regexp

func init() {
	// build the regexp first since scanMedia() depends on it
	var err error
	reTrack, err = regexp.Compile(RE_TRACK_NUMBER)
	if err != nil {
		os.Stderr.WriteString("ERROR: failed to compile track number regular expression: " + err.Error() + "\n")
	}
}

func PrepareMedia(path string) *Artists {
	artists := scanMedia(path)
	sortMedia(artists)
	return artists
}

func scanMedia(root string) *Artists {
	artists := new(Artists)
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
				id++
				artist.Id = strconv.Itoa(id)
				*artists = append(*artists, artist)
			}

			// does the album exist?
			foundAlbum, ok := artist.Albums.Contains(album)
			if ok {
				album = foundAlbum
			} else {
				id++
				album.Id = strconv.Itoa(id)
				*artist.Albums = append(*artist.Albums, album)
			}

			// append the track
			track.Filename = fname
			track.Path = fpath
			track.Id = strconv.Itoa(id)

			// clean up the track number
			track.Track = cleanUpTrack(len(*album.Tracks), track.Track)

			*album.Tracks = append(*album.Tracks, track)
		} else {
			scanMedia(filepath.Join(root, file.Name()))
		}
	}

	return artists
}

// sortMedia() sorts the media into output order.
func sortMedia(artists *Artists) {
	sort.Sort(artists)

	for _, artist := range *artists {
		sort.Sort(artist.Albums)
		for _, album := range *artist.Albums {
			sort.Sort(*album.Tracks)
		}
	}

	return
}

// cleanUpTrack() cleans the track number.
func cleanUpTrack(numberOfTracks int, track string) string {
	track = strings.TrimSpace(track)

	// no track number?  make one up
	if track == "" {
		return strconv.Itoa(numberOfTracks + 1)
	}

	// do we have a track number?
	if reTrack.MatchString(track) {
		// pull out the track number
		num := reTrack.FindStringSubmatch(track)[1]
		return num
	}

	// else we have some weird track, return it as is
	return track
}