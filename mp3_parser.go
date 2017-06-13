package media

import (
	id3 "github.com/mikkyang/id3-go"
	"os"
	"github.com/tcolgate/mp3"
	"time"
	"fmt"
)

const TAG_ALBUM = "TALB"
const TAG_ARTIST = "TPE1"
const TAG_DISC = "TPOS" // ? Part of Set ?
const TAG_PICTURE = "APIC"
const TAG_TITLE = "TIT2"
const TAG_TRACK = "TRCK"

func parseMp3File(fname string, fpath string) (artist *Artist, album *Album, track *Track) {
	artist = &Artist{
		Albums: make(Albums, 0),
	}
	album = &Album{
		Tracks: make(Tracks, 0),
	}
	track = &Track{}

	// parse the header
	func() {
		mp3File, err := id3.Open(fpath)
		if err != nil {
			os.Stderr.WriteString("ERROR: couldn't open or parse " + fpath + ": " + err.Error() + "\n")
		} else {
			defer mp3File.Close()

			album.Name = mp3File.Album()
			artist.Name = mp3File.Artist()
			pictureFrame := mp3File.Frame(TAG_PICTURE)
			track.Title = mp3File.Title()
			track.Track = mp3File.Frame(TAG_TRACK).String()

			switch i := pictureFrame.(type) {
			case nil:
				println("No picture for " + track.Track)
			default:
				fmt.Printf("Picture not the expected type: %v\n", i)
				fmt.Printf("pictureFrame: %v\n", pictureFrame)
			}
		}

		fmt.Printf("Artist=%v Album=%v Track=%v\n", artist, album, track)
	}()

	// parse the file content
	mp3f, err := os.Open(fpath)
	if err != nil {
		os.Stderr.WriteString("ERROR:  Couldn't open file " + fname + " for reading: %v" + err.Error() + "\n")
	}
	func() {
		skipped := 0
		defer func() {
			mp3f.Close()
		}()
		decoder := mp3.NewDecoder(mp3f)
		var trackDuration time.Duration
		var frame mp3.Frame
		for {
			if err := decoder.Decode(&frame, &skipped); err != nil {
				os.Stdout.WriteString("INFO: decode error: " + err.Error() + "\n")
				break
			}
			trackDuration += frame.Duration()
		}
		trackDuration = time.Duration(trackDuration.Seconds()) * time.Second
		track.Duration = trackDuration.String()
	}()

	return
}