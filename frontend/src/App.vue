<template>
  <main class="feed">
    <h1>视频列表</h1>
    <button :disabled="loading" @click="fetchVideos(true)">
      {{ loading ? '正在加载……' : loaded ? '刷新列表' : '加载视频' }}
    </button>

    <p v-if="errorMessage" role="alert">{{ errorMessage }}</p>
    <p v-if="loaded && videos.length === 0">暂时没有视频</p>

    <ul :aria-busy="loading">
      <li v-for="video in videos" :key="video.id">
        <h2>{{ video.title }}</h2>
        <p>{{ video.description || '暂无简介' }}</p>
        <small>视频 #{{ video.id }} · 作者 #{{ video.author_id }}</small>
      </li>
    </ul>

    <button v-if="loaded && hasMore" :disabled="loading" @click="fetchVideos(false)">
      {{ loading ? '正在加载……' : '加载更多' }}
    </button>
    <p v-if="loaded && !hasMore && videos.length > 0">已经到底了</p>
  </main>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'

const videos = ref([])
const loading = ref(false)
const loaded = ref(false)
const errorMessage = ref('')
const hasMore = ref(false)
const nextCursor = ref(0)

// 每批取一条，方便用当前两条视频观察翻页；之后可改为 5 或 10。
const pageSize = 1

async function fetchVideos(reset = true) {
  // 阻止重复请求；没有下一批时不使用 cursor=0 重头查询。
  if (loading.value || (!reset && (!loaded.value || !hasMore.value))) return

  loading.value = true
  errorMessage.value = ''

  try {
    const response = await axios.get('/api/v1/feed', {
      params: {
        limit: pageSize,
        cursor: reset ? 0 : nextCursor.value,
      },
    })
    const result = response.data
    if (result.code !== 0) {
      throw new Error(result.message || '加载失败')
    }

    // 请求成功后再更新列表和游标，失败时可以沿用原游标重试。
    const page = result.data
    videos.value = reset ? page.items : [...videos.value, ...page.items]
    hasMore.value = page.has_more
    nextCursor.value = page.next_cursor
    loaded.value = true
  } catch (error) {
    errorMessage.value = error.response?.data?.message || error.message || '加载失败，请稍后重试'
  } finally {
    // 成功或失败都要恢复按钮。
    loading.value = false
  }
}
</script>

<style scoped>
.feed {
  width: min(640px, 90vw);
  margin: 24px auto;
  text-align: left;
}
ul { list-style: none; padding: 0; }
li { padding: 16px 0; border-bottom: 1px solid #8885; }
h2 { font-size: 1.15rem; }
[role="alert"] { color: #c33; }
button:disabled { cursor: wait; opacity: 0.6; }
</style>