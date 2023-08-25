package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	// 0をシテイルすることで利用可能なポート番号を自動で割り当てる
	l, err := net.Listen("tcp", "localhost:0")
	// キャンセル可能なcontextオブジェクトを作成
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTPサーバーを起動する
	eg.Go(func() error {
		return run(ctx, l)
	})
	// Getリクエストの送信
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	// リッスンしているポート番号をログに出力
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %+v", err)
	}

	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %s, got %s", want, got)
	}

	// run関数に終了通知を送信する
	cancel()
	// run関数の戻り値を検証する
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
