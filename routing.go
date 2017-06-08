package media

import "net/http"

func ConstructRoutes() http.Handler {
	handler := http.DefaultServeMux

	handler.HandleFunc("/" + PATH_MEDIA, MediaHandler)
	handler.HandleFunc("/" + PATH_STATIC, StaticHandler)
	handler.HandleFunc("/" + PATH_IMAGES, StaticHandler)
	handler.HandleFunc("/" + PATH_ANDROID_HIRES, StaticHandler)
	handler.HandleFunc("/" + PATH_ANDROID_NORMAL, StaticHandler)
	handler.HandleFunc("/" + PATH_APPLE_PRECOMP, StaticHandler)
	handler.HandleFunc("/" + PATH_APPLE_TOUCH, StaticHandler)
	handler.HandleFunc("/" + PATH_FAVICON, StaticHandler)
	handler.HandleFunc("/", IndexHandler)

	return handler
}
