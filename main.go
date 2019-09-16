package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Quick Image Server.")
	fmt.Println("Author:zjyl1994@outlook.com.\nAppender:boommanpro@gmail.com")
	LoadConf()
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/fileUpload", FilerUploadHandler).Methods("POST")
	r.HandleFunc("/base64Upload", Base64UploadHandler).Methods("POST")
	r.HandleFunc("/{imgid}", DownloadHandler).Methods("GET")
	fmt.Printf("server start success: bindAddress:%s ,storage: %s", conf.ListenAddr, conf.Storage)
	err := http.ListenAndServe(conf.ListenAddr, r)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
