package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

func FilerUploadHandler(w http.ResponseWriter, r *http.Request) {
	//随机生成一个不存在的fileid
	var imgid string
	for {
		imgid = MakeImageID()
		if !FileExist(ImageID2Path(imgid)) {
			break
		}
	}
	//上传参数为uploadfile
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Upload Error."))
		return
	}
	defer file.Close()
	//检测文件类型
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Upload Error."))
		return
	}
	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" {
		w.Write([]byte("Error:Not JPEG."))
		return
	}
	//回绕文件指针
	log.Println(filetype)
	if _, err = file.Seek(0, 0); err != nil {
		log.Println(err)
	}
	//提前创建整棵存储树
	if err = BuildTree(imgid); err != nil {
		log.Println(err)
	}
	//log.Println(ImageID2Path(imgid))
	f, err := os.OpenFile(ImageID2Path(imgid), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Save Error."))
		return
	}
	defer f.Close()
	io.Copy(f, file)
	w.Write([]byte(imgid))
}

type Base64Content struct {
	Base64 string
}

func Base64UploadHandler(w http.ResponseWriter, r *http.Request) {
	var imgid string
	for {
		imgid = MakeImageID()
		if !FileExist(ImageID2Path(imgid)) {
			break
		}
	}
	if r.Body == nil {
		fmt.Fprintf(w, "Please send a request body:%d", 400)
		return
	}
	var err error
	var content Base64Content

	err = json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	bytes, err := base64.StdEncoding.DecodeString(content.Base64)
	if err != nil {
		fmt.Fprintf(w, "Base64 Content Error: %v", err)
		return
	}
	if err = BuildTree(imgid); err != nil {
		log.Println(err)
	}
	//log.Println(ImageID2Path(imgid))
	f, err := os.OpenFile(ImageID2Path(imgid), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Save Error."))
		return
	}
	defer f.Close()
	f.Write(bytes)

	w.Write([]byte(imgid))
	return
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageid := vars["imgid"]
	if len([]rune(imageid)) != 16 {
		w.Write([]byte("Error:ImageID incorrect."))
		return
	}
	imgpath := ImageID2Path(imageid)
	if !FileExist(imgpath) {
		w.Write([]byte("Error:Image Not Found."))
		return
	}
	http.ServeFile(w, r, imgpath)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><body><center><h1>It Works!</h1></center><hr><center>Quick Image Server</center></body></html>"))
}
