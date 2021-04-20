package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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

func generateAttachmentId() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func cleanAttachmentId(str string) string {
	i, _ := strconv.ParseInt(str, 10, 64)
	return strconv.FormatInt(i, 10)
}

func (a *AttachmentHandler) postAttachment(w http.ResponseWriter, r *http.Request) {
	attachmentId := generateAttachmentId()
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
	case "image/gif", "image/png", "image/jpeg", "image/webp":
	default:
		http.Error(w, fmt.Sprint("Invalid attachment mime type: ", mimeType), http.StatusBadRequest)
		return
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
	attachmentId := cleanAttachmentId(strings.TrimPrefix(r.URL.Path, "/gb2c/attachment/"))
	attachmentFilename := a.attachmentDirectory + "/" + attachmentId
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
