package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	defer stop()

	origin, err := url.Parse("https://api.telegram.org/")
	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup

	proxy := httputil.NewSingleHostReverseProxy(origin)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		defer wg.Done()

		wg.Add(1)

		r.Host = r.URL.Host

		proxy.ServeHTTP(w, r)
	})

	go func() {
		if err := http.ListenAndServe(":3883", nil); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-ctx.Done()

	wg.Wait()
}
