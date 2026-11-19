import test from 'node:test'
import assert from 'node:assert/strict'
import { createAuthSession, TOKEN_KEY } from './useAuth.js'

function memoryStorage(initial) {
  const values = new Map(initial ? [[TOKEN_KEY, initial]] : [])
  return { getItem: key => values.get(key), setItem: (key, value) => values.set(key, value), removeItem: key => values.delete(key) }
}
test('登录必须等待后端身份确认后才保存凭证', async () => {
  const storage = memoryStorage()
  let confirm
  const session = createAuthSession({storage, api: {login: async () => 'test-token', fetchMe: () => new Promise(resolve => {confirm = resolve})}})
  const pending = session.signIn(' tom ', 'password')
  await Promise.resolve()
  assert.equal(session.user.value, null)
  assert.equal(storage.getItem(TOKEN_KEY), undefined)
  confirm({id: 1})
  assert.equal(await pending, true)
  assert.equal(session.user.value.id, 1)
  assert.equal(storage.getItem(TOKEN_KEY), 'test-token')
  session.logout()
  assert.equal(session.user.value, null)
  assert.equal(storage.getItem(TOKEN_KEY), undefined)
})
test('刷新时失效凭证被清除，网络失败保留凭证供重试', async () => {
  for (const invalid of [true, false]) {
    const storage = memoryStorage('saved-token')
    const session = createAuthSession({storage, api: {fetchMe: async () => {throw invalid ? {response:{status:401}} : new Error('offline')}}})
    await session.restore()
    assert.equal(session.user.value, null)
    assert.equal(session.busy.value, false)
    assert.equal(storage.getItem(TOKEN_KEY), invalid ? undefined : 'saved-token')
    assert.ok(session.notice.value)
  }
})
test('退出后，先前在途的身份请求不能恢复登录', async () => {
  const storage = memoryStorage('saved-token')
  let confirm
  const session = createAuthSession({storage, api: {fetchMe: () => new Promise(resolve => {confirm = resolve})}})
  const pending = session.restore()
  session.logout()
  confirm({id: 1})
  await pending
  assert.equal(session.user.value, null)
  assert.equal(storage.getItem(TOKEN_KEY), undefined)
})
test('密码错误显示提示，身份验证失败也不能保存 Token', async () => {
  const storage = memoryStorage()
  const session = createAuthSession({storage, api: {login: async () => {throw {response:{status:401}}}}})
  assert.equal(await session.signIn('tom', 'wrong'), false)
  assert.equal(session.error.value, '用户名或密码错误')
  assert.equal(storage.getItem(TOKEN_KEY), undefined)
  const failedMe = createAuthSession({storage, api: {login: async () => 'candidate', fetchMe: async () => {throw new Error('offline')}}})
  assert.equal(await failedMe.signIn('tom','password'), false)
  assert.equal(storage.getItem(TOKEN_KEY), undefined)
})
