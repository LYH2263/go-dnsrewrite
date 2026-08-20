package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/dnsrewrite"
	"example.com/dnsrewrite/internal/api"
)

func main() {
	addr := flag.String("addr", ":8097", "HTTP 监听地址")
	web := flag.String("web", "web", "静态管理页目录")
	persist := flag.String("persist", "", "规则快照 JSON 路径（可选）")
	up := flag.String("upstream", "", "默认上游地址 host:port")
	timeout := flag.Duration("upstream-timeout", 5*time.Second, "上游超时")
	flag.Parse()

	opts := []dnsrewrite.Option{
		dnsrewrite.WithUpstreamTimeout(*timeout),
	}
	if *persist != "" {
		opts = append(opts, dnsrewrite.WithPersistPath(*persist))
	}
	if *up != "" {
		opts = append(opts, dnsrewrite.WithDefaultUpstream(*up))
	}

	eng := dnsrewrite.New(opts...)
	defer eng.Close()
	if *persist != "" {
		if err := eng.LoadPersist(); err != nil {
			log.Printf("load persist: %v", err)
		}
	}

	srv := api.New(eng, api.Options{WebDir: *web, AllowCORS: true})
	httpSrv := &http.Server{
		Addr:              *addr,
		Handler:           srv,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("dnsd 监听 %s", *addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	_ = httpSrv.Close()
}
