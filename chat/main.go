package main

import (
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
	// 戻り値はチェックすべき
	result := t.template.Execute(w, nil)
	log.Print(result)
}

func main() {
	// ルート
	http.Handle("/", &templateHandler{filename: "chat.html"})
	// Webサーバー開始します
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
