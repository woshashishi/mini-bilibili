package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// VideoTranscoder 将上传文件规范为浏览器常用的 H.264/AAC MP4。
// 当前采用同步处理，最多同时转一个文件；不影响普通查询请求。
type VideoTranscoder struct {
	binary string
	slots  chan struct{}
}

func NewVideoTranscoder(binary string) (*VideoTranscoder, error) {
	resolved, err := exec.LookPath(binary)
	if err != nil {
		return nil, fmt.Errorf("find ffmpeg: %w", err)
	}
	return &VideoTranscoder{binary: resolved, slots: make(chan struct{}, 1)}, nil
}

func (t *VideoTranscoder) Normalize(parent context.Context, src io.ReadSeeker) (*os.File, func(), error) {
	ctx, cancel := context.WithTimeout(parent, 2*time.Minute)
	defer cancel()
	select {
	case t.slots <- struct{}{}:
		defer func() { <-t.slots }()
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
	dir, err := os.MkdirTemp("", "mini-bilibili-video-*")
	if err != nil {
		return nil, nil, err
	}
	success := false
	defer func() {
		if !success {
			_ = os.RemoveAll(dir) // 仅清理本次创建的临时目录。
		}
	}()
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, nil, err
	}
	inputPath := filepath.Join(dir, "input.mp4")
	input, err := os.Create(inputPath)
	if err != nil {
		return nil, nil, err
	}
	_, copyErr := io.Copy(input, src)
	closeErr := input.Close()
	if copyErr != nil {
		return nil, nil, copyErr
	}
	if closeErr != nil {
		return nil, nil, closeErr
	}
	outputPath := filepath.Join(dir, "output.mp4")
	cmd := exec.CommandContext(ctx, t.binary,
		"-nostdin", "-hide_banner", "-loglevel", "error", "-xerror",
		"-threads", "1", "-protocol_whitelist", "file,pipe", "-i", inputPath,
		"-map", "0:v:0", "-map", "0:a:0?", "-map_metadata", "-1",
		"-filter_threads", "1", "-vf", "scale=trunc(iw/2)*2:trunc(ih/2)*2",
		"-c:v", "libx264", "-threads", "1", "-profile:v", "baseline", "-pix_fmt", "yuv420p",
		"-preset", "fast", "-crf", "23", "-c:a", "aac", "-b:a", "128k",
		"-movflags", "+faststart", "-n", outputPath,
	)
	// 不向客户端暴露本机路径或转码工具输出。
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return nil, nil, fmt.Errorf("transcode timeout or cancellation: %w", ctx.Err())
		}
		return nil, nil, fmt.Errorf("transcode video: %w", err)
	}
	output, err := os.Open(outputPath)
	if err != nil {
		return nil, nil, err
	}
	success = true
	cleanup := func() {
		_ = output.Close() // Windows 下先关闭，再删除临时文件。
		_ = os.RemoveAll(dir)
	}
	return output, cleanup, nil
}
