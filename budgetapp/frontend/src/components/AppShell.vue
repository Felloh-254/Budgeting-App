<template>
  <div class="shell">
    <aside class="sidebar">
      <div class="sidebar-top">
        <div class="sidebar-logo">
          <span>💰</span>
          <span class="logo-text">BudgetFlow</span>
        </div>
        <nav class="nav">
          <RouterLink to="/dashboard" class="nav-link">
            <span class="nav-icon">📊</span> Dashboard
          </RouterLink>
          <RouterLink to="/budgets" class="nav-link">
            <span class="nav-icon">🗂️</span> Budgets
          </RouterLink>
          <RouterLink to="/transactions" class="nav-link">
            <span class="nav-icon">💳</span> Transactions
          </RouterLink>
        </nav>
      </div>
      <div class="sidebar-bottom">
        <div class="user-info">
          <div class="user-avatar">{{ initials }}</div>
          <div>
            <div class="user-name">{{ auth.user?.name }}</div>
            <div class="user-email">{{ auth.user?.email }}</div>
          </div>
        </div>
        <button @click="logout" class="logout-btn" title="Sign out">⏻</button>
      </div>
    </aside>

    <main class="main-content">
      <slot />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()

const initials = computed(() => {
  const name = auth.user?.name || ''
  return name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
})

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.shell { display: flex; height: 100vh; overflow: hidden; }

.sidebar {
  width: 230px; flex-shrink: 0;
  background: var(--surface);
  border-right: 1px solid var(--border);
  display: flex; flex-direction: column;
  justify-content: space-between;
  padding: 24px 16px;
}
.sidebar-logo {
  display: flex; align-items: center; gap: 10px;
  font-size: 16px; font-weight: 700;
  padding: 0 8px; margin-bottom: 32px;
}
.sidebar-logo span:first-child { font-size: 22px; }

.nav { display: flex; flex-direction: column; gap: 4px; }
.nav-link {
  display: flex; align-items: center; gap: 10px;
  padding: 10px 12px; border-radius: 9px;
  font-size: 14px; font-weight: 500; color: var(--muted);
  transition: all .15s;
}
.nav-link:hover { background: var(--surface2); color: var(--text); }
.nav-link.router-link-active { background: rgba(124,106,247,.15); color: var(--accent2); }
.nav-icon { font-size: 16px; width: 20px; text-align: center; }

.sidebar-bottom { display: flex; align-items: center; gap: 10px; }
.user-info { display: flex; align-items: center; gap: 10px; flex: 1; min-width: 0; }
.user-avatar {
  width: 34px; height: 34px; border-radius: 50%;
  background: linear-gradient(135deg, var(--accent), var(--green));
  display: flex; align-items: center; justify-content: center;
  font-size: 12px; font-weight: 700; flex-shrink: 0;
}
.user-name { font-size: 13px; font-weight: 600; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.user-email { font-size: 11px; color: var(--muted); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.logout-btn {
  background: none; border: none; color: var(--muted); font-size: 18px;
  padding: 6px; border-radius: 8px; flex-shrink: 0; transition: all .15s;
}
.logout-btn:hover { color: var(--red); background: rgba(248,113,113,.1); }

.main-content { flex: 1; overflow-y: auto; padding: 32px; }
</style>
