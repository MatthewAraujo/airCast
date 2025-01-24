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
	"sync"

	"github.com/MatthewAraujo/airCast/utils"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

type Handler struct {
	conns map[*websocket.Conn]bool
	mu    sync.Mutex
}

func NewHandler() *Handler {
	return &Handler{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/video/{id}/stream", h.handleVideoStream).Methods(http.MethodGet)
	router.Handle("/ws", websocket.Handler(h.handleWS))

	absPath, err := filepath.Abs("./public")
	if err != nil {
		log.Fatalf("Erro ao obter o caminho absoluto: %v", err)
	}

	fileServer := http.FileServer(http.Dir(absPath))
	router.PathPrefix("/").Handler(http.StripPrefix("/api/v1", fileServer))
}

func (h *Handler) handleWS(ws *websocket.Conn) {
	log.Print("new incomming connection from cliente: ", ws.RemoteAddr())

	h.conns[ws] = true

	h.readLoop(ws)
}

func (h *Handler) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	for {
		n, err := ws.Read(buf)

		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
				break
			}

			log.Println("Read error:", err)
			break
		}

		msg := buf[:n]

		h.broadcastToWS(msg)
	}

	h.cleanupConnection(ws)
}

func (h *Handler) cleanupConnection(ws *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conns, ws)
	ws.Close()
	log.Println("Connection removed")
}

func (h *Handler) broadcastToWS(msg []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.conns {
		go func(conn *websocket.Conn, msg []byte) {
			_, err := conn.Write(msg)
			if err != nil {
				log.Println("Write error:", err)
			}
		}(conn, msg)
	}
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
		contentRange := fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize)

		w.Header().Set("Content-Range", contentRange)
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
			if err != nil && err == io.EOF {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("not able to get video info"))
				return
			}
			if n == 0 {
				break
			}
			w.Write(buffer[:n])
			bytesRead += int64(n)
		}

		h.broadcastToWS([]byte(utils.Int64ToString(bytesRead)))
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
