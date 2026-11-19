# 本地视频转码工具

Windows开发启动默认使用tools/ffmpeg.exe。该二进制仅保存在本机，已加入.gitignore。换电脑后需自行安装包含libx264、AAC编码器的FFmpeg，复制为此路径或用FFMPEG_PATH指定可执行文件。Linux默认从PATH查找ffmpeg。

上传时同步转码，最多一个转码任务；排队与处理合计最长两分钟。失败不生成正式视频。大文件后续可迁移到后台队列。

真实HEVC回归测试：设置VIDEO_TEST_INPUT为HEVC视频路径，VIDEO_TEST_FFMPEG为FFmpeg绝对路径，再运行go test ./internal/upload。未提供测试素材时跳过编解码集成测试。
