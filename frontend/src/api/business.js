import axios from 'axios'

export function requestMessage(error) {
  if (error.response?.status === 401) return '登录已失效，请重新登录'
  if (typeof error.response?.data?.message === 'string') return error.response.data.message
  if (error.code === 'ECONNABORTED') return '请求超时，请检查网络后重试'
  return error.message || '请求失败，请稍后重试'
}

// 每次请求读取当前凭证，退出或切换账号后不会继续使用旧 Token。
export function createBusinessAPI({ getToken, onUnauthorized, adapter, baseURL = '/api/v1' }) {
  const client = axios.create({ baseURL, timeout: 15000, ...(adapter ? { adapter } : {}) })
  client.interceptors.request.use(config => {
    if (config.protected) {
      const token = getToken()
      if (!token) throw new Error('请先登录')
      config.sessionToken = token
      config.headers.Authorization = 'Bearer ' + token
    }
    return config
  })
  client.interceptors.response.use(response => {
    if (response.data?.code !== 0) throw new Error(response.data?.message || '服务器响应异常')
    return response.data.data
  }, error => {
    if (error.response?.status === 401 && error.config?.sessionToken) onUnauthorized(error.config.sessionToken)
    return Promise.reject(error)
  })
  return {
    register: (username, password) => client.post('/auth/register', { username, password }),
    detail: (id, signal) => client.get(`/videos/${id}`, { signal }),
    count: (id, signal) => client.get(`/videos/${id}/likes/count`, { signal }),
    status: (id, signal) => client.get(`/videos/${id}/like`, { protected: true, signal }),
    setLike: (id, liked) => client.request({ url: `/videos/${id}/like`, method: liked ? 'put' : 'delete', protected: true }),
    upload: (file, onUploadProgress) => {
      const form = new FormData()
      form.append('file', file)
      // 浏览器自动设置含 boundary 的 multipart Content-Type，不手工拼接。
      return client.post('/uploads/videos', form, { protected: true, timeout: 180000, onUploadProgress })
    },
    publish: fields => client.post('/videos', fields, { protected: true }),
  }
}
