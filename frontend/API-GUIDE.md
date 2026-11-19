# 前端操作怎样到达 Go 后端

这是前后端分离的单体项目。Vue负责收集输入和展示结果；Gin负责HTTP请求、身份验证和业务处理；MySQL保存用户、视频信息及点赞关系；视频文件保存在磁盘。浏览器不会直接连接MySQL，当前业务也没有依赖Redis。

## 一、请求的三个部分

接口不是一个网址，而是“HTTP方法 + 路径 + 请求/响应约定”。

- 路径参数：`/videos/3`里的3，Gin用`c.Param("id")`读取。
- 查询参数：`/feed?cursor=3&limit=8`，Gin用`c.DefaultQuery`读取。
- 请求体：注册、登录和发布用JSON；上传用multipart表单传二进制文件。
- 请求头：受保护请求带`Authorization: Bearer <Token>`，它不在URL或JSON里。

开发时浏览器请求前端所在地址的`/api/...`，Vite代理转发给`127.0.0.1:8000`。`/uploads`和`/healthz`也有代理。生产环境需要由反向代理承担转发，打包后的dist不包含Vite开发代理。

## 二、全部接口与页面入口

| 页面操作 | HTTP接口 | 需要Token | 核心结果 |
|---|---|---|---|
| 页脚检查服务 | GET /healthz | 否 | 服务返回ok；不代表数据库、磁盘或视频解码全部正常 |
| 登录弹窗切换注册 | POST /api/v1/auth/register | 否 | 创建用户，返回id和username |
| 登录 | POST /api/v1/auth/login | 否 | 返回access_token和token_type |
| 刷新恢复身份 | GET /api/v1/me | 是 | 返回当前用户id |
| 首页与加载更多 | GET /api/v1/feed?cursor=0&limit=8 | 否 | items、has_more、next_cursor |
| 打开真实视频卡片 | GET /api/v1/videos/:id | 否 | 视频最新详情和播放地址 |
| 视频弹窗点赞数量 | GET /api/v1/videos/:id/likes/count | 否 | count |
| 视频弹窗点赞状态 | GET /api/v1/videos/:id/like | 是 | liked |
| 点赞 | PUT /api/v1/videos/:id/like | 是 | 将当前用户设为已点赞 |
| 取消点赞 | DELETE /api/v1/videos/:id/like | 是 | 将当前用户设为未点赞 |
| 投稿中上传文件 | POST /api/v1/uploads/videos | 是 | 服务器生成的filename |
| 投稿中发布作品 | POST /api/v1/videos | 是 | 新建视频记录的id |
| 播放与拖动进度 | GET /uploads/videos/文件名.mp4 | 否 | 视频字节，Range请求可返回206 |

## 三、后端各层分工

```text
Vue表单/按钮
  → Axios发送HTTP请求
  → Vite开发代理
  → Gin路由匹配方法和路径
  → RequireAuth（需要登录的接口）
  → Handler：解析参数、调用Service、将错误转换为HTTP响应
  → Service：业务校验、安排执行顺序
  → Repository：通过GORM操作MySQL
  → JSON响应
  → Vue更新响应式状态，页面重新渲染
```

上传的最后一层是LocalStorage写磁盘，静态播放直接读取文件，不经过视频Repository。

## 四、注册、登录、JWT

注册先校验用户名和密码，查询用户名是否存在，bcrypt生成带盐的密码哈希，最后写users表。数据库保存PasswordHash，不保存可直接读取的明文密码。bcrypt是密码哈希，不是能解密还原密码的加密。

登录通过用户名查出用户，再用bcrypt.CompareHashAndPassword比对输入密码和哈希。正确时Handler调用TokenManager.Generate，把user_id、签发者和过期时间放入Claims，用HS256及服务端密钥签名。

JWT通常可以被解码查看内容，签名用来验证它没有被篡改。密钥只在后端，不能放进Vue源码。前端保存Token，不需要也不应获得JWT_SECRET。

每次保护请求都由中间件拆出Bearer Token，Parse验证HS256、签名、issuer、exp及user_id，再通过c.Set(UserIDKey, id)传给Handler。发布作者、点赞用户由这一步确定，不能相信前端提交的user_id。

前端登录后调用/me确认身份，才显示登录状态；localStorage用于刷新后恢复。401会清除对应会话，网络故障保留凭证供重试。退出只删除浏览器凭证，已经签发的JWT仍会在到期前有效，当前没有服务端注销接口。

## 五、上传和发布为什么是两步

1. 选择文件后，前端创建FormData，将文件放入名为file的字段。浏览器自动生成multipart boundary，不要手动把Content-Type固定成application/json或缺boundary的multipart值。
2. 后端先限制整个请求体，再解析multipart。ParseMultipartForm的8MiB是文件内存缓冲阈值，不是允许上传的最大大小；较大文件会使用临时磁盘文件。
3. ValidateMP4检查大小、扩展名及开头内容的MIME。读取文件头后Seek回开头，避免保存时丢失最开始的字节。
4. FFmpeg先把上传文件转成H.264/AAC，成功后LocalStorage生成随机文件名并保存转换结果；出错会清理临时文件，不返回上传成功。上传完成返回filename。
5. 前端根据当前访问地址构造`/uploads/videos/filename`的完整HTTP地址，连同title和description提交发布接口。后端要求video_url为HTTP/HTTPS绝对地址。
6. 发布Service校验字段，作者ID取自JWT，Repository写入videos表。此时Feed才能查到这条作品。

上传文件和写视频记录不是一个数据库事务。上传成功但未发布，文件已经在磁盘上；当前后端没有清理孤立文件的接口。前端在发布失败时保留上传结果，不会强制重复上传。视频链接模式则跳过上传，直接发布已有HTTP/HTTPS链接。

开发时生成的完整地址包含当前访问的主机和端口。换部署域名后应规划统一媒体域名或后端存储键，不能认为localhost地址能被别人的电脑访问。

## 六、Feed分页

按id降序取数据，cursor=0代表第一页。Repository添加`id < cursor`条件，Service向数据库多取一条判断has_more，再截取用户要求的数量。

例如有ID 3、2、1，limit=2：先查3、2、1，返回3、2及next_cursor=2；下一页查询id<2，返回1，has_more=false。没有下一页时停止请求，否则cursor=0会回到第一页。

目前按ID倒序是最新发布流，没有个性化推荐算法。前端搜索也只过滤已加载视频，分区和全站搜索尚无后端接口。

## 七、点赞的四个接口

likes表的一行表示“某用户点赞了某视频”。联合唯一索引(user_id,video_id)阻止同一人重复生成多条记录。PUT重复提交遇到重复键时仍视为成功；DELETE匹配不到记录也视为成功，这就是当前点赞/取消的幂等行为。

liked回答“我是否点过”，count回答“总共有多少人点过”，两者需要不同查询。弹窗打开时查询详情、数量和当前用户状态；游客只查询公开内容。操作成功后再查询数量，前端不单纯依靠本地+1/-1。

查询失败不能当成0，也不能当成未点赞。页面显示重试，未确认状态时禁止继续切换。网络失败也不等于数据库没执行：请求可能已写入，只是响应丢失；重新读取状态才能确认。

## 八、播放为什么不返回JSON

Feed/详情给的是video_url字符串。浏览器的video元素拿到地址后会另发HTTP请求，后端静态路由把URL映射到storage/videos里的文件。Range指定要读取的字节区间，206表示分段返回，播放器可据此拖动进度。

200/206只证明文件传输成功。MP4是容器，内部还可能是H.264或H.265；浏览器不支持编码仍无法播放。上传流程现已增加FFmpeg自动转码为H.264/AAC。已有HTTP/HTTPS外部链接不会经过本机转码，仍需提供可播放的文件。

## 九、HTTP状态和业务code

HTTP状态描述请求结果，例如201创建成功、400参数错误、401身份无效、404资源不存在、409用户名冲突、413上传过大、500服务器异常。JSON里的code是项目自己的细分编号，成功统一为0，不能把HTTP200直接填成业务code200。

前端同时处理HTTP错误和业务错误：提示参数问题；401退出对应会话；服务或网络故障提供重试；提交过程中禁用按钮，避免连续重复点击。当前发布接口没有幂等键，网络超时后重新发布仍可能重复，后续可单独改进。

## 十、从哪些文件读起

- frontend/src/App.vue：登录注册、Feed、弹窗入口与发布后刷新。
- frontend/src/api/business.js：方法、路径、JSON/FormData、Token和401统一处理。
- frontend/src/composables/useAuth.js：凭证保存、恢复、退出。
- frontend/src/components/VideoDialog.vue：详情、播放、点赞四接口。
- frontend/src/components/PublishDialog.vue：选择文件、上传进度、发布表单。
- internal/server/router.go：入口路由，哪些挂了鉴权中间件。
- internal/middleware：Token到可信用户ID。
- internal/user、video、like：Handler→Service→Repository。
- internal/upload、storage：文件校验和磁盘保存。

## 更新：上传自动转码
上传现在先校验，再通过FFmpeg同步转换为H.264视频和AAC音频，成功后才保存并返回文件名。最多同时转换一个文件，排队和转码合计最多两分钟，请求取消会终止子进程；失败不生成正式文件，临时目录会清理。不会自动修改已发布的历史视频，也不转换外部链接。转换后文件哈希与原文件不同，属于正常现象。
本机工具位于tools/ffmpeg.exe，换电脑需准备带libx264和AAC编码器的FFmpeg；可用FFMPEG_PATH指定路径。当前仍是学习版同步方案，大视频后续可使用队列和后台Worker。
