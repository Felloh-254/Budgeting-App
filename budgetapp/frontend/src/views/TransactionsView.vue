<template>
  <AppShell>
    <div class="page-header">
      <div>
        <p class="eyebrow">All activity</p>
        <h1>Transactions</h1>
      </div>
      <button class="btn btn-primary" @click="openCreate">+ Add transaction</button>
    </div>

    <!-- Filter bar -->
    <div class="filter-bar">
      <div class="search-wrap">
        <svg class="search-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        <input v-model="search" class="search-input" placeholder="Search transactions…" />
      </div>
      <select v-model="filterType" class="form-control filt">
        <option value="">All types</option>
        <option value="income">Income</option>
        <option value="expense">Expense</option>
      </select>
      <select v-model="filterCategory" class="form-control filt">
        <option value="">All categories</option>
        <option v-for="c in usedCats" :key="c">{{ c }}</option>
      </select>
    </div>

    <!-- Summary row -->
    <div class="sum-row">
      <div class="sum-item">
        <span class="sum-lbl">Income</span>
        <span class="sum-val num income">{{ fmt(filteredIncome) }}</span>
      </div>
      <div class="sum-div"></div>
      <div class="sum-item">
        <span class="sum-lbl">Expenses</span>
        <span class="sum-val num expense">{{ fmt(filteredExpenses) }}</span>
      </div>
      <div class="sum-div"></div>
      <div class="sum-item">
        <span class="sum-lbl">Net</span>
        <span class="sum-val num" :class="filteredNet >= 0 ? 'income' : 'expense'">{{ fmt(filteredNet) }}</span>
      </div>
      <span class="sum-count">{{ filtered.length }} transaction{{ filtered.length !== 1 ? 's' : '' }}</span>
    </div>

    <!-- Loading -->
    <div v-if="store.loading" class="skel-list">
      <div v-for="i in 6" :key="i" class="skel" style="height:52px"></div>
    </div>

    <!-- Empty -->
    <div v-else-if="!filtered.length" class="empty-state" style="margin-top:48px">
      <div class="icon">💳</div>
      <strong>{{ store.transactions.length ? 'No matches' : 'No transactions yet' }}</strong>
      <p>{{ store.transactions.length ? 'Try adjusting your filters' : 'Add your first transaction above' }}</p>
    </div>

    <!-- Table -->
    <div v-else class="txn-panel">
      <table class="txn-table">
        <thead>
          <tr>
            <th>Description</th>
            <th>Category</th>
            <th>Budget</th>
            <th class="right">Date</th>
            <th class="right">Amount</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in filtered" :key="t.id" class="txn-row">
            <td class="desc-cell">
              <span class="type-pip" :class="t.type"></span>
              <span class="desc-text">{{ t.title }}</span>
              <span v-if="t.note" class="note-text">{{ t.note }}</span>
            </td>
            <td class="cat-cell">{{ t.category }}</td>
            <td class="bud-cell">{{ budgetName(t.budget_id) }}</td>
            <td class="date-cell num right">{{ fmtDate(t.date) }}</td>
            <td class="amount-cell num right" :class="t.type">
              {{ t.type === 'income' ? '+' : '−' }}{{ fmt(t.amount) }}
            </td>
            <td class="actions-cell">
              <div class="row-actions">
                <button class="icon-btn" @click="openEdit(t)">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                </button>
                <button class="icon-btn danger" @click="confirmDelete(t)">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/><path d="M10 11v6M14 11v6"/></svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create / Edit Modal -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ editing ? 'Edit transaction' : 'New transaction' }}</h3>
          <button class="icon-btn" @click="showModal = false">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </div>
        <div class="form-stack">
          <div class="type-toggle">
            <button class="toggle-btn" :class="{ active: form.type === 'expense' }" @click="form.type = 'expense'">Expense</button>
            <button class="toggle-btn income-active" :class="{ active: form.type === 'income' }" @click="form.type = 'income'">Income</button>
          </div>
          <div class="form-group">
            <label>Description</label>
            <input v-model="form.title" class="form-control" placeholder="e.g. Grocery run" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Amount (Ksh)</label>
              <input v-model.number="form.amount" type="number" min="0" step="0.01" class="form-control" placeholder="0.00" />
            </div>
            <div class="form-group">
              <label>Date</label>
              <input v-model="form.date" type="date" class="form-control" />
            </div>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Category</label>
              <select v-model="form.category" class="form-control">
                <option v-for="c in categories" :key="c">{{ c }}</option>
              </select>
            </div>
            <div class="form-group">
              <label>Budget (optional)</label>
              <select v-model="form.budget_id" class="form-control">
                <option :value="null">None</option>
                <option v-for="b in budgetStore.budgets" :key="b.id" :value="b.id">{{ b.name }}</option>
              </select>
            </div>
          </div>
          <div class="form-group">
            <label>Note (optional)</label>
            <input v-model="form.note" class="form-control" placeholder="Add a note…" />
          </div>
        </div>
        <div v-if="formError" class="form-err">{{ formError }}</div>
        <div class="modal-footer">
          <button class="btn btn-ghost btn-sm" @click="showModal = false">Cancel</button>
          <button class="btn btn-primary btn-sm" @click="save" :disabled="saving">{{ saving ? 'Saving…' : editing ? 'Update' : 'Add' }}</button>
        </div>
      </div>
    </div>

    <!-- Delete confirm -->
    <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
      <div class="modal" style="max-width:380px">
        <div class="modal-header">
          <h3>Delete transaction</h3>
          <button class="icon-btn" @click="deleteTarget = null">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </div>
        <p class="confirm-text">Delete <strong>{{ deleteTarget.title }}</strong>?</p>
        <div class="modal-footer">
          <button class="btn btn-ghost btn-sm" @click="deleteTarget = null">Cancel</button>
          <button class="btn btn-danger btn-sm" @click="doDelete" :disabled="saving">Delete</button>
        </div>
      </div>
    </div>
  </AppShell>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import AppShell from '../components/AppShell.vue'
import { useTransactionStore } from '../stores/transactions'
import { useBudgetStore } from '../stores/budgets'

const store = useTransactionStore()
const budgetStore = useBudgetStore()

const search = ref(''); const filterType = ref(''); const filterCategory = ref('')
const showModal = ref(false); const editing = ref(null)
const deleteTarget = ref(null); const saving = ref(false); const formError = ref('')

const categories = ['Food & Dining','Housing','Transport','Health','Entertainment','Shopping','Utilities','Education','Travel','Salary','Freelance','Investment','Other']
const form = reactive({ title: '', amount: '', type: 'expense', category: 'Food & Dining', date: today(), note: '', budget_id: null })

function today() { return new Date().toISOString().split('T')[0] }
const fmt = n => new Intl.NumberFormat('en-KE', { style: 'currency', currency: 'KES', currencyDisplay: 'narrowSymbol', maximumFractionDigits: 0 }).format(n || 0)
const fmtDate = d => new Date(d + 'T12:00').toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
const budgetName = id => id ? (budgetStore.budgets.find(b => b.id === id)?.name || '') : ''

const usedCats = computed(() => [...new Set(store.transactions.map(t => t.category))])
const filtered = computed(() => {
  let list = store.transactions
  if (search.value) { const q = search.value.toLowerCase(); list = list.filter(t => t.title.toLowerCase().includes(q) || t.category.toLowerCase().includes(q) || (t.note||'').toLowerCase().includes(q)) }
  if (filterType.value) list = list.filter(t => t.type === filterType.value)
  if (filterCategory.value) list = list.filter(t => t.category === filterCategory.value)
  return list
})
const filteredIncome = computed(() => filtered.value.filter(t => t.type === 'income').reduce((s,t) => s+t.amount, 0))
const filteredExpenses = computed(() => filtered.value.filter(t => t.type === 'expense').reduce((s,t) => s+t.amount, 0))
const filteredNet = computed(() => filteredIncome.value - filteredExpenses.value)

function openCreate() { editing.value = null; Object.assign(form, { title: '', amount: '', type: 'expense', category: 'Food & Dining', date: today(), note: '', budget_id: null }); formError.value = ''; showModal.value = true }
function openEdit(t) { editing.value = t.id; Object.assign(form, { title: t.title, amount: t.amount, type: t.type, category: t.category, date: t.date, note: t.note||'', budget_id: t.budget_id }); formError.value = ''; showModal.value = true }
function confirmDelete(t) { deleteTarget.value = t }

async function save() {
  if (!form.title || !form.amount) { formError.value = 'Description and amount are required'; return }
  saving.value = true; formError.value = ''
  try {
    const payload = { ...form, budget_id: form.budget_id || null }
    editing.value ? await store.updateTransaction(editing.value, payload) : await store.createTransaction(payload)
    showModal.value = false
  } catch(e) { formError.value = e.response?.data?.error || 'Failed to save' }
  finally { saving.value = false }
}
async function doDelete() {
  saving.value = true
  try { await store.deleteTransaction(deleteTarget.value.id); deleteTarget.value = null }
  finally { saving.value = false }
}
onMounted(async () => { await Promise.all([store.fetchTransactions(), budgetStore.fetchBudgets()]) })
</script>

<style scoped>
.eyebrow { font-size: 12px; color: var(--ink-3); margin-bottom: 4px; text-transform: uppercase; letter-spacing: .06em; font-weight: 500; }

/* Filter bar */
.filter-bar { display: flex; gap: 10px; margin-bottom: 14px; }
.search-wrap { position: relative; flex: 1; }
.search-icon { position: absolute; left: 12px; top: 50%; transform: translateY(-50%); color: var(--ink-3); pointer-events: none; }
.search-input {
  width: 100%; background: white; border: 1.5px solid var(--line); color: var(--ink);
  border-radius: 6px; padding: 9px 13px 9px 34px; font-size: 14px; outline: none;
  transition: border-color .14s; font-family: var(--font-ui);
}
.search-input:focus { border-color: var(--blue); }
.filt { width: 160px; flex-shrink: 0; }

/* Summary row */
.sum-row {
  display: flex; align-items: center; gap: 20px;
  background: white; border: 1px solid var(--line); border-radius: 8px;
  padding: 13px 20px; margin-bottom: 20px;
}
.sum-item { display: flex; flex-direction: column; gap: 1px; }
.sum-lbl { font-size: 11px; color: var(--ink-3); text-transform: uppercase; letter-spacing: .06em; }
.sum-val { font-size: 15px; font-weight: 500; }
.sum-val.income { color: var(--green); }
.sum-val.expense { color: var(--red); }
.sum-div { width: 1px; height: 30px; background: var(--line); }
.sum-count { margin-left: auto; font-size: 12px; color: var(--ink-3); font-family: var(--font-num); }

.skel-list { display: flex; flex-direction: column; gap: 6px; }
.skel { background: var(--line); border-radius: 6px; animation: pulse 1.5s ease-in-out infinite; }
@keyframes pulse { 0%,100%{opacity:.5} 50%{opacity:.25} }

/* Table */
.txn-panel { background: white; border: 1px solid var(--line); border-radius: 10px; overflow: hidden; }
.txn-table { width: 100%; border-collapse: collapse; }
.txn-table thead tr { border-bottom: 1px solid var(--line); }
.txn-table th { font-size: 11px; font-weight: 500; text-transform: uppercase; letter-spacing: .06em; color: var(--ink-3); padding: 12px 16px; text-align: left; }
.txn-table th.right { text-align: right; }
.txn-table tbody .txn-row { border-bottom: 1px solid var(--line); transition: background .1s; }
.txn-table tbody .txn-row:last-child { border-bottom: none; }
.txn-table tbody .txn-row:hover { background: var(--paper); }
.txn-table td { padding: 11px 16px; font-size: 13.5px; vertical-align: middle; }

.desc-cell { display: flex; align-items: center; gap: 10px; }
.type-pip { width: 6px; height: 6px; border-radius: 50%; flex-shrink: 0; }
.type-pip.income { background: var(--green); }
.type-pip.expense { background: var(--ink-3); }
.desc-text { font-weight: 500; }
.note-text { font-size: 12px; color: var(--ink-3); margin-left: 4px; }
.cat-cell { color: var(--ink-3); font-size: 13px; }
.bud-cell { color: var(--ink-3); font-size: 12.5px; }
.date-cell { color: var(--ink-3); font-size: 12.5px; }
.amount-cell { font-weight: 500; }
.amount-cell.income { color: var(--green); }
.amount-cell.expense { color: var(--ink); }
.right { text-align: right; }

.actions-cell { width: 72px; }
.row-actions { display: flex; gap: 2px; opacity: 0; transition: opacity .13s; }
.txn-row:hover .row-actions { opacity: 1; }
.icon-btn { background: none; border: none; color: var(--ink-3); padding: 5px; border-radius: 5px; display: flex; align-items: center; justify-content: center; transition: all .13s; }
.icon-btn:hover { background: var(--paper); color: var(--ink); }
.icon-btn.danger:hover { background: var(--red-dim); color: var(--red); }

/* Modal form */
.form-stack { display: flex; flex-direction: column; gap: 14px; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }

.type-toggle { display: flex; background: var(--paper); border-radius: 7px; padding: 3px; gap: 3px; }
.toggle-btn { flex: 1; padding: 7px; border: none; background: none; font-size: 13.5px; font-weight: 500; color: var(--ink-3); border-radius: 5px; transition: all .13s; }
.toggle-btn.active { background: white; color: var(--ink); box-shadow: 0 1px 3px rgba(0,0,0,.08); }
.toggle-btn.income-active.active { color: var(--green); }

.form-err { background: var(--red-dim); border: 1px solid #FECACA; color: var(--red); border-radius: 6px; padding: 9px 13px; font-size: 13px; }
.confirm-text { font-size: 14px; color: var(--ink-2); }
.confirm-text strong { color: var(--ink); }
</style>
