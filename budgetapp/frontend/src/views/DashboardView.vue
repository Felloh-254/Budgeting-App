<template>
  <AppShell>
    <div class="page-header">
      <div>
        <h1>Dashboard</h1>
        <p class="subtitle">{{ greeting }}, {{ auth.user?.name?.split(' ')[0] }} 👋</p>
      </div>
      <span class="date-badge">{{ todayLabel }}</span>
    </div>

    <div v-if="loading" class="loading-grid">
      <div v-for="i in 3" :key="i" class="skeleton-card"></div>
    </div>

    <template v-else>
      <!-- Summary Cards -->
      <div class="summary-grid">
        <div class="stat-card stat-balance">
          <div class="stat-label">Net Balance</div>
          <div class="stat-value" :class="summary.balance >= 0 ? 'positive' : 'negative'">
            {{ fmt(summary.balance) }}
          </div>
          <div class="stat-sub">Total across all time</div>
        </div>
        <div class="stat-card stat-income">
          <div class="stat-label">↑ Total Income</div>
          <div class="stat-value positive">{{ fmt(summary.total_income) }}</div>
        </div>
        <div class="stat-card stat-expense">
          <div class="stat-label">↓ Total Expenses</div>
          <div class="stat-value negative">{{ fmt(summary.total_expenses) }}</div>
        </div>
      </div>

      <!-- Budget Progress -->
      <div class="section-header">
        <h2>Budget Progress</h2>
        <RouterLink to="/budgets" class="see-all">See all →</RouterLink>
      </div>
      <div v-if="summary.budget_stats?.length" class="budgets-grid">
        <div
          v-for="b in summary.budget_stats"
          :key="b.name"
          class="budget-progress-card card"
        >
          <div class="bp-top">
            <div class="bp-dot" :style="{ background: b.color }"></div>
            <span class="bp-name">{{ b.name }}</span>
            <span class="bp-cat">{{ b.category }}</span>
          </div>
          <div class="bp-amounts">
            <span>{{ fmt(b.spent) }} spent</span>
            <span class="muted">of {{ fmt(b.amount) }}</span>
          </div>
          <div class="progress-track">
            <div
              class="progress-fill"
              :style="{
                width: Math.min(100, (b.spent / b.amount) * 100) + '%',
                background: b.color,
                opacity: pct(b) > 90 ? 1 : 0.75
              }"
            ></div>
          </div>
          <div class="bp-pct" :class="{ danger: pct(b) > 90, warn: pct(b) > 70 && pct(b) <= 90 }">
            {{ pct(b).toFixed(0) }}% used
          </div>
        </div>
      </div>
      <div v-else class="empty-state">
        <div class="icon">🗂️</div>
        <strong>No budgets yet</strong>
        <p><RouterLink to="/budgets">Create your first budget</RouterLink></p>
      </div>

      <!-- Monthly Chart -->
      <div v-if="chartData.length" class="section-header" style="margin-top:32px">
        <h2>Monthly Overview</h2>
      </div>
      <div v-if="chartData.length" class="card chart-card">
        <div class="chart-legend">
          <span class="legend-dot income"></span> Income
          <span class="legend-dot expense" style="margin-left:16px"></span> Expenses
        </div>
        <div class="bar-chart">
          <div v-for="m in chartData" :key="m.month" class="bar-group">
            <div class="bars">
              <div
                class="bar income-bar"
                :style="{ height: barHeight(m.income) + 'px' }"
                :title="fmt(m.income)"
              ></div>
              <div
                class="bar expense-bar"
                :style="{ height: barHeight(m.expense) + 'px' }"
                :title="fmt(m.expense)"
              ></div>
            </div>
            <div class="bar-label">{{ shortMonth(m.month) }}</div>
          </div>
        </div>
      </div>
    </template>
  </AppShell>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AppShell from '../components/AppShell.vue'
import { useAuthStore } from '../stores/auth'
import api from '../api/client'

const auth = useAuthStore()
const loading = ref(true)
const summary = ref({ balance: 0, total_income: 0, total_expenses: 0, budget_stats: [], monthly_data: [] })

const todayLabel = new Date().toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric' })
const greeting = new Date().getHours() < 12 ? 'Good morning' : new Date().getHours() < 18 ? 'Good afternoon' : 'Good evening'

const chartData = computed(() => [...(summary.value.monthly_data || [])].reverse())
const maxBar = computed(() => Math.max(...chartData.value.flatMap(m => [m.income, m.expense]), 1))

function fmt(n) {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0 }).format(n || 0)
}
function pct(b) { return b.amount ? (b.spent / b.amount) * 100 : 0 }
function barHeight(val) { return Math.max(4, (val / maxBar.value) * 140) }
function shortMonth(ym) {
  if (!ym) return ''
  const [y, m] = ym.split('-')
  return new Date(y, m - 1).toLocaleString('default', { month: 'short' })
}

onMounted(async () => {
  try {
    const res = await api.get('/summary')
    summary.value = res.data
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.subtitle { color: var(--muted); font-size: 14px; margin-top: 4px; }
.date-badge {
  background: var(--surface2); border: 1px solid var(--border);
  border-radius: 20px; padding: 6px 14px; font-size: 13px; color: var(--muted);
}

.loading-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; margin-bottom: 28px; }
.skeleton-card {
  height: 100px; background: var(--surface); border-radius: var(--radius);
  border: 1px solid var(--border);
  animation: pulse 1.4s ease-in-out infinite;
}
@keyframes pulse { 0%,100%{opacity:.6} 50%{opacity:.3} }

.summary-grid { display: grid; grid-template-columns: 1.4fr 1fr 1fr; gap: 16px; margin-bottom: 32px; }
.stat-card {
  background: var(--surface); border: 1px solid var(--border);
  border-radius: var(--radius); padding: 20px 22px;
  position: relative; overflow: hidden;
}
.stat-card::before {
  content: ''; position: absolute; inset: 0;
  background: linear-gradient(135deg, rgba(124,106,247,.06) 0%, transparent 60%);
  pointer-events: none;
}
.stat-label { font-size: 12px; color: var(--muted); font-weight: 500; text-transform: uppercase; letter-spacing: .5px; margin-bottom: 8px; }
.stat-value { font-size: 26px; font-weight: 700; letter-spacing: -.5px; }
.stat-sub { font-size: 12px; color: var(--muted); margin-top: 6px; }
.positive { color: var(--green); }
.negative { color: var(--red); }
.muted { color: var(--muted); }

.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.section-header h2 { font-size: 16px; font-weight: 600; }
.see-all { font-size: 13px; color: var(--accent2); }
.see-all:hover { text-decoration: underline; }

.budgets-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 16px; }
.budget-progress-card { display: flex; flex-direction: column; gap: 10px; }
.bp-top { display: flex; align-items: center; gap: 8px; }
.bp-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
.bp-name { font-size: 14px; font-weight: 600; flex: 1; }
.bp-cat { font-size: 11px; color: var(--muted); background: var(--surface2); padding: 2px 8px; border-radius: 20px; }
.bp-amounts { display: flex; align-items: center; gap: 6px; font-size: 13px; font-weight: 500; }
.progress-track { height: 6px; background: var(--surface2); border-radius: 3px; overflow: hidden; }
.progress-fill { height: 100%; border-radius: 3px; transition: width .4s ease; }
.bp-pct { font-size: 12px; color: var(--muted); }
.bp-pct.warn { color: var(--yellow); }
.bp-pct.danger { color: var(--red); }

.chart-card { padding: 24px; }
.chart-legend { display: flex; align-items: center; font-size: 13px; color: var(--muted); margin-bottom: 20px; }
.legend-dot { width: 10px; height: 10px; border-radius: 50%; display: inline-block; margin-right: 6px; }
.legend-dot.income { background: var(--green); }
.legend-dot.expense { background: var(--accent); }
.bar-chart { display: flex; align-items: flex-end; gap: 12px; height: 160px; padding-bottom: 24px; position: relative; }
.bar-group { display: flex; flex-direction: column; align-items: center; gap: 6px; flex: 1; }
.bars { display: flex; align-items: flex-end; gap: 4px; }
.bar { width: 18px; border-radius: 4px 4px 0 0; transition: height .3s ease; min-height: 4px; }
.income-bar { background: var(--green); opacity: .8; }
.expense-bar { background: var(--accent); opacity: .8; }
.bar-label { font-size: 11px; color: var(--muted); }
</style>
