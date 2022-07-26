package controller

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// This is the root directory of uploaded files
var base = "/home/kjlb/src/webdav/dav/data"

func Upload(file *multipart.FileHeader, filePath string) (string, error) {

	// filePathのディレクトリが存在しない場合は作成する
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.MkdirAll(base+"/"+filePath, 0755)
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	n := fmt.Sprintf("%d-%s", time.Now().UTC().Unix(), file.Filename)
	fmt.Println(base + filePath + n)
	dst := fmt.Sprintf("%s/%s/%s", base, filePath, n)
	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)

	return n, err
}

func Download(n string) (string, []byte, error) {
	dst := fmt.Sprintf("%s/%s", base, n)
	b, err := ioutil.ReadFile(dst)
	if err != nil {
		return "", nil, err
	}
	m := http.DetectContentType(b[:512])

	return m, b, nil
}

// サーバーのファイル一覧を返却する
func List() ([]string, error) {
	files, err := ioutil.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names, nil
}
