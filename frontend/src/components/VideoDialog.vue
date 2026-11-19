<script setup>
import { onBeforeUnmount, ref, watch } from 'vue'
import { playbackURL } from '../api/feed'
import { requestMessage } from '../api/business'

const props = defineProps({ videoId: Number, user: Object, api: Object })
const emit = defineEmits(['close', 'login', 'notify'])
const detail = ref(null)
const detailError = ref('')
const count = ref(null)
const liked = ref(null)
const interactionError = ref('')
const loading = ref(false)
const saving = ref(false)
const playbackError = ref(false)
let controller
let revision = 0
async function load() {
  const attempt = ++revision
  controller?.abort()
  controller = new AbortController()
  const signal = controller.signal
  detail.value = null
  detailError.value = ''
  interactionError.value = ''
  playbackError.value = false
  count.value = null
  liked.value = null
  loading.value = true
  const requests = [props.api.detail(props.videoId, signal), props.api.count(props.videoId, signal)]
  if (props.user) requests.push(props.api.status(props.videoId, signal))
  const result = await Promise.allSettled(requests)
  if (attempt !== revision) return
  loading.value = false
  if (result[0].status === 'fulfilled') detail.value = result[0].value
  else detailError.value = requestMessage(result[0].reason)
  if (result[1].status === 'fulfilled') count.value = result[1].value.count
  else interactionError.value = '点赞数加载失败，请重试'
  if (result[2]?.status === 'fulfilled') liked.value = result[2].value.liked
  else if (result[2]) interactionError.value = requestMessage(result[2].reason)
}
async function changeLike() {
  if (!props.user) { emit('login'); return }
  if (saving.value || liked.value === null) return
  const id = props.videoId
  const account = props.user.id
  const desired = !liked.value
  saving.value = true
  interactionError.value = ''
  try {
    await props.api.setLike(id, desired)
    if (id !== props.videoId || account !== props.user?.id) return
    liked.value = desired
    emit('notify', desired ? '点赞成功' : '取消点赞成功')
    // 重新读取数据库数量，避免只在前端加减造成数据漂移。
    const result = await props.api.count(id)
    if (id === props.videoId && account === props.user?.id) count.value = result.count
  } catch (error) {
    if (id === props.videoId && account === props.user?.id) {
      interactionError.value = requestMessage(error) + '；请重新加载状态后确认结果'
      liked.value = null
    }
  } finally { saving.value = false }
}
watch(() => [props.videoId, props.user?.id], load, { immediate: true })
onBeforeUnmount(() => { revision++; controller?.abort() })
</script>

<template>
  <div class="overlay" @click.self="emit('close')">
    <section class="modal video-modal" role="dialog" aria-modal="true" aria-label="视频详情">
      <button class="close-modal" aria-label="关闭" @click="emit('close')">×</button>
      <div v-if="loading" class="dialog-state" role="status">正在加载视频详情…</div>
      <div v-else-if="detailError" class="dialog-state" role="alert"><p>{{ detailError }}</p><button class="primary-button" @click="load">重新加载</button></div>
      <template v-else-if="detail">
        <div class="real-player"><video :key="detail.id" :src="playbackURL(detail.video_url)" controls playsinline preload="metadata" @error="playbackError = true"></video><p v-if="playbackError" role="alert">视频无法播放，请检查文件是否存在及编码是否兼容。</p></div>
        <div class="video-modal-body">
          <span class="modal-category">已发布视频</span><h2>{{ detail.title }}</h2>
          <p>作者 #{{ detail.author_id }} · {{ new Date(detail.created_at).toLocaleString('zh-CN') }}</p>
          <p class="real-description">{{ detail.description || '作者还没有填写简介' }}</p>
          <div class="modal-actions"><button :class="{ liked: liked === true }" :disabled="saving || (user && liked === null)" @click="changeLike">{{ saving ? '正在保存…' : liked ? '♥ 已点赞' : '♡ 点赞' }}</button><span>{{ count === null ? '点赞数暂不可用' : count + ' 人点赞' }}</span><span v-if="!user">登录后参与互动</span></div>
          <div v-if="interactionError" class="feed-error" role="alert"><span>{{ interactionError }}</span><button :disabled="saving" @click="load">重新加载</button></div>
        </div>
      </template>
    </section>
  </div>
</template>
