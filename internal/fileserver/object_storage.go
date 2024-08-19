package fileserver

import (
	"fmt"
	"net/http"
)

type ObjectStorageFileServer struct {
	ObjectStorageURL string
}

func (osfs *ObjectStorageFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileURL := fmt.Sprintf("%s%s", osfs.ObjectStorageURL, r.URL.Path)
	http.Redirect(w, r, fileURL, http.StatusFound)
}
