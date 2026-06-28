import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api/client'

export const useBudgetStore = defineStore('budgets', () => {
  const budgets = ref([])
  const loading = ref(false)

  async function fetchBudgets() {
    loading.value = true
    try {
      const res = await api.get('/budgets')
      budgets.value = res.data || []
    } finally {
      loading.value = false
    }
  }

  async function createBudget(data) {
    const res = await api.post('/budgets', data)
    budgets.value.unshift({ ...res.data, spent: 0 })
    return res.data
  }

  async function updateBudget(id, data) {
    const res = await api.put(`/budgets/${id}`, data)
    const idx = budgets.value.findIndex((b) => b.id === id)
    if (idx !== -1) budgets.value[idx] = { ...budgets.value[idx], ...res.data }
    return res.data
  }

  async function deleteBudget(id) {
    await api.delete(`/budgets/${id}`)
    budgets.value = budgets.value.filter((b) => b.id !== id)
  }

  return { budgets, loading, fetchBudgets, createBudget, updateBudget, deleteBudget }
})
