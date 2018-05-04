package httputils

import (
	"net/http"
	"strings"
	"io"
	"fmt"
	"os"
	"bytes"
	"io/ioutil"
	"errors"
)

func HttpGet(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func HttpGet2(url string) (data []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		return ioutil.ReadAll(response.Body)
	}
	return nil, errors.New("http code != 200")
}

func HttpPost(url string, data string) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Sia-Agent")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return http.DefaultClient.Do(req)
}

//接收文件
func RecvUploadFile(r *http.Request, dst string) {
	reader, _ := r.MultipartReader();
	for {
		p, err := reader.NextPart();
		if err == io.EOF {
			break
		}
		if err != nil{
			return
		}

		name := p.FormName()
		if name == "" {
			continue
		}
		filename := p.FileName()
		fmt.Println(filename)

		_, hasContentTypeHeader := p.Header["Content-Type"]
		if !hasContentTypeHeader && filename == "" {
			var b bytes.Buffer
			_, err := io.CopyN(&b, p, 1024*40)
			if err != nil && err != io.EOF {
				return
			}
			r.MultipartForm.Value[name] = append(r.MultipartForm.Value[name], b.String())
		}

		dstFile, err := os.Create(dst)
		if err != nil {
			return
		}
		io.Copy(dstFile, p)
		dstFile.Close()
	}
}