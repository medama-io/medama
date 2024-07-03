package services

import (
	"hash/crc32"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-faster/errors"
	generate "github.com/medama-io/medama"
)

type SPAHandler struct {
	clientServer http.Handler
	indexFile    []byte
	indexETag    string
	fileETags    map[string]string
}

func SetupAssetHandler(mux *http.ServeMux) error {
	client, err := generate.SPAClient()
	if err != nil {
		return errors.Wrap(err, "failed to create spa client")
	}

	handler, err := NewSPAHandler(client)
	if err != nil {
		return err
	}

	mux.Handle("/", handler)
	return nil
}

func NewSPAHandler(client fs.FS) (*SPAHandler, error) {
	clientServer := http.FileServer(http.FS(client))

	indexFile, err := readFile(client, "index.html")
	if err != nil {
		return nil, errors.Wrap(err, "could not read index.html")
	}

	handler := &SPAHandler{
		clientServer: clientServer,
		indexFile:    indexFile,
		indexETag:    generateETag(indexFile),
		fileETags:    make(map[string]string),
	}

	if err := handler.precomputeFileETags(client); err != nil {
		return nil, errors.Wrap(err, "failed to precompute asset ETags")
	}

	return handler, nil
}

func (h *SPAHandler) precomputeFileETags(client fs.FS) error {
	return fs.WalkDir(client, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && (isAssetPath("/"+path) || isRootFile(path)) {
			content, err := readFile(client, path)
			if err != nil {
				return err
			}
			h.fileETags["/"+path] = generateETag(content)
		}
		return nil
	})
}

func (h *SPAHandler) serveFile(w http.ResponseWriter, r *http.Request, filePath string) {
	if etag, ok := h.fileETags[filePath]; ok {
		w.Header().Set("ETag", etag)

		// 1 year for most asset files.
		cacheControl := "public, max-age=31536000, immutable"

		// 24 hours for root favicon files and tracker script.
		if isRootFile(strings.TrimPrefix(filePath, "/")) {
			cacheControl = "public, max-age=86400, must-revalidate"
		}
		w.Header().Set("Cache-Control", cacheControl)

		// Return 304 if the ETag matches.
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		// Get Content-Type from file extension.
		contentType := mime.TypeByExtension(filepath.Ext(filePath))
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}

	}

	h.clientServer.ServeHTTP(w, r)
}

func (h *SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uPath := path.Clean(r.URL.Path)
	_, exists := h.fileETags[uPath]

	// Serve index.html to all routes that are not /api.
	if uPath == "/" || !exists {
		h.serveIndexHTML(w, r)
		return
	}

	h.serveFile(w, r, uPath)
}

func (h *SPAHandler) serveIndexHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("ETag", h.indexETag)

	if r.Header.Get("If-None-Match") == h.indexETag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	if _, err := w.Write(h.indexFile); err != nil {
		http.Error(w, "could not serve index.html", http.StatusInternalServerError)
	}
}

func isAssetPath(path string) bool {
	return strings.HasPrefix(path, "/assets/") ||
		strings.HasPrefix(path, "/favicon.ico") ||
		strings.HasPrefix(path, "/manifest")
}

func isRootFile(path string) bool {
	return !strings.Contains(path, "/") && path != "index.html"
}

func readFile(filesystem fs.FS, file string) ([]byte, error) {
	f, err := filesystem.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func generateETag(content []byte) string {
	crc := crc32.ChecksumIEEE(content)
	return strconv.FormatUint(uint64(crc), 16)
}
