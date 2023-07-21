package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	if r := recover(); r != nil {
		f, err := os.Create(fmt.Sprint("log-", time.Now().UnixNano(), ".txt"))
		if err == nil {
			defer f.Close()
			f.WriteString(err.Error())
		}

		main()
	}

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
				panic(err)
			}
		}
	}()

	<-ctx.Done()

	wg.Wait()
}
