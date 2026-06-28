<template>
  <AppShell>
    <div class="page-header">
      <h1>Transactions</h1>
      <button class="btn btn-primary" @click="openCreate">+ Add Transaction</button>
    </div>

    <!-- Filters -->
    <div class="filters card">
      <input v-model="search" class="form-control" placeholder="🔍  Search transactions..." style="flex:1" />
      <select v-model="filterType" class="form-control filter-select">
        <option value="">All types</option>
        <option value="income">Income</option>
        <option value="expense">Expense</option>
      </select>
      <select v-model="filterCategory" class="form-control filter-select">
        <option value="">All categories</option>
        <option v-for="c in usedCategories" :key="c">{{ c }}</option>
      </select>
    </div>

    <!-- Summary strip -->
    <div class="txn-summary">
      <div class="txn-sum-item">
        <span class="label">Filtered income</span>
        <span class="value green">{{ fmt(filteredIncome) }}</span>
      </div>
      <div class="divider"></div>
      <div class="txn-sum-item">
        <span class="label">Filtered expenses</span>
        <span class="value red">{{ fmt(filteredExpenses) }}</span>
      </div>
      <div class="divider"></div>
      <div class="txn-sum-item">
        <span class="label">Net</span>
        <span class="value" :class="filteredNet >= 0 ? 'green' : 'red'">{{ fmt(filteredNet) }}</span>
      </div>
    </div>

    <div v-if="store.loading" class="loading-list">
      <div v-for="i in 5" :key="i" class="skeleton-row"></div>
    </div>

    <div v-else-if="filtered.length === 0" class="empty-state">
      <div class="icon">💳</div>
      <strong>{{ store.transactions.length ? 'No matches' : 'No transactions yet' }}</strong>
      <p>{{ store.transactions.length ? 'Try adjusting your filters' : 'Add your first transaction' }}</p>
    </div>

    <div v-else class="txn-list card">
      <div
        v-for="t in filtered"
        :key="t.id"
        class="txn-row"
      >
        <div class="txn-icon" :class="t.type">
          {{ t.type === 'income' ? '↑' : '↓' }}
        </div>
        <div class="txn-main">
          <span class="txn-title">{{ t.title }}</span>
          <span class="txn-meta">{{ t.category }} · {{ fmtDate(t.date) }}<span v-if="t.note"> · {{ t.note }}</span></span>
        </div>
        <div class="txn-right">
          <span class="txn-amount" :class="t.type">
            {{ t.type === 'income' ? '+' : '-' }}{{ fmt(t.amount) }}
          </span>
          <div class="txn-actions">
            <button class="icon-btn" @click="openEdit(t)">✏️</button>
            <button class="icon-btn danger" @click="confirmDelete(t)">🗑️</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create / Edit Modal -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ editing ? 'Edit Transaction' : 'New Transaction' }}</h3>
          <button class="icon-btn" @click="showModal = false">✕</button>
        </div>

        <div class="form-body">
          <div class="form-group">
            <label>Title</label>
            <input v-model="form.title" class="form-control" placeholder="e.g. Grocery run" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Amount ($)</label>
              <input v-model.number="form.amount" type="number" min="0" step="0.01" class="form-control" placeholder="0.00" />
            </div>
            <div class="form-group">
              <label>Type</label>
              <select v-model="form.type" class="form-control">
                <option value="expense">Expense</option>
                <option value="income">Income</option>
              </select>
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
              <label>Date</label>
              <input v-model="form.date" type="date" class="form-control" />
            </div>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Budget (optional)</label>
              <select v-model="form.budget_id" class="form-control">
                <option :value="null">No budget</option>
                <option v-for="b in budgetStore.budgets" :key="b.id" :value="b.id">{{ b.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label>Note (optional)</label>
              <input v-model="form.note" class="form-control" placeholder="Optional note" />
            </div>
          </div>
        </div>

        <div v-if="formError" class="form-error">{{ formError }}</div>
        <div class="modal-footer">
          <button class="btn btn-ghost" @click="showModal = false">Cancel</button>
          <button class="btn btn-primary" @click="save" :disabled="saving">
            {{ saving ? 'Saving...' : editing ? 'Update' : 'Add' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirm -->
    <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
      <div class="modal" style="max-width:380px">
        <div class="modal-header">
          <h3>Delete Transaction</h3>
          <button class="icon-btn" @click="deleteTarget = null">✕</button>
        </div>
        <p style="color:var(--muted);font-size:14px">
          Delete <strong style="color:var(--text)">{{ deleteTarget.title }}</strong>?
        </p>
        <div class="modal-footer">
          <button class="btn btn-ghost" @click="deleteTarget = null">Cancel</button>
          <button class="btn btn-danger" @click="doDelete" :disabled="saving">Delete</button>
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

const search = ref('')
const filterType = ref('')
const filterCategory = ref('')
const showModal = ref(false)
const editing = ref(null)
const deleteTarget = ref(null)
const saving = ref(false)
const formError = ref('')

const categories = ['Food & Dining', 'Housing', 'Transport', 'Health', 'Entertainment', 'Shopping', 'Utilities', 'Education', 'Travel', 'Salary', 'Freelance', 'Investment', 'Other']

const form = reactive({
  title: '', amount: '', type: 'expense', category: 'Food & Dining',
  date: new Date().toISOString().split('T')[0], note: '', budget_id: null
})

const filtered = computed(() => {
  let list = store.transactions
  if (search.value) {
    const q = search.value.toLowerCase()
    list = list.filter(t => t.title.toLowerCase().includes(q) || t.category.toLowerCase().includes(q) || (t.note || '').toLowerCase().includes(q))
  }
  if (filterType.value) list = list.filter(t => t.type === filterType.value)
  if (filterCategory.value) list = list.filter(t => t.category === filterCategory.value)
  return list
})

const usedCategories = computed(() => [...new Set(store.transactions.map(t => t.category))])
const filteredIncome = computed(() => filtered.value.filter(t => t.type === 'income').reduce((s, t) => s + t.amount, 0))
const filteredExpenses = computed(() => filtered.value.filter(t => t.type === 'expense').reduce((s, t) => s + t.amount, 0))
const filteredNet = computed(() => filteredIncome.value - filteredExpenses.value)

function fmt(n) {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0 }).format(n || 0)
}
function fmtDate(d) {
  return new Date(d + 'T12:00:00').toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function openCreate() {
  editing.value = null
  Object.assign(form, {
    title: '', amount: '', type: 'expense', category: 'Food & Dining',
    date: new Date().toISOString().split('T')[0], note: '', budget_id: null
  })
  formError.value = ''
  showModal.value = true
}
function openEdit(t) {
  editing.value = t.id
  Object.assign(form, { title: t.title, amount: t.amount, type: t.type, category: t.category, date: t.date, note: t.note || '', budget_id: t.budget_id })
  formError.value = ''
  showModal.value = true
}
function confirmDelete(t) { deleteTarget.value = t }

async function save() {
  if (!form.title || !form.amount) { formError.value = 'Title and amount are required'; return }
  saving.value = true
  formError.value = ''
  try {
    const payload = { ...form, budget_id: form.budget_id || null }
    if (editing.value) {
      await store.updateTransaction(editing.value, payload)
    } else {
      await store.createTransaction(payload)
    }
    showModal.value = false
  } catch (e) {
    formError.value = e.response?.data?.error || 'Failed to save'
  } finally {
    saving.value = false
  }
}

async function doDelete() {
  saving.value = true
  try {
    await store.deleteTransaction(deleteTarget.value.id)
    deleteTarget.value = null
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await Promise.all([store.fetchTransactions(), budgetStore.fetchBudgets()])
})
</script>

<style scoped>
.filters { display: flex; gap: 12px; align-items: center; margin-bottom: 16px; padding: 14px 16px; }
.filter-select { width: 160px; flex-shrink: 0; }

.txn-summary {
  display: flex; align-items: center; gap: 20px;
  background: var(--surface); border: 1px solid var(--border);
  border-radius: var(--radius); padding: 14px 20px;
  margin-bottom: 20px;
}
.txn-sum-item { display: flex; flex-direction: column; gap: 2px; }
.txn-sum-item .label { font-size: 11px; color: var(--muted); text-transform: uppercase; letter-spacing: .4px; }
.txn-sum-item .value { font-size: 16px; font-weight: 700; }
.green { color: var(--green); }
.red { color: var(--red); }
.divider { width: 1px; height: 36px; background: var(--border); }

.loading-list { display: flex; flex-direction: column; gap: 8px; }
.skeleton-row { height: 64px; background: var(--surface); border-radius: var(--radius); border: 1px solid var(--border); animation: pulse 1.4s ease-in-out infinite; }
@keyframes pulse { 0%,100%{opacity:.6} 50%{opacity:.3} }

.txn-list { padding: 8px; display: flex; flex-direction: column; gap: 2px; }
.txn-row {
  display: flex; align-items: center; gap: 14px;
  padding: 12px 12px; border-radius: 8px; transition: background .12s;
}
.txn-row:hover { background: var(--surface2); }
.txn-icon {
  width: 36px; height: 36px; border-radius: 50%; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  font-weight: 700; font-size: 16px;
}
.txn-icon.income { background: rgba(34,211,160,.15); color: var(--green); }
.txn-icon.expense { background: rgba(124,106,247,.15); color: var(--accent2); }

.txn-main { flex: 1; min-width: 0; }
.txn-title { display: block; font-size: 14px; font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.txn-meta { font-size: 12px; color: var(--muted); }

.txn-right { display: flex; align-items: center; gap: 12px; flex-shrink: 0; }
.txn-amount { font-size: 15px; font-weight: 600; min-width: 90px; text-align: right; }
.txn-amount.income { color: var(--green); }
.txn-amount.expense { color: var(--red); }
.txn-actions { display: flex; gap: 2px; opacity: 0; transition: opacity .15s; }
.txn-row:hover .txn-actions { opacity: 1; }
.icon-btn { background: none; border: none; font-size: 13px; padding: 4px 6px; border-radius: 5px; cursor: pointer; opacity: .7; transition: all .12s; }
.icon-btn:hover { opacity: 1; background: var(--surface); }

.form-body { display: flex; flex-direction: column; gap: 16px; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.form-error { background: rgba(248,113,113,.12); border: 1px solid rgba(248,113,113,.3); color: var(--red); border-radius: 8px; padding: 9px 14px; font-size: 13px; }
</style>
