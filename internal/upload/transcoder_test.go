package upload

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// 设置 VIDEO_TEST_INPUT 和 VIDEO_TEST_FFMPEG 后运行真实编解码回归。
func TestNormalizeHEVCAndCleanup(t *testing.T) {
	inputPath := os.Getenv("VIDEO_TEST_INPUT")
	binary := os.Getenv("VIDEO_TEST_FFMPEG")
	if inputPath == "" || binary == "" {
		t.Skip("set VIDEO_TEST_INPUT and VIDEO_TEST_FFMPEG for codec integration test")
	}
	transcoder, err := NewVideoTranscoder(binary)
	if err != nil {
		t.Fatal(err)
	}
	input, err := os.Open(inputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	output, cleanup, err := transcoder.Normalize(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	data, err := io.ReadAll(output)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte("avc1")) {
		t.Fatal("output is not H.264")
	}
	dir := filepath.Dir(output.Name())
	cleanup()
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("temporary directory was not removed")
	}
	_, _, err = transcoder.Normalize(context.Background(), bytes.NewReader([]byte("invalid media")))
	if err == nil {
		t.Fatal("invalid media accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _, err = transcoder.Normalize(ctx, input)
	if err == nil {
		t.Fatal("canceled conversion succeeded")
	}
}
