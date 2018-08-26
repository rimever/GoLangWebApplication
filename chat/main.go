package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// templateは1つのテンプレートを表します。
type templateHandler struct {
	once     sync.Once
	filename string
	template *template.Template
}

// ServeHTTPはHTTPリクエストを処理します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		currentDirectoryPath, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		t.template = template.Must(template.ParseFiles(filepath.Join(currentDirectoryPath+"/chat/templates", t.filename)))
	})
	err := t.template.Execute(w,r)
	if err != nil {
		log.Fatal(err)
	}
}



func main() {
	var addr = flag.String("addr",":8080","アプリケーションのアドレス")
	flag.Parse()
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room",r)
	//チャットルームを開始します
	go r.run()
	// Webサーバー開始します
	log.Println("Webサーバーを開始します。ポート:",*addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
