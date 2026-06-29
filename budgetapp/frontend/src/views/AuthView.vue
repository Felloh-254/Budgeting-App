<template>
  <div class="auth-wrap">
    <div class="auth-left">
      <div class="auth-brand">BudgetFlow</div>
      <div class="auth-pitch">
        <p class="pitch-eyebrow">Personal finance</p>
        <h1>Know where<br>your money goes.</h1>
        <p class="pitch-sub">Track budgets, log transactions, and see your financial picture at a glance.</p>
      </div>
      <div class="auth-stat-row">
        <div class="auth-stat">
          <span class="stat-val">Ksh 0</span>
          <span class="stat-lbl">hidden fees</span>
        </div>
        <div class="auth-stat-div"></div>
        <div class="auth-stat">
          <span class="stat-val">100%</span>
          <span class="stat-lbl">your data</span>
        </div>
        <div class="auth-stat-div"></div>
        <div class="auth-stat">
          <span class="stat-val">1 file</span>
          <span class="stat-lbl">SQLite, local</span>
        </div>
      </div>
    </div>

    <div class="auth-right">
      <div class="auth-card">
        <h2>{{ isLogin ? 'Sign in' : 'Create account' }}</h2>
        <p class="card-sub">{{ isLogin ? 'Good to have you back.' : 'Get started in seconds.' }}</p>

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
            {{ loading ? 'Just a moment…' : isLogin ? 'Sign in' : 'Create account' }}
          </button>
        </form>

        <p class="switch">
          {{ isLogin ? "Don't have an account?" : 'Already have an account?' }}
          <button @click="isLogin = !isLogin; error = ''">{{ isLogin ? 'Sign up' : 'Sign in' }}</button>
        </p>
      </div>
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
  error.value = ''; loading.value = true
  try {
    if (isLogin.value) await auth.login(form.email, form.password)
    else await auth.register(form.name, form.email, form.password)
    router.push('/dashboard')
  } catch (e) {
    error.value = e.response?.data?.error || 'Something went wrong'
  } finally { loading.value = false }
}
</script>

<style scoped>
.auth-wrap { display: flex; min-height: 100vh; }

.auth-left {
  flex: 1; background: var(--ink); color: var(--paper);
  padding: 48px 56px;
  display: flex; flex-direction: column; justify-content: space-between;
}
.auth-brand { font-size: 15px; font-weight: 600; letter-spacing: -.01em; opacity: .5; }

.auth-pitch { margin: auto 0; }
.pitch-eyebrow {
  font-size: 11px; font-weight: 500; letter-spacing: .1em; text-transform: uppercase;
  color: var(--blue); background: rgba(45,91,227,.2); display: inline-block;
  padding: 4px 10px; border-radius: 4px; margin-bottom: 20px;
}
.auth-pitch h1 { font-size: 40px; font-weight: 600; line-height: 1.15; letter-spacing: -.03em; margin-bottom: 16px; }
.pitch-sub { font-size: 15px; color: rgba(247,246,242,.55); line-height: 1.6; max-width: 340px; }

.auth-stat-row { display: flex; align-items: center; gap: 24px; }
.auth-stat { display: flex; flex-direction: column; gap: 2px; }
.stat-val { font-family: var(--font-num); font-size: 18px; font-weight: 500; }
.stat-lbl { font-size: 11px; color: rgba(247,246,242,.4); text-transform: uppercase; letter-spacing: .06em; }
.auth-stat-div { width: 1px; height: 28px; background: rgba(247,246,242,.12); }

.auth-right {
  width: 420px; flex-shrink: 0; background: var(--paper);
  display: flex; align-items: center; justify-content: center; padding: 48px 40px;
}
.auth-card { width: 100%; }
.auth-card h2 { font-size: 22px; font-weight: 600; letter-spacing: -.02em; margin-bottom: 4px; }
.card-sub { font-size: 13.5px; color: var(--ink-3); margin-bottom: 28px; }

.auth-form { display: flex; flex-direction: column; gap: 14px; }
.btn-full { width: 100%; justify-content: center; padding: 11px; font-size: 14px; margin-top: 6px; }
.auth-error {
  background: var(--red-dim); border: 1px solid #FECACA;
  color: var(--red); border-radius: 6px; padding: 10px 13px; font-size: 13px;
}
.switch { font-size: 13px; color: var(--ink-3); margin-top: 20px; text-align: center; }
.switch button { background: none; border: none; color: var(--blue); font-size: 13px; font-weight: 500; margin-left: 4px; }
.switch button:hover { text-decoration: underline; }

@media (max-width: 768px) {
  .auth-left { display: none; }
  .auth-right { width: 100%; }
}
</style>
