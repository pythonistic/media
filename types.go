package media


const TAG_TRACK = "TRCK"

type MediaContent struct {
	Filename string
	Path string
	Artist string
	Album string
	Title string
	Duration string
}

type PageContext struct {
	MediaList map[string]*MediaContent
}