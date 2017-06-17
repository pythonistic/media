package media

import (
	id3 "github.com/mikkyang/id3-go"
	"os"
	"github.com/tcolgate/mp3"
	"time"
	"github.com/mikkyang/id3-go/v2"
	"strings"
)

const TAG_ALBUM = "TALB"
const TAG_ARTIST = "TPE1"
const TAG_DISC = "TPOS" // ? Part of Set ?
const TAG_PICTURE = "APIC"
const TAG_TITLE = "TIT2"
const TAG_TRACK = "TRCK"

func parseMp3File(fname string, fpath string) (artist *Artist, album *Album, track *Track) {
	artist = &Artist{
		Albums: new(Albums),
	}
	album = &Album{
		Tracks: new(Tracks),
	}
	track = &Track{}

	// parse the header
	func() {
		mp3File, err := id3.Open(fpath)
		if err != nil {
			os.Stderr.WriteString("ERROR: couldn't open or parse " + fpath + ": " + err.Error() + "\n")
		} else {
			defer mp3File.Close()

			album.Name = strings.TrimSpace(mp3File.Album())
			artist.Name = strings.TrimSpace(mp3File.Artist())
			pictureFrame := mp3File.Frame(TAG_PICTURE)
			track.Title = strings.TrimSpace(mp3File.Title())
			track.Track = strings.TrimSpace(mp3File.Frame(TAG_TRACK).String())

			//switch i := pictureFrame.(type) {
			//case nil:
			//	println("No picture for " + track.Track)
			//case v2.ImageFrame:
			//	imageFrame := pictureFrame.(v2.ImageFrame)
			//	album.CoverArtBytes = imageFrame.Data()
			//	album.CoverArtMimeType = imageFrame.Encoding()
			//default:
			//	fmt.Printf("Picture not the expected type: %v\n", i)
			//	fmt.Printf("pictureFrame: %v\n", pictureFrame)
			//}
			if pictureFrame != nil {
				imageFrame := pictureFrame.(*v2.ImageFrame)
				album.CoverArtBytes = imageFrame.Data()
				album.CoverArtMimeType = imageFrame.Encoding()
				album.CoverArtPath = "/" + PATH_ALBUM_ART + "/" + artist.Name + "/" + album.Name
			}
		}
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