<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { categories, videos } from './data/feed'
import { fetchFeedPage } from './api/feed'
import { createBusinessAPI, requestMessage } from './api/business'
import VideoDialog from './components/VideoDialog.vue'
import PublishDialog from './components/PublishDialog.vue'
import { useAuth } from './composables/useAuth'

const activeCategory = ref('全部')
const activeNav = ref('首页')
const query = ref('')
const search = ref('')
const feedVideos = ref([])
const feedLoading = ref(false)
const feedError = ref('')
const feedLoaded = ref(false)
const retryReset = ref(true)
const hasMore = ref(false)
const nextCursor = ref(0)
let controller
const selectedVideo = ref(null)
const dialog = ref('')
const toast = ref('')
let toastTimer
const loginName = ref('')
const loginPassword = ref('')
const { user: currentUser, busy: authBusy, error: loginError, notice: authNotice, signIn, restore: restoreSession, logout, getToken, expire } = useAuth()
const api = createBusinessAPI({ getToken, onUnauthorized: expire })
const authMode = ref('login')
const registerBusy = ref(false)
const confirmPassword = ref('')
const uploadBusy = ref(false)
const healthStatus = ref('检查服务')
const authWorking = computed(() => authBusy.value || registerBusy.value)
async function checkHealth() {
  healthStatus.value = '检查中…'
  try {
    const response = await fetch('/healthz', { signal: AbortSignal.timeout(5000) })
    const result = await response.json()
    healthStatus.value = response.ok && result.status === 'ok' ? '后端在线' : '服务异常'
  } catch { healthStatus.value = '后端暂不可用' }
}
function switchAuth(mode) { authMode.value = mode; loginError.value = ''; loginPassword.value = ''; confirmPassword.value = '' }
async function submitAuth() {
  if (authWorking.value) return
  if (authMode.value === 'login') { await submitLogin(); return }
  const name = loginName.value.trim()
  const password = loginPassword.value
  const size = new TextEncoder().encode(password).length
  if (!/^[A-Za-z0-9_]{3,32}$/.test(name)) { loginError.value = '用户名应为3至32位字母、数字或下划线'; return }
  if (size < 8 || size > 72) { loginError.value = '密码应为8至72字节'; return }
  if (password !== confirmPassword.value) { loginError.value = '两次输入的密码不一致'; return }
  registerBusy.value = true
  loginError.value = ''
  try {
    await api.register(name, password)
    loginName.value = name
    switchAuth('login')
    notify('注册成功，请登录')
  } catch (error) { loginError.value = requestMessage(error) }
  finally { registerBusy.value = false }
}
async function onPublished(id) {
  dialog.value = ''
  query.value = ''; search.value = ''
  notify('发布成功，视频 #' + id)
  await loadFeed(true)
}
async function submitLogin() {
  const success = await signIn(loginName.value, loginPassword.value)
  loginPassword.value = ''
  if (success) { dialog.value = ''; notify('登录成功') }
}
function signOut() { logout(); loginPassword.value = ''; notify('已退出登录') }
function openLogin() { selectedVideo.value = null; switchAuth('login'); dialog.value = 'login' }
watch(dialog, () => { loginPassword.value = ''; confirmPassword.value = ''; loginError.value = '' })
const featuredIndex = ref(0)
const features = [
  { image: 'mountain', tag: '生活，不止一种打开方式', title: '把向往的生活，\n拍成自己的故事。', subtitle: '记录热爱 · 分享日常 · 发现同频的你', color: '#294e43' },
  { image: 'city', tag: '城市漫游计划', title: '等一盏灯亮起，\n听一座城的心跳。', subtitle: '带上镜头，收藏每个闪光的瞬间', color: '#38394e' },
  { image: 'desk', tag: '今天也学到一点新东西', title: '让每一份好奇，\n都有新的回响。', subtitle: '一起探索，一起把想法变成现实', color: '#334956' },
]
const feature = computed(() => features[featuredIndex.value])
const displayedVideos = computed(() => feedVideos.value.filter(video =>
  !search.value || [video.title, video.description, String(video.author_id)].join(' ').toLowerCase().includes(search.value.toLowerCase())
))
async function loadFeed(reset = false) {
  if (feedLoading.value || (!reset && feedLoaded.value && !hasMore.value)) return
  feedLoading.value = true
  feedError.value = ''
  retryReset.value = reset
  controller = new AbortController()
  try {
    const page = await fetchFeedPage(reset ? 0 : nextCursor.value, 8, controller.signal)
    const incoming = page.items.map(video => ({ ...video, realVideo: true }))
    const previous = reset ? [] : feedVideos.value
    const seen = new Set(previous.map(video => video.id))
    feedVideos.value = [...previous, ...incoming.filter(video => !seen.has(video.id))]
    hasMore.value = page.has_more
    nextCursor.value = page.next_cursor
    feedLoaded.value = true
  } catch (error) {
    if (error.code !== 'ERR_CANCELED') feedError.value = '视频加载失败，请确认后端已启动后重试。'
  } finally {
    feedLoading.value = false
  }
}
function displayDate(value) {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? '' : date.toLocaleDateString('zh-CN')
}
onMounted(() => { loadFeed(true); restoreSession(); checkHealth() })
function notify(message) {
  toast.value = message
  clearTimeout(toastTimer)
  toastTimer = setTimeout(() => { toast.value = '' }, 3000)
}
function chooseCategory(category) { if (category !== '全部') { notify('分区筛选待接入，当前展示全部视频'); return }; activeCategory.value = category }
function runSearch() { search.value = query.value.trim() }
function setNav(nav) { activeNav.value = nav; if (nav === '首页') { query.value = ''; search.value = ''; chooseCategory('全部') } else notify(`${nav}页面即将上线`) }
function scrollToTop() { window.scrollTo({ top: 0, behavior: 'smooth' }) }
function closeOverlay() { if (uploadBusy.value) { notify('正在提交，请等待完成'); return }; selectedVideo.value = null; dialog.value = '' }
function onKeydown(event) { if (event.key === 'Escape') closeOverlay() }
window.addEventListener('keydown', onKeydown)
onBeforeUnmount(() => { window.removeEventListener('keydown', onKeydown); clearTimeout(toastTimer); controller?.abort() })
// 后端已有业务接口已接入；分类、收藏、历史等尚无对应接口。
</script>

<template>
  <div class="site">
    <header class="topbar">
      <button class="brand" aria-label="Mini Bili 首页" @click="setNav('首页')"><svg viewBox="0 0 48 42" fill="none" aria-hidden="true"><path d="m14 3 5 7M34 3l-5 7"/><rect x="5" y="11" width="38" height="27" rx="7"/><path d="m15 21 5 2m13-2-5 2m-9 8h10"/></svg><span>mini<span class="brand-bili">bili</span><small>干杯，生活！</small></span></button>
      <nav class="main-nav" aria-label="主导航"><button v-for="nav in ['首页', '番剧', '直播', '游戏中心', '会员购']" :key="nav" :class="{ active: activeNav === nav }" @click="setNav(nav)">{{ nav }}</button></nav>
      <form class="search-box" role="search" @submit.prevent="runSearch"><input v-model="query" placeholder="发现你感兴趣的视频" aria-label="搜索视频"/><button aria-label="搜索"><svg viewBox="0 0 24 24" fill="none"><circle cx="10.5" cy="10.5" r="6.5"/><path d="m16 16 4 4"/></svg></button></form>
      <div class="header-actions"><div class="account-controls"><template v-if="currentUser"><span class="account-id">用户 #{{ currentUser.id }}</span><button class="logout-button" :disabled="uploadBusy" @click="signOut">退出</button></template><button v-else class="login-avatar" :disabled="authBusy" @click="openLogin">{{ authBusy ? '验证中' : '登录' }}</button></div><button class="header-tool" @click="notify('消息中心即将上线')"><span>♧</span>消息</button><button class="header-tool" @click="notify('收藏夹即将上线')"><span>☆</span>收藏</button><button class="header-tool" @click="notify('观看历史即将上线')"><span>◷</span>历史</button><button class="upload-button" @click="dialog = 'upload'"><svg viewBox="0 0 24 24" fill="none"><path d="M12 16V4m-5 5 5-5 5 5M4 15v5h16v-5"/></svg>投稿</button></div>
    </header>

    <section class="sky-banner" aria-label="首页横幅"><div class="banner-cloud cloud-one"></div><div class="banner-cloud cloud-two"></div><div class="banner-copy"><span>在这里，每一份热爱都值得被看见</span><h1>生活有趣，一起 <em>发现</em>。</h1><p>你的日常，也可以是别人的灵感。</p></div><div class="banner-art" aria-hidden="true"><div class="orbit orbit-one"></div><div class="orbit orbit-two"></div><div class="float-card card-music">♫</div><div class="float-card card-heart">♥</div><div class="float-card card-play">▶</div><div class="tv"><div class="antenna left"></div><div class="antenna right"></div><div class="tv-screen"><i></i><i></i><b>⌣</b></div><div class="tv-foot"></div></div><span class="spark spark-one">✦</span><span class="spark spark-two">✧</span></div><span class="banner-note">HELLO, INSPIRATION.</span></section>

    <main class="page-content"><div v-if="authNotice" class="auth-notice" role="status"><span>{{ authNotice }}</span><button :disabled="authBusy" @click="restoreSession">重试</button><button @click="authNotice = ''" aria-label="关闭提示">×</button></div>
      <section class="channels" aria-label="视频频道"><div class="quick-channels"><button @click="notify('动态页面即将上线')"><span class="quick-icon pink">✿</span>动态</button><button @click="notify('热门排行待接入')"><span class="quick-icon orange">♨</span>热门</button><button @click="notify('分区目录即将上线')"><span class="quick-icon blue">▦</span>频道</button></div><div class="category-grid"><button v-for="category in categories.slice(1)" :key="category" :class="{ selected: activeCategory === category }" @click="chooseCategory(category)">{{ category }}<span v-if="category === '科技'" class="new-dot"></span></button><button @click="notify('更多分区正在准备中')">更多 <span>⌄</span></button></div><div class="channel-links"><button @click="notify('专栏即将上线')">▤ 专栏</button><button @click="notify('活动中心即将上线')">⚑ 活动</button><button @click="notify('社区即将上线')">♧ 社区</button><button @click="notify('新歌榜即将上线')">♫ 新歌</button></div></section>

      <p class="demo-note">精选区域 · 布局演示，真实视频见下方</p>
      <section class="discovery" aria-label="精选布局演示"><button class="feature-card" :style="{ backgroundColor: feature.color }" @click="selectedVideo = videos[featuredIndex === 0 ? 0 : featuredIndex === 1 ? 4 : 1]"><img :src="`/images/${feature.image}.jpg`" alt="精选视频封面"/><div class="feature-shade"></div><div class="feature-copy"><span class="feature-tag">{{ feature.tag }}</span><h2>{{ feature.title }}</h2><p>{{ feature.subtitle }}</p></div><div class="feature-bottom"><span>编辑精选 <b>✦</b></span><span class="feature-dots"><i v-for="(_, index) in features" :key="index" :class="{ current: featuredIndex === index }" @click.stop="featuredIndex = index"></i></span></div></button>
        <div class="editor-picks"><div class="section-label"><h2>发现好内容 <span>DISCOVER</span></h2><span>每一帧，都有热爱</span></div><div class="pick-grid"><button v-for="video in videos.slice(1, 3)" :key="video.id" class="pick-card" @click="selectedVideo = video"><div class="pick-image" :style="{ backgroundColor: video.color }"><img :src="`/images/${video.image}.jpg`" :alt="video.title"/><span class="pick-pill">{{ video.category }}精选</span><span class="duration">{{ video.duration }}</span><span class="hover-play">▶</span></div><h3>{{ video.title }}</h3><p><span class="up-tag">UP</span> {{ video.author }}</p></button></div></div>
        <aside class="ranking"><div class="rank-heading"><h2>大家都在看 <span>HOT</span></h2><button @click="notify('完整排行榜即将上线')">排行榜</button></div><button v-for="(video, index) in [videos[3], videos[5], videos[4], videos[2]]" :key="video.id" class="rank-row" @click="selectedVideo = video"><span :class="['rank-number', `rank-${index}`]">{{ index + 1 }}</span><span><strong>{{ video.title }}</strong><small>{{ video.views }}人看过</small></span></button><div class="rank-footer">好内容，值得被更多人发现 <span>✧</span></div></aside>
      </section>

      <section class="recommendations" :aria-busy="feedLoading">
        <div class="feed-heading">
          <div class="feed-title"><span class="tiny-tv">▣</span><h2>{{ search ? '已加载视频的搜索结果' : '最新视频' }}</h2><span class="preview-label">真实 Feed</span></div>
          <button class="refresh-button" :disabled="feedLoading" @click="loadFeed(true)">⟳ 刷新</button>
        </div>
        <div class="filter-bar"><button v-for="category in categories" :key="category" :class="{ chosen: activeCategory === category }" @click="chooseCategory(category)">{{ category === '全部' ? '全部' : category }}</button><span>按发布时间排列 · 分区待开放</span></div>
        <p v-if="search" class="feed-hint">仅搜索已加载的视频；加载更多可扩大结果范围。</p>
        <div v-if="feedError" class="feed-error" role="alert"><span>{{ feedError }}</span><button :disabled="feedLoading" @click="loadFeed(retryReset)">重试</button></div>
        <div v-if="displayedVideos.length" class="video-grid">
          <article v-for="video in displayedVideos" :key="video.id" class="video-card">
            <button class="thumbnail real-thumbnail" :aria-label="'播放：' + video.title" @click="selectedVideo = video"><span class="cover-symbol">▷</span><strong>{{ video.title }}</strong><span class="video-stats"><span>视频 #{{ video.id }}</span><span>点击播放</span></span></button>
            <button class="video-title" @click="selectedVideo = video">{{ video.title }}</button>
            <div class="video-meta"><span class="up-tag">UP</span><span>作者 #{{ video.author_id }}</span><span class="meta-date">· {{ displayDate(video.created_at) }}</span></div>
          </article>
        </div>
        <div v-else-if="feedLoaded && !feedLoading" class="empty-state"><span>▣</span><h3>{{ search ? '已加载的视频中没有匹配结果' : '还没有发布的视频' }}</h3><button v-if="search" @click="query = ''; runSearch()">查看全部</button><p v-else>发布第一个视频后，它会显示在这里。</p></div>
        <div class="load-area">
          <p v-if="feedLoading" role="status">正在加载视频…</p>
          <button v-else-if="feedLoaded && hasMore" @click="loadFeed(false)">加载更多视频 ⌄</button>
          <p v-else-if="feedLoaded && feedVideos.length">已经看完全部视频了</p>
        </div>
      </section>
      <footer class="site-footer"><span class="footer-brand">minibili</span><span>记录生活的闪光，连接同频的热爱。</span><span>独立学习项目 · 非哔哩哔哩官方网站</span><button class="health-button" @click="checkHealth">{{ healthStatus }}</button></footer>
    </main>
    <button class="back-top" aria-label="回到顶部" @click="scrollToTop">⌃<span>顶部</span></button>

    <VideoDialog v-if="selectedVideo?.realVideo" :video-id="selectedVideo.id" :user="currentUser" :api="api" @close="closeOverlay" @login="openLogin" @notify="notify" />
    <PublishDialog v-if="dialog === 'upload'" :api="api" :user="currentUser" @close="closeOverlay" @login="openLogin" @busy="uploadBusy = $event" @published="onPublished" />
    <div v-if="(selectedVideo && !selectedVideo.realVideo) || dialog === 'login'" class="overlay" @click.self="closeOverlay">
      <section class="modal" :class="{ 'video-modal': selectedVideo }" role="dialog" aria-modal="true" :aria-label="selectedVideo ? selectedVideo.title : authMode === 'login' ? '登录' : '注册'">
        <button class="close-modal" aria-label="关闭弹窗" @click="closeOverlay">×</button>
        <template v-if="selectedVideo">
          <div class="preview-cover"><img :src="'/images/' + selectedVideo.image + '.jpg'" :alt="selectedVideo.title"/><div><span class="big-play">▶</span><p>精选布局演示，请到最新视频区观看真实作品</p></div></div>
          <div class="video-modal-body"><span class="modal-category">布局演示</span><h2>{{ selectedVideo.title }}</h2><p>{{ selectedVideo.author }} · 演示内容</p></div>
        </template>
        <template v-else>
          <div class="login-illustration">▣</div><h2 class="modal-title">欢迎来到 minibili</h2>
          <div class="auth-tabs"><button :disabled="authWorking" :class="{ active: authMode === 'login' }" @click="switchAuth('login')">登录</button><button :disabled="authWorking" :class="{ active: authMode === 'register' }" @click="switchAuth('register')">注册</button></div>
          <form class="login-form" @submit.prevent="submitAuth">
            <label>用户名<input v-model="loginName" :disabled="authWorking" placeholder="请输入用户名" autocomplete="username" required/></label>
            <label>密码<input v-model="loginPassword" :disabled="authWorking" type="password" placeholder="请输入密码" :autocomplete="authMode === 'login' ? 'current-password' : 'new-password'" required/></label>
            <label v-if="authMode === 'register'">确认密码<input v-model="confirmPassword" :disabled="authWorking" type="password" placeholder="再次输入密码" autocomplete="new-password" required/></label>
            <p v-if="loginError" class="login-error" role="alert">{{ loginError }}</p>
            <button class="primary-button" type="submit" :disabled="authWorking">{{ authWorking ? '正在提交…' : authMode === 'login' ? '登录' : '创建账号' }}</button>
          </form>
          <p class="form-note">{{ authMode === 'register' ? '用户名3至32位字母、数字或下划线，密码8至72字节。' : '登录状态将在此浏览器保留。' }}</p>
        </template>
      </section>
    </div>
    <Transition name="toast"><div v-if="toast" class="toast-message" role="status">{{ toast }}</div></Transition>
  </div>
</template>
