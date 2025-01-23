package video

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/MatthewAraujo/airCast/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/video/{id}/stream", h.handleVideoStream).Methods(http.MethodGet)

	absPath, err := filepath.Abs("./public")
	if err != nil {
		log.Fatalf("Erro ao obter o caminho absoluto: %v", err)
	}

	fileServer := http.FileServer(http.Dir(absPath))
	router.PathPrefix("/").Handler(http.StripPrefix("/api/v1", fileServer))
}

func (h *Handler) handleVideoStream(w http.ResponseWriter, r *http.Request) {
	log.Print("getting hit")
	videoID := mux.Vars(r)["id"]
	videoPath := fmt.Sprintf("./videos/%s.mp4", videoID)

	file, err := os.Open(videoPath)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("video not found"))
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("not able to get video info"))
		return
	}
	fileSize := fileInfo.Size()

	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		rangeParts := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
		start, _ := strconv.ParseInt(rangeParts[0], 10, 64)
		var end int64 = fileSize - 1
		if len(rangeParts) > 1 && rangeParts[1] != "" {
			end, _ = strconv.ParseInt(rangeParts[1], 10, 64)
		}

		chunkSize := end - start + 1
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", chunkSize))
		w.WriteHeader(http.StatusPartialContent)

		file.Seek(start, 0)

		buffer := make([]byte, 1024*8) // Buffer de 8KB
		bytesRead := int64(0)
		for {
			if bytesRead >= chunkSize {
				break
			}
			n, err := file.Read(buffer)
			if err != nil && err.Error() != "EOF" {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("not able to get video info"))
				return
			}
			if n == 0 {
				break
			}
			w.Write(buffer[:n])
			bytesRead += int64(n)
		}
	} else {
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
		w.WriteHeader(http.StatusOK)

		_, err := file.Seek(0, 0)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error getting the start of the video"))
			return
		}
		_, err = io.Copy(w, file)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error sending the video"))
		}
	}
}
