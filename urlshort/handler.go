package urlshort

import (
	"net/http"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.String()

		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, url, 301)
		}

        fallback.ServeHTTP(w, r)
    }
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

type Route struct {
	Path, Url string
}

func parseYaml(yml []byte) ([]Route, error) {
	var routes []Route
	if err := yaml.Unmarshal(yml, &routes); err != nil {
		return nil, err
	}
	return routes, nil
}

func buildMap(routes []Route) map[string]string {
	routeMap := make(map[string]string)
	for _, route := range routes {
		routeMap[route.Path] = route.Url
	}
	return routeMap
}