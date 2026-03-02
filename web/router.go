package web

import (
	"mime"
	"net/http"

	ws "github.com/andeya/pholcus/common/websocket"
)

func init() {
	mime.AddExtensionType(".css", "text/css")
}

// Router registers HTTP and WebSocket routes.
func Router() {
	http.Handle("/ws", ws.Handler(wsHandle))
	http.Handle("/ws/log", ws.Handler(wsLogHandle))
	http.HandleFunc("/", web)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.FS(viewsSubFS()))))
}
