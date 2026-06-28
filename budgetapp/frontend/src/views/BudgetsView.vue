<template>
  <AppShell>
    <div class="page-header">
      <h1>Budgets</h1>
      <button class="btn btn-primary" @click="openCreate">+ New Budget</button>
    </div>

    <div v-if="store.loading" class="loading-grid">
      <div v-for="i in 4" :key="i" class="skeleton-card"></div>
    </div>

    <div v-else-if="store.budgets.length === 0" class="empty-state">
      <div class="icon"></div>
      <strong>No budgets yet</strong>
      <p>Create a budget to start tracking your spending</p>
    </div>

    <div v-else class="budgets-grid">
      <div
        v-for="b in store.budgets"
        :key="b.id"
        class="budget-card card"
      >
        <div class="bc-header">
          <div class="bc-color" :style="{ background: b.color }"></div>
          <div class="bc-info">
            <h3>{{ b.name }}</h3>
            <span class="bc-cat">{{ b.category }}</span>
          </div>
          <div class="bc-actions">
            <button class="icon-btn" @click="openEdit(b)">✏️</button>
            <button class="icon-btn danger" @click="confirmDelete(b)">🗑️</button>
          </div>
        </div>

        <div class="bc-amounts">
          <div>
            <div class="label">Budget</div>
            <div class="amount">{{ fmt(b.amount) }}</div>
          </div>
          <div>
            <div class="label">Spent</div>
            <div class="amount spent" :class="{ over: b.spent > b.amount }">{{ fmt(b.spent) }}</div>
          </div>
          <div>
            <div class="label">Remaining</div>
            <div class="amount" :class="b.amount - b.spent >= 0 ? 'remain' : 'over'">
              {{ fmt(b.amount - b.spent) }}
            </div>
          </div>
        </div>

        <div class="progress-track">
          <div
            class="progress-fill"
            :style="{
              width: Math.min(100, pct(b)) + '%',
              background: pct(b) > 90 ? 'var(--red)' : pct(b) > 70 ? 'var(--yellow)' : b.color
            }"
          ></div>
        </div>
        <div class="bc-pct">{{ pct(b).toFixed(1) }}% of budget used</div>
      </div>
    </div>

    <!-- Create / Edit Modal -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ editing ? 'Edit Budget' : 'New Budget' }}</h3>
          <button class="icon-btn" @click="showModal = false">✕</button>
        </div>

        <div class="form-body">
          <div class="form-group">
            <label>Name</label>
            <input v-model="form.name" class="form-control" placeholder="e.g. Groceries" required />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Amount ($)</label>
              <input v-model.number="form.amount" type="number" min="0" step="0.01" class="form-control" placeholder="500" />
            </div>
            <div class="form-group">
              <label>Category</label>
              <select v-model="form.category" class="form-control">
                <option v-for="c in categories" :key="c">{{ c }}</option>
              </select>
            </div>
          </div>
          <div class="form-group">
            <label>Color</label>
            <div class="color-picker">
              <button
                v-for="c in colors"
                :key="c"
                class="color-swatch"
                :style="{ background: c, outline: form.color === c ? '2px solid white' : 'none' }"
                @click="form.color = c"
              ></button>
            </div>
          </div>
        </div>

        <div v-if="formError" class="form-error">{{ formError }}</div>
        <div class="modal-footer">
          <button class="btn btn-ghost" @click="showModal = false">Cancel</button>
          <button class="btn btn-primary" @click="save" :disabled="saving">
            {{ saving ? 'Saving...' : editing ? 'Update' : 'Create' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirm -->
    <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
      <div class="modal" style="max-width:380px">
        <div class="modal-header">
          <h3>Delete Budget</h3>
          <button class="icon-btn" @click="deleteTarget = null">✕</button>
        </div>
        <p style="color:var(--muted);font-size:14px">
          Are you sure you want to delete <strong style="color:var(--text)">{{ deleteTarget.name }}</strong>?
          This cannot be undone.
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
import { ref, reactive, onMounted } from 'vue'
import AppShell from '../components/AppShell.vue'
import { useBudgetStore } from '../stores/budgets'

const store = useBudgetStore()

const showModal = ref(false)
const editing = ref(null)
const deleteTarget = ref(null)
const saving = ref(false)
const formError = ref('')

const categories = ['Food & Dining', 'Housing', 'Transport', 'Health', 'Entertainment', 'Shopping', 'Utilities', 'Education', 'Travel', 'Other']
const colors = ['#7c6af7', '#22d3a0', '#f87171', '#fbbf24', '#38bdf8', '#f472b6', '#a3e635', '#fb923c', '#e879f9', '#34d399']

const form = reactive({ name: '', amount: '', category: 'Food & Dining', color: '#7c6af7' })

function pct(b) { return b.amount ? (b.spent / b.amount) * 100 : 0 }
function fmt(n) {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0 }).format(n || 0)
}

function openCreate() {
  editing.value = null
  Object.assign(form, { name: '', amount: '', category: 'Food & Dining', color: '#7c6af7' })
  formError.value = ''
  showModal.value = true
}
function openEdit(b) {
  editing.value = b.id
  Object.assign(form, { name: b.name, amount: b.amount, category: b.category, color: b.color })
  formError.value = ''
  showModal.value = true
}
function confirmDelete(b) { deleteTarget.value = b }

async function save() {
  if (!form.name || !form.amount) { formError.value = 'Name and amount are required'; return }
  saving.value = true
  formError.value = ''
  try {
    if (editing.value) {
      await store.updateBudget(editing.value, { ...form })
    } else {
      await store.createBudget({ ...form })
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
    await store.deleteBudget(deleteTarget.value.id)
    deleteTarget.value = null
  } finally {
    saving.value = false
  }
}

onMounted(() => store.fetchBudgets())
</script>

<style scoped>
.loading-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 16px; }
.skeleton-card { height: 180px; background: var(--surface); border-radius: var(--radius); border: 1px solid var(--border); animation: pulse 1.4s ease-in-out infinite; }
@keyframes pulse { 0%,100%{opacity:.6} 50%{opacity:.3} }

.budgets-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 16px; }
.budget-card { display: flex; flex-direction: column; gap: 14px; }

.bc-header { display: flex; align-items: center; gap: 12px; }
.bc-color { width: 12px; height: 38px; border-radius: 4px; flex-shrink: 0; }
.bc-info { flex: 1; min-width: 0; }
.bc-info h3 { font-size: 15px; font-weight: 600; }
.bc-cat { font-size: 11px; color: var(--muted); background: var(--surface2); padding: 2px 8px; border-radius: 20px; }
.bc-actions { display: flex; gap: 4px; }
.icon-btn {
  background: none; border: none; font-size: 14px;
  padding: 5px 7px; border-radius: 6px; cursor: pointer; opacity: .6; transition: all .15s;
}
.icon-btn:hover { opacity: 1; background: var(--surface2); }
.icon-btn.danger:hover { background: rgba(248,113,113,.1); }

.bc-amounts { display: grid; grid-template-columns: repeat(3,1fr); gap: 12px; }
.label { font-size: 11px; color: var(--muted); margin-bottom: 3px; }
.amount { font-size: 15px; font-weight: 600; }
.spent { color: var(--muted); }
.remain { color: var(--green); }
.over { color: var(--red); }

.progress-track { height: 6px; background: var(--surface2); border-radius: 3px; overflow: hidden; }
.progress-fill { height: 100%; border-radius: 3px; transition: width .4s; }
.bc-pct { font-size: 12px; color: var(--muted); }

.form-body { display: flex; flex-direction: column; gap: 16px; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.color-picker { display: flex; gap: 8px; flex-wrap: wrap; }
.color-swatch { width: 28px; height: 28px; border-radius: 50%; border: none; cursor: pointer; outline-offset: 3px; transition: transform .1s; }
.color-swatch:hover { transform: scale(1.15); }
.form-error { background: rgba(248,113,113,.12); border: 1px solid rgba(248,113,113,.3); color: var(--red); border-radius: 8px; padding: 9px 14px; font-size: 13px; }
</style>
