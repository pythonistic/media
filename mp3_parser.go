package media

import (
	//"github.com/bogem/id3v2"
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

func parseMp3File(fname string, fpath string) (media *MediaContent) {
	media = &MediaContent{
		Filename: fname,
		Path: fpath,
	}

	// parse the header
	func() {
		//tag, err := id3v2.Open(fpath, id3v2.Options{Parse: true})
		mp3File, err := id3.Open(fpath)
		if err != nil {
			os.Stderr.WriteString("ERROR: couldn't open or parse " + fpath + ": " + err.Error() + "\n")
		} else {
			//defer tag.Close()
			defer mp3File.Close()

			//fmt.Printf("tag: %d %v\n", tag.Count(), tag)
			fmt.Printf("mp3file: %v\n", mp3File)

			//if tag.Count() > 0 {
				fmt.Printf("All Frmes: %v\n", mp3File.AllFrames())

				media.Album = mp3File.Album()
				media.Artist = mp3File.Artist()
				pictureFrame := mp3File.Frame(TAG_PICTURE)
				media.Title = mp3File.Title()
				media.Track = mp3File.Frame(TAG_TRACK).String()
				//albumFrame := tag.GetTextFrame(TAG_ALBUM)
				//artistFrame := tag.GetTextFrame(TAG_ARTIST)
				//pictureFrame := tag.GetLastFrame(TAG_PICTURE)
				//titleFrame := tag.GetTextFrame(TAG_TITLE)
				//trackFrame := tag.GetTextFrame(TAG_TRACK)
				//albumFrame := tag.GetLastFrame(TAG_ALBUM).(id3v2.TextFrame)
				//artistFrame := tag.GetLastFrame(TAG_ARTIST).(id3v2.TextFrame)
				//pictureFrame := tag.GetLastFrame(TAG_PICTURE)
				//titleFrame := tag.GetLastFrame(TAG_TITLE).(id3v2.TextFrame)
				//trackFrame := tag.GetLastFrame(TAG_TRACK).(id3v2.TextFrame)


				//fmt.Printf("Album %v %v\n", albumFrame, tag.Album())
				//fmt.Printf("Artist %v %v\n", artistFrame, tag.Artist())
				//fmt.Printf("Picture %v\n", pictureFrame)
				//fmt.Printf("Title %v %v\n", titleFrame, tag.Title())
				//fmt.Printf("Track %v\n", trackFrame)
				//
				//media.Album = strings.TrimSpace(albumFrame.Text)
				//media.Artist = strings.TrimSpace(artistFrame.Text)
				//media.Title = strings.TrimSpace(titleFrame.Text)
				//media.Track = strings.TrimSpace(trackFrame.Text)

				switch i := pictureFrame.(type) {
				case nil:
					println("No picture for " + media.Track)
				//case id3v2.PictureFrame:
				//	pf := pictureFrame.(id3v2.PictureFrame)
				//	media.Picture = pf.Picture
				//	media.PictureMimeType = pf.MimeType
				default:
					fmt.Printf("Picture not the expected type: %v\n", i)
					fmt.Printf("pictureFrame: %v\n", pictureFrame)
				}
				//if pictureFrame != nil {
				//	media.Picture = pictureFrame.(id3v2.PictureFrame).Picture
				//	media.PictureMimeType = pictureFrame.(id3v2.PictureFrame).MimeType
				//
				//}
			//}
		}

		fmt.Printf("Media: %v\n", media)
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
		media.Duration = trackDuration.String()
	}()

	return
}