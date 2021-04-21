package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const maxAttachmentSize = 4 * 1024 * 1024

type AttachmentHandler struct {
	attachmentDirectory string
}

func (a *AttachmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		a.postAttachment(w, r)
	} else {
		a.getAttachment(w, r)
	}
}

func generateAttachmentIdBase() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

var attachmentIdRegex = regexp.MustCompile(`\d+(\.[a-zA-Z0-9]{1,4})?`)

func isAttachmentId(attachmentId string) bool {
	return attachmentIdRegex.MatchString(attachmentId)
}

func (a *AttachmentHandler) deleteOldAttachments() {
	cleanupTime := time.Now().Add(-48 * time.Hour)
	err := filepath.WalkDir(a.attachmentDirectory, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			info, err := d.Info()
			if nil != err {
				return err
			}
			if info.ModTime().Before(cleanupTime) {
				os.Remove(path)
			}
		}
		return nil
	})
	if nil != err {
		log.Println("Error during attachment cleanup: ", err)
	}

}

func (a *AttachmentHandler) postAttachment(w http.ResponseWriter, r *http.Request) {
	a.deleteOldAttachments()
	attachmentId := generateAttachmentIdBase()
	attachmentFormFile, _, err := r.FormFile("attachment")
	if nil != err {
		http.Error(w, fmt.Sprintf("No attachment file in request: %s", err), http.StatusBadRequest)
		return
	}
	defer attachmentFormFile.Close()

	header := make([]byte, 512)
	_, err = attachmentFormFile.Read(header)
	if nil != err {
		http.Error(w, "Cannot read attachment header", http.StatusBadRequest)
		return
	}

	mimeType := http.DetectContentType(header)
	switch mimeType {
	case "image/gif":
		attachmentId += ".gif"
	case "image/png":
		attachmentId += ".png"
	case "image/jpeg":
		attachmentId += ".jpeg"
	case "image/webp":
		attachmentId += ".webp"
	case "audio/mpeg":
		attachmentId += ".mp3"
	case "application/ogg":
		attachmentId += ".ogg"
	case "video/webm":
		attachmentId += ".webm"
	}

	attachmentDestFile, err := os.Create(a.attachmentDirectory + "/" + attachmentId)
	if nil != err {
		http.Error(w, "Cannot create attachment file", http.StatusInternalServerError)
		log.Println("Cannot create attachment file: ", err)
		return
	}

	attachmentFormFile.Seek(0, io.SeekStart)
	io.Copy(attachmentDestFile, io.LimitReader(attachmentFormFile, maxAttachmentSize))

	attachmentURL := url.URL{
		Scheme: remoteProtocol(r),
		Host:   remoteHost(r),
		Path:   "/gb2c/attachment/" + attachmentId,
	}

	w.Header().Set("Location", attachmentURL.String())
	w.WriteHeader(http.StatusCreated)

}

func (a *AttachmentHandler) getAttachment(w http.ResponseWriter, r *http.Request) {
	attachmentId := strings.TrimPrefix(r.URL.Path, "/gb2c/attachment/")
	if !isAttachmentId(attachmentId) {
		http.Error(w, "Attachment id is not valid", http.StatusBadRequest)
		return
	}
	attachmentFilename := a.attachmentDirectory + "/" + attachmentId
	if len(filepath.Ext(attachmentId)) == 0 {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+attachmentId)
	}
	http.ServeFile(w, r, attachmentFilename)
}

//RegisterAttachmentAPI setup http attachment API endpoints
func RegisterAttachmentAPI() {
	attachmentHandler := &AttachmentHandler{attachmentDirectory: os.Getenv("GB2C_ATTACHMENT_DIRECTORY")}
	if len(attachmentHandler.attachmentDirectory) == 0 {
		log.Println("GB2C_ATTACHMENT_DIRECTORY is not defined, attachments will be stored in temporary directory")
		attachmentHandler.attachmentDirectory, _ = os.MkdirTemp("", "gb2c_attachments")
	}

	http.Handle("/gb2c/attachment/", attachmentHandler)
}
