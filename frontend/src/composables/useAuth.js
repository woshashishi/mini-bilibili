import { ref } from 'vue'
import * as authAPI from '../api/auth.js'

export const TOKEN_KEY = 'minibili.access_token'
export function createAuthSession({ api = authAPI, storage } = {}) {
  const user = ref(null)
  const busy = ref(false)
  const error = ref('')
  const notice = ref('')
  let token = ''
  let revision = 0
  function persist(value) {
    try {
      if (!storage) throw new Error('Storage unavailable')
      if (value) storage.setItem(TOKEN_KEY, value)
      else storage.removeItem(TOKEN_KEY)
    } catch { if (value) notice.value = '浏览器无法保存登录状态，关闭页面后需重新登录' }
  }
  function logout() {
    revision++
    token = ''
    user.value = null
    busy.value = false
    error.value = ''
    notice.value = ''
    persist('')
  }
  async function signIn(username, password) {
    if (busy.value) return false
    const attempt = ++revision
    busy.value = true
    error.value = ''
    notice.value = ''
    try {
      const candidate = await api.login(username.trim(), password)
      // 身份由后端确认后，才能显示已登录并保存凭证。
      const identity = await api.fetchMe(candidate)
      if (attempt !== revision) return false
      token = candidate
      user.value = identity
      persist(token)
      return true
    } catch (failure) {
      if (attempt === revision) error.value = authAPI.authErrorMessage(failure)
      return false
    } finally { if (attempt === revision) busy.value = false }
  }
  async function restore() {
    if (busy.value) return
    try { token = storage?.getItem(TOKEN_KEY) || '' } catch { token = '' }
    if (!token) return
    const attempt = ++revision
    busy.value = true
    notice.value = ''
    try {
      const identity = await api.fetchMe(token)
      if (attempt === revision) user.value = identity
    } catch (failure) {
      if (attempt !== revision) return
      if (failure.response?.status === 401) {
        logout()
        notice.value = '登录已失效，请重新登录'
      } else {
        notice.value = '暂时无法确认登录状态，请重试'
      }
    } finally { if (attempt === revision) busy.value = false }
  }
  function getToken() { return token }
  function expire(expectedToken) {
    if (!expectedToken || expectedToken !== token) return
    logout()
    notice.value = '登录已失效，请重新登录'
  }
  return { user, busy, error, notice, signIn, restore, logout, getToken, expire }
}
export function useAuth() {
  let storage
  try { storage = window.localStorage } catch { /* 仍可使用当前页面会话 */ }
  return createAuthSession({ storage })
}
