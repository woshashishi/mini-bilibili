import axios from 'axios'

const client = axios.create({ baseURL: '/api/v1', timeout: 10000 })
export async function login(username, password) {
  const { data } = await client.post('/auth/login', { username, password })
  if (data.code !== 0 || typeof data.data?.access_token !== 'string' || !data.data.access_token || data.data.token_type !== 'Bearer') throw new Error('登录响应格式错误')
  return data.data.access_token
}
export async function fetchMe(token) {
  const { data } = await client.get('/me', { headers: { Authorization: 'Bearer ' + token } })
  if (data.code !== 0 || !Number.isSafeInteger(data.data?.id) || data.data.id <= 0) throw new Error('身份响应格式错误')
  return data.data
}
export function authErrorMessage(error) {
  const status = error.response?.status
  if (status === 401) return '用户名或密码错误'
  if (status === 400) return '请填写用户名和密码'
  if (!error.response && error.isAxiosError) return '无法连接服务器，请稍后重试'
  return '登录暂时失败，请稍后重试'
}
