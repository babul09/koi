package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var file string
		switch r.URL.Path {
		case "/", "/index.html":
			file = "index.html"
		case "/game.wasm":
			w.Header().Set("Content-Type", "application/wasm")
			file = "game.wasm"
		case "/wasm_exec.js":
			file = "wasm_exec.js"
		default:
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, file)
	})
	log.Println("Motion DNA at http://localhost:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}
