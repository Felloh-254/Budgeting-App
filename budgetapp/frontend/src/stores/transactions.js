import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api/client'

export const useTransactionStore = defineStore('transactions', () => {
  const transactions = ref([])
  const loading = ref(false)

  async function fetchTransactions() {
    loading.value = true
    try {
      const res = await api.get('/transactions')
      transactions.value = res.data || []
    } finally {
      loading.value = false
    }
  }

  async function createTransaction(data) {
    const res = await api.post('/transactions', data)
    transactions.value.unshift(res.data)
    return res.data
  }

  async function updateTransaction(id, data) {
    const res = await api.put(`/transactions/${id}`, data)
    const idx = transactions.value.findIndex((t) => t.id === id)
    if (idx !== -1) transactions.value[idx] = res.data
    return res.data
  }

  async function deleteTransaction(id) {
    await api.delete(`/transactions/${id}`)
    transactions.value = transactions.value.filter((t) => t.id !== id)
  }

  return { transactions, loading, fetchTransactions, createTransaction, updateTransaction, deleteTransaction }
})
