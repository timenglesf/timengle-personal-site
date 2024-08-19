package fileserver

import "net/http"

type EmbeddedFileServer struct {
	FS http.FileSystem
}

func (efs *EmbeddedFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.FileServer(efs.FS).ServeHTTP(w, r)
}
