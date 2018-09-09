package main

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
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

	res, err := http.Post("http://127.0.0.1:1029/task/add", w.FormDataContentType(), &buf)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}
	return nil
}
