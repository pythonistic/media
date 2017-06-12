package media


type MediaContent struct {
	Filename        string
	Path            string
	Artist          string
	Album           string
	Title           string
	Duration        string
	Track           string
	Picture         []byte
	PictureMimeType string
}

type PageContext struct {
	MediaList map[string]*MediaContent
}