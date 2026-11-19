import test from 'node:test'
import assert from 'node:assert/strict'
import axios from 'axios'
import { createBusinessAPI } from './business.js'

test('公开接口不发送Token，保护接口读取最新Token并使用正确HTTP方法', async () => {
  const calls = []
  let token = 'first'
  const api = createBusinessAPI({ getToken: () => token, onUnauthorized: () => {}, adapter: async config => {
    calls.push(config)
    return {data:{code:0,data:{liked:true,count:1}},status:200,statusText:'OK',headers:{},config}
  }})
  await api.detail(3)
  await api.setLike(3,true)
  token = 'second'
  await api.setLike(3,false)
  assert.equal(calls[0].headers.Authorization, undefined)
  assert.equal(calls[1].method,'put')
  assert.equal(calls[1].headers.Authorization,'Bearer first')
  assert.equal(calls[2].method,'delete')
  assert.equal(calls[2].headers.Authorization,'Bearer second')
})
test('无凭证不发送保护请求，401报告该次请求使用的旧凭证', async () => {
  let token = ''
  let expired
  let calls = 0
  const api = createBusinessAPI({getToken:() => token,onUnauthorized: value => {expired = value},adapter:async config => {
    calls++
    token = 'new-session'
    throw new axios.AxiosError('Unauthorized','ERR_BAD_REQUEST',config,null,{status:401,data:{code:40102},config})
  }})
  await assert.rejects(api.status(3),/请先登录/)
  assert.equal(calls,0)
  token = 'old-session'
  await assert.rejects(api.status(3))
  assert.equal(expired,'old-session')
  assert.equal(token,'new-session')
})
