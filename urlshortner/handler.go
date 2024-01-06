package urlshortner

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will redirect the url
func MapHandler(pathToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if find matching url

		path, ok := pathToUrls[r.URL.Path]

		if ok {
			redirectHandler := http.RedirectHandler(path, 301)
			redirectHandler.ServeHTTP(w, r)

			// redirect
		} else {
			//else fallback
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse url in yaml and redirect if user calls them
func YAMLHandler(yamlData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var urls []YamlUrl
	err := yaml.Unmarshal(yamlData, &urls)

	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {}, err
	}

	urlMap := make(map[string]string)

	for _, url := range urls {
		_, ok := urlMap[url.Path]
		if !ok {
			urlMap[url.Path] = url.Url
		}
	}

	fmt.Printf("%v", urlMap)

	return MapHandler(urlMap, fallback), nil
}

type YamlUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}
