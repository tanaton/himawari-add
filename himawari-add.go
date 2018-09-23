package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)
}

func main() {
	err := run()
	if err != nil {
		log.Warn(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	if len(os.Args) < 3 {
		return errors.New("引数がおかしい")
	}
	p := os.Args[1]

	buf := bytes.Buffer{}
	w := multipart.NewWriter(&buf)
	_, filename := filepath.Split(p)
	ferr := w.WriteField("filename", filename)
	if ferr != nil {
		return ferr
	}
	w.Close() // 閉じることでPOSTデータが出来上がる模様

	res, err := http.Post("http://127.0.0.1:10616/task/add", w.FormDataContentType(), &buf)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}
	var data interface{}
	if err = json.Unmarshal([]byte(os.Args[2]), &data); err != nil {
		return err
	}
	log.WithField("data", data).Info("タスク追加に成功しました。")
	return nil
}
