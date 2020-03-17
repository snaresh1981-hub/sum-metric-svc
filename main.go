package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
	"log"
	"net/http"
	"os"
)

var (
	svc   = new(SumMetricService)
	cache = Cache{}
	router *mux.Router
	port   string
)

func GetEnv(key string, fallback string) string {
	v := os.Getenv(key)
	if len(v) != 0 {
		return v
	}
	return fallback
}


func AddSecurityHeaderMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Cache-Control", "no-cache, no-store")
		w.Header().Set("Expires", "0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")

		secureMiddleware := secure.New(secure.Options{
			FrameDeny:             true,
			ContentTypeNosniff:    true,
			BrowserXssFilter:      true,
			ContentSecurityPolicy: "default-src 'self'",
			ReferrerPolicy:        "same-origin",
		})

		secureMiddleware.Handler(next).ServeHTTP(w, r)

	})
}

func initRouter() {
	router.Handle("/", http.HandlerFunc(svc.getInfo)).Methods("GET")
	router.Handle("/metric/{key}/sum", http.HandlerFunc(svc.getMetricsSum)).Methods("GET")
	router.Handle("/metric/{key}", http.HandlerFunc(svc.postMetricData)).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(svc.notFound)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})
	h := AddSecurityHeaderMiddleWare(c.Handler(router))
	http.Handle("/", h)
	port = GetEnv("SUM_METRIC_SERVICE_PORT", "8080")
	log.Printf("Sum Metric Service listening on PORT %s", port)
}

func main() {
	router = mux.NewRouter()
	initRouter()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
