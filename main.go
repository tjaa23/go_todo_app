package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}

	// コマンド実行時の引数でポート番号を受け取る
	p := os.Args[1]
	l, err := net.Listen("tcp", "localhost:"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}

	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to run: %s", err)
	}
}

func run(ctx context.Context, l net.Listener) error {

	// HTTPリクエストを受け付ける
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTPサーバーを起動する
	eg.Go(func() error {
		// http.ErrServeClosedはshutdownが正常に終了したことを示すので異常ではない
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %s", err)
			return err
		}
		return nil
	})

	// チャネルからの終了通知を待機する
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %s", err)
	}

	// Goメソッドで起動した別ゴルーチンの終了を待つ
	return eg.Wait()
}
