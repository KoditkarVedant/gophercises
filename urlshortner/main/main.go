package main

import (
	"fmt"
	"net/http"

	"github.com/KoditkarVedant/gophercises/urlshortner"
)

func main() {
	server := server()

	urlMap := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortner.MapHandler(urlMap, server)

	yamlUrlMap := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

	yamlHandler, err := urlshortner.YAMLHandler([]byte(yamlUrlMap), mapHandler)

	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":8081", yamlHandler)
}

func server() *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Running...")
	})
	return server
}
