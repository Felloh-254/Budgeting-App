<template>
  <AppShell>
    <div class="page-header">
      <div>
        <p class="eyebrow">{{ todayLabel }}</p>
        <h1>Overview</h1>
      </div>
    </div>

    <div v-if="loading" class="skeleton-grid">
      <div class="skel" style="height:96px"></div>
      <div class="skel" style="height:96px"></div>
      <div class="skel" style="height:96px"></div>
    </div>

    <template v-else>
      <!-- Balance strip -->
      <div class="balance-strip">
        <div class="balance-main">
          <span class="balance-label">Net balance</span>
          <span class="balance-value num-lg" :class="summary.balance >= 0 ? 'pos' : 'neg'">
            {{ fmt(summary.balance) }}
          </span>
        </div>
        <div class="balance-stats">
          <div class="bstat">
            <span class="bstat-arrow up">↑</span>
            <div>
              <div class="bstat-label">Income</div>
              <div class="bstat-val num pos">{{ fmt(summary.total_income) }}</div>
            </div>
          </div>
          <div class="bstat-div"></div>
          <div class="bstat">
            <span class="bstat-arrow dn">↓</span>
            <div>
              <div class="bstat-label">Expenses</div>
              <div class="bstat-val num neg">{{ fmt(summary.total_expenses) }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Monthly chart + budget progress side by side -->
      <div class="two-col">
        <!-- Bar chart -->
        <div class="panel" v-if="chartData.length">
          <div class="panel-head">
            <span class="panel-title">Monthly cash flow</span>
            <div class="legend">
              <span class="leg income"></span> Income
              <span class="leg expense"></span> Expenses
            </div>
          </div>
          <div class="bar-chart">
            <div v-for="m in chartData" :key="m.month" class="bar-col">
              <div class="bar-pair">
                <div class="bar income-bar" :style="{ height: bh(m.income) + 'px' }" :title="fmt(m.income)"></div>
                <div class="bar expense-bar" :style="{ height: bh(m.expense) + 'px' }" :title="fmt(m.expense)"></div>
              </div>
              <div class="bar-lbl">{{ shortMonth(m.month) }}</div>
            </div>
          </div>
        </div>
        <div class="panel" v-else style="align-items:center;justify-content:center;display:flex">
          <div class="empty-state" style="padding:32px">
            <div class="icon">📊</div>
            <strong>No data yet</strong>
            <p>Add transactions to see your chart</p>
          </div>
        </div>

        <!-- Budget table -->
        <div class="panel">
          <div class="panel-head">
            <span class="panel-title">Budget progress</span>
            <RouterLink to="/budgets" class="panel-link">Manage →</RouterLink>
          </div>
          <div v-if="summary.budget_stats?.length" class="budget-rows">
            <div v-for="b in summary.budget_stats" :key="b.name" class="budget-row">
              <div class="budget-row-top">
                <div class="brow-name">
                  <span class="brow-dot" :style="{ background: b.color }"></span>
                  {{ b.name }}
                </div>
                <div class="brow-nums">
                  <span class="num">{{ fmt(b.spent) }}</span>
                  <span class="brow-of">/ {{ fmt(b.amount) }}</span>
                </div>
              </div>
              <div class="prog-track">
                <div class="prog-fill" :style="{
                  width: Math.min(100, (b.spent/b.amount)*100) + '%',
                  background: pct(b) > 90 ? 'var(--red)' : pct(b) > 70 ? '#D97706' : b.color
                }"></div>
              </div>
            </div>
          </div>
          <div v-else class="empty-state" style="padding:32px">
            <div class="icon">🗂️</div>
            <strong>No budgets</strong>
            <p><RouterLink to="/budgets" style="color:var(--blue)">Create one</RouterLink> to track spending</p>
          </div>
        </div>
      </div>

      <!-- Recent transactions -->
      <div class="panel" style="margin-top:24px" v-if="recentTxns.length">
        <div class="panel-head">
          <span class="panel-title">Recent transactions</span>
          <RouterLink to="/transactions" class="panel-link">All transactions →</RouterLink>
        </div>
        <table class="txn-table">
          <thead>
            <tr>
              <th>Description</th>
              <th>Category</th>
              <th>Date</th>
              <th class="right">Amount</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="t in recentTxns" :key="t.id">
              <td class="txn-title-cell">
                <span class="txn-dot" :class="t.type"></span>
                {{ t.title }}
              </td>
              <td class="cat-cell">{{ t.category }}</td>
              <td class="date-cell num">{{ fmtDate(t.date) }}</td>
              <td class="amount-cell num right" :class="t.type">
                {{ t.type === 'income' ? '+' : '−' }}{{ fmt(t.amount) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </AppShell>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AppShell from '../components/AppShell.vue'
import { useTransactionStore } from '../stores/transactions'
import api from '../api/client'

const txnStore = useTransactionStore()
const loading = ref(true)
const summary = ref({ balance: 0, total_income: 0, total_expenses: 0, budget_stats: [], monthly_data: [] })

const todayLabel = new Date().toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric', year: 'numeric' })
const chartData = computed(() => [...(summary.value.monthly_data || [])].reverse())
const maxBar = computed(() => Math.max(...chartData.value.flatMap(m => [m.income, m.expense]), 1))
const recentTxns = computed(() => (txnStore.transactions || []).slice(0, 8))

function bh(v) { return Math.max(3, (v / maxBar.value) * 120) }
function pct(b) { return b.amount ? (b.spent / b.amount) * 100 : 0 }
function fmt(n) { return new Intl.NumberFormat('en-KE', { style: 'currency', currency: 'KES', currencyDisplay: 'narrowSymbol', maximumFractionDigits: 0 }).format(n || 0) }
function fmtDate(d) { return new Date(d + 'T12:00').toLocaleDateString('en-US', { month: 'short', day: 'numeric' }) }
function shortMonth(ym) { if (!ym) return ''; const [y,m] = ym.split('-'); return new Date(y,m-1).toLocaleString('default',{month:'short'}) }

onMounted(async () => {
  try {
    const [sum] = await Promise.all([api.get('/summary'), txnStore.fetchTransactions()])
    summary.value = sum.data
  } finally { loading.value = false }
})
</script>

<style scoped>
.eyebrow { font-size: 12px; color: var(--ink-3); margin-bottom: 4px; text-transform: uppercase; letter-spacing: .06em; font-weight: 500; }

.skeleton-grid { display: grid; grid-template-columns: repeat(3,1fr); gap: 16px; }
.skel { background: var(--line); border-radius: 8px; animation: pulse 1.5s ease-in-out infinite; }
@keyframes pulse { 0%,100%{opacity:.5} 50%{opacity:.25} }

/* Balance strip */
.balance-strip {
  display: flex; align-items: center; justify-content: space-between;
  background: var(--ink); color: var(--paper);
  border-radius: 10px; padding: 28px 36px; margin-bottom: 24px;
}
.balance-label { display: block; font-size: 11px; text-transform: uppercase; letter-spacing: .08em; opacity: .45; margin-bottom: 6px; }
.balance-value { font-size: 42px; font-weight: 400; letter-spacing: -.03em; line-height: 1; }
.balance-value.pos { color: #6EE7B7; }
.balance-value.neg { color: #FCA5A5; }

.balance-stats { display: flex; align-items: center; gap: 28px; }
.bstat { display: flex; align-items: center; gap: 10px; }
.bstat-arrow { font-size: 18px; width: 32px; height: 32px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-weight: 700; }
.bstat-arrow.up { background: rgba(110,231,183,.15); color: #6EE7B7; }
.bstat-arrow.dn { background: rgba(252,165,165,.15); color: #FCA5A5; }
.bstat-label { font-size: 11px; opacity: .45; text-transform: uppercase; letter-spacing: .06em; margin-bottom: 2px; }
.bstat-val { font-size: 17px; font-weight: 500; letter-spacing: -.01em; }
.bstat-val.pos { color: #6EE7B7; }
.bstat-val.neg { color: #FCA5A5; }
.bstat-div { width: 1px; height: 36px; background: rgba(247,246,242,.1); }

/* Two-col grid */
.two-col { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }

/* Panel */
.panel { background: white; border: 1px solid var(--line); border-radius: 10px; padding: 22px 24px; }
.panel-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.panel-title { font-size: 13px; font-weight: 600; text-transform: uppercase; letter-spacing: .06em; color: var(--ink-3); }
.panel-link { font-size: 13px; color: var(--blue); font-weight: 500; }
.panel-link:hover { text-decoration: underline; }

/* Bar chart */
.legend { display: flex; align-items: center; gap: 14px; font-size: 12px; color: var(--ink-3); }
.leg { display: inline-block; width: 8px; height: 8px; border-radius: 2px; margin-right: 4px; }
.leg.income { background: var(--green); }
.leg.expense { background: var(--blue); }

.bar-chart { display: flex; align-items: flex-end; gap: 8px; height: 140px; padding-bottom: 24px; position: relative; }
.bar-col { display: flex; flex-direction: column; align-items: center; gap: 6px; flex: 1; }
.bar-pair { display: flex; gap: 3px; align-items: flex-end; }
.bar { width: 14px; border-radius: 3px 3px 0 0; min-height: 3px; transition: height .4s cubic-bezier(.4,0,.2,1); }
.income-bar { background: var(--green); opacity: .75; }
.expense-bar { background: var(--blue); opacity: .65; }
.bar-lbl { font-size: 10.5px; color: var(--ink-3); font-family: var(--font-num); }

/* Budget rows */
.budget-rows { display: flex; flex-direction: column; gap: 14px; }
.budget-row { display: flex; flex-direction: column; gap: 7px; }
.budget-row-top { display: flex; justify-content: space-between; align-items: center; }
.brow-name { display: flex; align-items: center; gap: 8px; font-size: 13.5px; font-weight: 500; }
.brow-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.brow-nums { display: flex; align-items: baseline; gap: 4px; font-size: 13px; }
.brow-of { font-size: 12px; color: var(--ink-3); }
.prog-track { height: 4px; background: var(--line); border-radius: 2px; overflow: hidden; }
.prog-fill { height: 100%; border-radius: 2px; transition: width .5s cubic-bezier(.4,0,.2,1); }

/* Transactions table */
.txn-table { width: 100%; border-collapse: collapse; }
.txn-table thead tr { border-bottom: 1px solid var(--line); }
.txn-table th {
  font-size: 11px; font-weight: 500; text-transform: uppercase; letter-spacing: .06em;
  color: var(--ink-3); padding: 0 0 10px; text-align: left;
}
.txn-table th.right { text-align: right; }
.txn-table tbody tr { border-bottom: 1px solid var(--line); transition: background .1s; }
.txn-table tbody tr:last-child { border-bottom: none; }
.txn-table tbody tr:hover { background: var(--paper); }
.txn-table td { padding: 11px 0; font-size: 13.5px; }

.txn-title-cell { display: flex; align-items: center; gap: 10px; font-weight: 500; }
.txn-dot { width: 6px; height: 6px; border-radius: 50%; flex-shrink: 0; }
.txn-dot.income { background: var(--green); }
.txn-dot.expense { background: var(--ink-3); }
.cat-cell { color: var(--ink-3); font-size: 13px; }
.date-cell { color: var(--ink-3); font-size: 12.5px; }
.amount-cell.income { color: var(--green); }
.amount-cell.expense { color: var(--ink); }
.right { text-align: right; }

/* Utility */
.pos { color: var(--green); }
.neg { color: var(--red); }
.num { font-family: var(--font-num); }
.num-lg { font-family: var(--font-num); }
</style>
