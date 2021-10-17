package file_upload

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func StartWatching(watchedDir string, endpoint string, token string, deleteFile bool) {
	sigChan := make(chan os.Signal, 1)
	closeChan := make(chan interface{}, 1)

	wg := new(sync.WaitGroup)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
		return
	}
	defer func(watcher *fsnotify.Watcher) {
		if err = watcher.Close(); err != nil {
			log.Println(err)
		}
	}(watcher)

	if err = watcher.Add(watchedDir); err != nil {
		log.Println(err)
		return
	}

	wg.Add(1)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				func() {
					if !ok {
						return
					}
					if event.Op != fsnotify.Write {
						return
					}
					wg.Add(1)
					go func() {
						defer wg.Done()
						time.Sleep(10 * time.Second)
						UploadFile(endpoint, event.Name, token, deleteFile)
						log.Println(event)
					}()

				}()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println(err)
			case _, ok := <-closeChan:
				if !ok {
					wg.Done()
					return
				}
			}
		}
	}()

	<-sigChan
	close(closeChan)

	wg.Wait()

	log.Println("Shutting down...")
}

func UploadFile(endpoint string, fileName string, token string, deleteFile bool) {
	url := endpoint + "/api/documents/post_document/"

	b := new(bytes.Buffer)
	multiPartWriter := multipart.NewWriter(b)

	var formFile io.Writer
	var err error
	var file *os.File
	c := &http.Client{}
	var req *http.Request
	var resp *http.Response

	if formFile, err = multiPartWriter.CreateFormFile("document", fileName); err != nil {
		log.Println(err)
		return
	}

	if file, err = os.Open(fileName); err != nil {
		log.Println(err)
		return
	}
	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			log.Println(err)
		}
	}(file)

	if _, err = io.Copy(formFile, file); err != nil {
		log.Println(err)
		return
	}

	//https://stackoverflow.com/questions/20205796/post-data-using-the-content-type-multipart-form-data
	if err = multiPartWriter.Close(); err != nil {
		log.Println(err)
		return
	}

	if req, err = http.NewRequest(http.MethodPost, url, b); err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept", "version=2")

	resp, err = c.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode == http.StatusOK && deleteFile {
		if err = os.Remove(fileName); err != nil {
			log.Println(err)
		}
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(resp)
	}
}
