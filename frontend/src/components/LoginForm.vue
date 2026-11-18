<template>
  <form @submit.prevent="login">
    <h2>登录</h2>

    <label>
      用户名
      <input
          v-model="username"
          autocomplete="username"
          required
      />
    </label>

    <label>
      密码
      <input
          v-model="password"
          type="password"
          autocomplete="current-password"
          required
      />
    </label>

    <button type="submit" :disabled="loading">
      {{ loading ? '正在登录……' : '登录' }}
    </button>

    <p v-if="errorMessage" role="alert">{{ errorMessage }}</p>
  </form>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'

const emit = defineEmits(['logged-in'])

const username = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')

async function login() {
  // TODO 1：loading 为 true 时，直接 return，阻止重复提交
      if (loading.value){
        return
      }
  // TODO 2：设置 loading=true，并清空错误提示
      loading.value=true
  try {
    // TODO 3：axios.post 请求 /api/v1/auth/login
    // 第二个参数为对象：
        const response=await axios.post("/api/v1/auth/login",{
          username:username.value.trim(),
          password:password.value
        })
    // username：username.value.trim()
    // password：password.value，不要去空格

    // TODO 4：检查返回的业务码
    // 非 0 时，throw new Error(...)
        if (response.data!==0){
          throw  new  Error(response.data.message||'登陆失败')
        }
    // TODO 5：取出 data.access_token
    // 检查它是非空字符串，否则抛出错误
    const token = response.data.data?.access_token
    if (typeof token !== 'string' || token.trim() === '') {
      throw new Error('登录失败：服务器未返回有效 Token')
    }
    // TODO 6：emit('logged-in', token)
    // 通知父组件登录成功
    // 清空 password.value
    emit('logged-in', token)
    password.value = ''
  } catch (error) {
    // TODO 7：显示后端错误提示，或通用网络错误提示
    // 可以参考 App.vue 的 catch
    errorMessage.value =
        error.response?.data?.message ||
        error.message ||
        '网络错误，请检查网络连接后重试'
  } finally {
    // TODO 8：恢复 loading=false
    loading.value = false
  }
}
</script>