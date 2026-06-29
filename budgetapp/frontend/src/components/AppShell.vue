<template>
  <div class="shell">
    <header class="topbar">
      <div class="topbar-left">
        <div class="wordmark">Budget<span>Flow</span></div>
        <nav class="nav">
          <RouterLink to="/dashboard" class="nav-link">Overview</RouterLink>
          <RouterLink to="/budgets" class="nav-link">Budgets</RouterLink>
          <RouterLink to="/transactions" class="nav-link">Transactions</RouterLink>
        </nav>
      </div>
      <div class="topbar-right">
        <span class="user-name">{{ auth.user?.name }}</span>
        <button class="sign-out" @click="logout">Sign out</button>
      </div>
    </header>

    <main class="content">
      <slot />
    </main>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
const auth = useAuthStore()
const router = useRouter()
function logout() { auth.logout(); router.push('/login') }
</script>

<style scoped>
.shell { display: flex; flex-direction: column; min-height: 100vh; background: var(--paper); }

.topbar {
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 40px; height: 56px;
  background: white;
  border-bottom: 1px solid var(--line);
  position: sticky; top: 0; z-index: 50;
}

.topbar-left { display: flex; align-items: center; gap: 36px; }

.wordmark {
  font-size: 15px; font-weight: 600; letter-spacing: -.02em; color: var(--ink);
  flex-shrink: 0;
}
.wordmark span { color: var(--blue); }

.nav { display: flex; gap: 2px; }
.nav-link {
  font-size: 13.5px; font-weight: 500; color: var(--ink-3);
  padding: 6px 12px; border-radius: 5px; transition: all .13s;
}
.nav-link:hover { color: var(--ink); background: var(--paper); }
.nav-link.router-link-active { color: var(--ink); background: var(--paper); }

.topbar-right { display: flex; align-items: center; gap: 16px; }
.user-name { font-size: 13px; color: var(--ink-3); }
.sign-out {
  font-size: 13px; color: var(--ink-3); background: none; border: none;
  padding: 0; transition: color .13s;
}
.sign-out:hover { color: var(--red); }

.content { flex: 1; padding: 40px; max-width: 1140px; width: 100%; margin: 0 auto; }
</style>
