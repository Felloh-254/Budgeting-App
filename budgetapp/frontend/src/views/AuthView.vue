<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-logo">
        <span class="logo-icon"></span>
        <h1>BudgetFlow</h1>
        <p>{{ isLogin ? 'Welcome back' : 'Create your account' }}</p>
      </div>

      <form @submit.prevent="submit" class="auth-form">
        <div v-if="!isLogin" class="form-group">
          <label>Full name</label>
          <input v-model="form.name" class="form-control" placeholder="Jane Smith" required />
        </div>
        <div class="form-group">
          <label>Email</label>
          <input v-model="form.email" type="email" class="form-control" placeholder="you@example.com" required />
        </div>
        <div class="form-group">
          <label>Password</label>
          <input v-model="form.password" type="password" class="form-control" placeholder="••••••••" required />
        </div>

        <div v-if="error" class="auth-error">{{ error }}</div>

        <button type="submit" class="btn btn-primary btn-full" :disabled="loading">
          <span v-if="loading"></span>
          {{ isLogin ? 'Sign in' : 'Create account' }}
        </button>
      </form>

      <p class="auth-toggle">
        {{ isLogin ? "Don't have an account?" : 'Already have an account?' }}
        <button @click="isLogin = !isLogin; error = ''">
          {{ isLogin ? 'Sign up' : 'Sign in' }}
        </button>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()

const isLogin = ref(true)
const loading = ref(false)
const error = ref('')
const form = reactive({ name: '', email: '', password: '' })

async function submit() {
  error.value = ''
  loading.value = true
  try {
    if (isLogin.value) {
      await auth.login(form.email, form.password)
    } else {
      await auth.register(form.name, form.email, form.password)
    }
    router.push('/dashboard')
  } catch (e) {
    error.value = e.response?.data?.error || 'Something went wrong'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: radial-gradient(ellipse at 60% 0%, rgba(124,106,247,.18) 0%, transparent 60%),
              radial-gradient(ellipse at 20% 100%, rgba(34,211,160,.1) 0%, transparent 50%),
              var(--bg);
  padding: 24px;
}
.auth-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 20px;
  padding: 40px;
  width: 100%;
  max-width: 420px;
}
.auth-logo { text-align: center; margin-bottom: 32px; }
.logo-icon { font-size: 40px; display: block; margin-bottom: 10px; }
.auth-logo h1 { font-size: 24px; font-weight: 700; letter-spacing: -.5px; }
.auth-logo p { color: var(--muted); font-size: 14px; margin-top: 4px; }
.auth-form { display: flex; flex-direction: column; gap: 16px; }
.btn-full { width: 100%; justify-content: center; padding: 12px; font-size: 15px; margin-top: 4px; }
.auth-error {
  background: rgba(248,113,113,.12); border: 1px solid rgba(248,113,113,.3);
  color: var(--red); border-radius: 8px; padding: 10px 14px; font-size: 13px;
}
.auth-toggle { text-align: center; font-size: 13px; color: var(--muted); margin-top: 20px; }
.auth-toggle button {
  background: none; border: none; color: var(--accent2); font-size: 13px;
  cursor: pointer; margin-left: 4px; font-weight: 500;
}
.auth-toggle button:hover { text-decoration: underline; }
</style>
