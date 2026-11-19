<script setup>
import { computed, ref, watch } from 'vue'
import { requestMessage } from '../api/business'

const props = defineProps({ api: Object, user: Object })
const emit = defineEmits(['close', 'login', 'published', 'busy'])
const mode = ref('file')
const file = ref(null)
const title = ref('')
const description = ref('')
const externalURL = ref('')
const filename = ref('')
const error = ref('')
const busy = ref(false)
const progress = ref(0)
const publishing = ref(false)
const publishedID = ref(0)
const uploadedURL = computed(() => filename.value ? new URL('/uploads/videos/' + encodeURIComponent(filename.value), window.location.origin).href : '')
watch(busy, value => emit('busy', value), { flush: 'sync' })
function selectFile(event) {
  file.value = event.target.files?.[0] || null
  filename.value = ''
  progress.value = 0
  error.value = ''
}
async function upload() {
  if (busy.value || publishedID.value) return
  if (!props.user) { emit('login'); return }
  error.value = ''
  if (!file.value || !/\.mp4$/i.test(file.value.name)) { error.value = '请选择 MP4 文件'; return }
  if (file.value.size <= 0 || file.value.size > 100 * 1024 * 1024) { error.value = '视频大小应大于 0 且不超过 100MiB'; return }
  busy.value = true
  progress.value = 0
  try {
    const result = await props.api.upload(file.value, event => {
      if (event.total) progress.value = Math.min(100, Math.round(event.loaded / event.total * 100))
    })
    if (typeof result.filename !== 'string' || !/^[a-f0-9]+\.mp4$/i.test(result.filename)) throw new Error('上传响应异常')
    filename.value = result.filename
  } catch (failure) { error.value = requestMessage(failure) }
  finally { busy.value = false }
}
async function publish() {
  if (busy.value || publishedID.value) return
  if (!props.user) { emit('login'); return }
  error.value = ''
  if ([...title.value.trim()].length < 1 || [...title.value.trim()].length > 100) { error.value = '标题必须为 1 至 100 个字符'; return }
  if ([...description.value.trim()].length > 1000) { error.value = '简介不能超过 1000 个字符'; return }
  const videoURL = mode.value === 'file' ? uploadedURL.value : externalURL.value.trim()
  if (!videoURL) { error.value = mode.value === 'file' ? '请先上传文件，再发布视频' : '请填写视频地址'; return }
  try {
    const url = new URL(videoURL)
    if (!['http:', 'https:'].includes(url.protocol) || new TextEncoder().encode(videoURL).length > 2048) throw new Error()
  } catch { error.value = '请填写有效的 HTTP 或 HTTPS 视频地址'; return }
  busy.value = true
  publishing.value = true
  try {
    const result = await props.api.publish({ title: title.value.trim(), description: description.value.trim(), video_url: videoURL })
    publishedID.value = result.id
  } catch (failure) { error.value = requestMessage(failure) }
  finally { busy.value = false; publishing.value = false }
  if (publishedID.value) emit('published', publishedID.value)
}
</script>

<template>
  <div class="overlay" @click.self="!busy && emit('close')">
    <section class="modal publish-modal" role="dialog" aria-modal="true" aria-label="发布视频">
      <button class="close-modal" :disabled="busy" aria-label="关闭" @click="emit('close')">×</button>
      <h2 class="modal-title">把热爱分享给世界</h2><p class="modal-description">上传作品，填写信息，再发布到首页</p>
      <template v-if="!user"><p class="form-note">请先登录后投稿。</p><button class="primary-button" @click="emit('login')">前往登录</button></template>
      <form v-else class="publish-form" @submit.prevent="publish">
        <fieldset :disabled="busy || !!publishedID">
          <div class="auth-tabs"><button type="button" :class="{ active: mode === 'file' }" @click="mode = 'file'">上传文件</button><button type="button" :class="{ active: mode === 'url' }" @click="mode = 'url'">使用视频链接</button></div>
          <div v-if="mode === 'file'" class="upload-zone"><label>视频文件（MP4，最多100MiB）<input type="file" accept=".mp4,video/mp4" @change="selectFile"/></label><p v-if="file">{{ file.name }} · {{ (file.size / 1024 / 1024).toFixed(2) }} MiB</p><button class="secondary-button" type="button" :disabled="!file || !!filename" @click="upload">{{ filename ? '已上传' : '上传文件' }}</button><p v-if="busy && !publishing" role="status">传输进度 {{ progress }}%{{ progress === 100 ? '，服务器正在转换为可播放格式…' : '' }}</p><p v-if="filename" class="upload-success">上传成功，可以填写信息并发布。</p></div>
          <label v-else>视频地址<input v-model="externalURL" type="url" placeholder="https://…/video.mp4" required/></label>
          <label>标题<input v-model="title" placeholder="给视频起个标题" required/><small>{{ [...title].length }}/100</small></label>
          <label>简介<textarea v-model="description" rows="3" placeholder="分享视频背后的故事"></textarea><small>{{ [...description].length }}/1000</small></label>
        </fieldset>
        <p v-if="error" class="login-error" role="alert">{{ error }}</p>
        <button class="primary-button" type="submit" :disabled="busy || !!publishedID">{{ publishing ? '正在发布…' : busy ? '正在上传…' : '发布视频' }}</button>
        <p class="form-note">上传成功后仍需发布，作品才会进入首页。</p>
      </form>
    </section>
  </div>
</template>
