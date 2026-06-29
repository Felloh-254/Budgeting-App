<template>
  <AppShell>
    <div class="page-header">
      <div>
        <p class="eyebrow">Spending limits</p>
        <h1>Budgets</h1>
      </div>
      <button class="btn btn-primary" @click="openCreate">+ New budget</button>
    </div>

    <div v-if="store.loading" class="skel-grid">
      <div v-for="i in 4" :key="i" class="skel" style="height:160px"></div>
    </div>

    <div v-else-if="!store.budgets.length" class="empty-state" style="margin-top:60px">
      <div class="icon">🗂️</div>
      <strong>No budgets yet</strong>
      <p>Create a budget to start tracking your spending</p>
    </div>

    <div v-else class="budget-grid">
      <div v-for="b in store.budgets" :key="b.id" class="budget-card">
        <div class="bc-accent" :style="{ background: b.color }"></div>
        <div class="bc-body">
          <div class="bc-top">
            <div>
              <div class="bc-name">{{ b.name }}</div>
              <div class="bc-cat">{{ b.category }}</div>
            </div>
            <div class="bc-menu">
              <button class="icon-btn" @click="openEdit(b)" title="Edit">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
              </button>
              <button class="icon-btn danger" @click="confirmDelete(b)" title="Delete">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/><path d="M10 11v6M14 11v6"/><path d="M9 6V4h6v2"/></svg>
              </button>
            </div>
          </div>

          <div class="bc-amounts">
            <div>
              <div class="amt-label">Spent</div>
              <div class="amt-val num" :class="b.spent > b.amount ? 'over' : ''">{{ fmt(b.spent) }}</div>
            </div>
            <div>
              <div class="amt-label">Budget</div>
              <div class="amt-val num">{{ fmt(b.amount) }}</div>
            </div>
            <div>
              <div class="amt-label">Left</div>
              <div class="amt-val num" :class="b.amount - b.spent < 0 ? 'over' : 'remain'">{{ fmt(b.amount - b.spent) }}</div>
            </div>
          </div>

          <div class="prog-track">
            <div class="prog-fill" :style="{
              width: Math.min(100, pct(b)) + '%',
              background: pct(b) > 90 ? 'var(--red)' : pct(b) > 70 ? '#D97706' : b.color
            }"></div>
          </div>
          <div class="bc-pct" :class="{ warn: pct(b) > 70, danger: pct(b) > 90 }">
            {{ pct(b).toFixed(0) }}% used
          </div>
        </div>
      </div>
    </div>

    <!-- Create / Edit Modal -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ editing ? 'Edit budget' : 'New budget' }}</h3>
          <button class="icon-btn" @click="showModal = false">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </div>
        <div class="form-stack">
          <div class="form-group">
            <label>Name</label>
            <input v-model="form.name" class="form-control" placeholder="e.g. Groceries" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Amount (Ksh)</label>
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
            <div class="color-row">
              <button v-for="c in colors" :key="c" class="cswatch"
                :style="{ background: c, outline: form.color === c ? `2px solid ${c}` : 'none', outlineOffset: '2px' }"
                @click="form.color = c">
              </button>
            </div>
          </div>
        </div>
        <div v-if="formError" class="form-err">{{ formError }}</div>
        <div class="modal-footer">
          <button class="btn btn-ghost btn-sm" @click="showModal = false">Cancel</button>
          <button class="btn btn-primary btn-sm" @click="save" :disabled="saving">{{ saving ? 'Saving…' : editing ? 'Update' : 'Create' }}</button>
        </div>
      </div>
    </div>

    <!-- Delete confirm -->
    <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
      <div class="modal" style="max-width:380px">
        <div class="modal-header">
          <h3>Delete budget</h3>
          <button class="icon-btn" @click="deleteTarget = null">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </div>
        <p class="confirm-text">Delete <strong>{{ deleteTarget.name }}</strong>? This can't be undone.</p>
        <div class="modal-footer">
          <button class="btn btn-ghost btn-sm" @click="deleteTarget = null">Cancel</button>
          <button class="btn btn-danger btn-sm" @click="doDelete" :disabled="saving">Delete</button>
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
const showModal = ref(false); const editing = ref(null)
const deleteTarget = ref(null); const saving = ref(false); const formError = ref('')

const categories = ['Food & Dining','Housing','Transport','Health','Entertainment','Shopping','Utilities','Education','Travel','Other']
const colors = ['#2D5BE3','#16A34A','#DC2626','#D97706','#7C3AED','#0891B2','#BE185D','#059669','#EA580C','#6B7280']
const form = reactive({ name: '', amount: '', category: 'Food & Dining', color: '#2D5BE3' })

const pct = b => b.amount ? (b.spent / b.amount) * 100 : 0
const fmt = n => new Intl.NumberFormat('en-KE', { style: 'currency', currency: 'KES', currencyDisplay: 'narrowSymbol', maximumFractionDigits: 0 }).format(n || 0)

function openCreate() { editing.value = null; Object.assign(form, { name: '', amount: '', category: 'Food & Dining', color: '#2D5BE3' }); formError.value = ''; showModal.value = true }
function openEdit(b) { editing.value = b.id; Object.assign(form, { name: b.name, amount: b.amount, category: b.category, color: b.color }); formError.value = ''; showModal.value = true }
function confirmDelete(b) { deleteTarget.value = b }

async function save() {
  if (!form.name || !form.amount) { formError.value = 'Name and amount are required'; return }
  saving.value = true; formError.value = ''
  try {
    editing.value ? await store.updateBudget(editing.value, { ...form }) : await store.createBudget({ ...form })
    showModal.value = false
  } catch (e) { formError.value = e.response?.data?.error || 'Failed to save' }
  finally { saving.value = false }
}
async function doDelete() {
  saving.value = true
  try { await store.deleteBudget(deleteTarget.value.id); deleteTarget.value = null }
  finally { saving.value = false }
}
onMounted(() => store.fetchBudgets())
</script>

<style scoped>
.eyebrow { font-size: 12px; color: var(--ink-3); margin-bottom: 4px; text-transform: uppercase; letter-spacing: .06em; font-weight: 500; }

.skel-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 16px; }
.skel { background: var(--line); border-radius: 10px; animation: pulse 1.5s ease-in-out infinite; }
@keyframes pulse { 0%,100%{opacity:.5} 50%{opacity:.25} }

.budget-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 16px; }
.budget-card { background: white; border: 1px solid var(--line); border-radius: 10px; overflow: hidden; }
.bc-accent { height: 4px; }
.bc-body { padding: 18px 20px; display: flex; flex-direction: column; gap: 14px; }

.bc-top { display: flex; justify-content: space-between; align-items: flex-start; }
.bc-name { font-size: 15px; font-weight: 600; letter-spacing: -.01em; }
.bc-cat { font-size: 12px; color: var(--ink-3); margin-top: 2px; }

.bc-menu { display: flex; gap: 4px; }
.icon-btn { background: none; border: none; color: var(--ink-3); padding: 5px; border-radius: 5px; display: flex; align-items: center; justify-content: center; transition: all .13s; }
.icon-btn:hover { background: var(--paper); color: var(--ink); }
.icon-btn.danger:hover { background: var(--red-dim); color: var(--red); }

.bc-amounts { display: grid; grid-template-columns: repeat(3,1fr); }
.amt-label { font-size: 11px; color: var(--ink-3); text-transform: uppercase; letter-spacing: .05em; margin-bottom: 3px; }
.amt-val { font-size: 14.5px; font-weight: 500; }
.amt-val.over { color: var(--red); }
.amt-val.remain { color: var(--green); }

.prog-track { height: 4px; background: var(--line); border-radius: 2px; overflow: hidden; }
.prog-fill { height: 100%; border-radius: 2px; transition: width .5s cubic-bezier(.4,0,.2,1); }
.bc-pct { font-size: 11.5px; color: var(--ink-3); font-family: var(--font-num); }
.bc-pct.warn { color: #D97706; }
.bc-pct.danger { color: var(--red); }

.form-stack { display: flex; flex-direction: column; gap: 14px; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.color-row { display: flex; gap: 8px; flex-wrap: wrap; }
.cswatch { width: 26px; height: 26px; border-radius: 50%; border: none; cursor: pointer; transition: transform .12s; }
.cswatch:hover { transform: scale(1.15); }
.form-err { background: var(--red-dim); border: 1px solid #FECACA; color: var(--red); border-radius: 6px; padding: 9px 13px; font-size: 13px; }
.confirm-text { font-size: 14px; color: var(--ink-2); margin-bottom: 4px; }
.confirm-text strong { color: var(--ink); }
</style>
