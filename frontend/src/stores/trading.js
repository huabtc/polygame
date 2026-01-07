import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/axios'

export const useTradingStore = defineStore('trading', () => {
  const orders = ref([])
  const positions = ref([])
  const loading = ref(false)

  async function placeOrder(orderData) {
    loading.value = true
    try {
      const data = await api.post('/trading/orders', orderData)
      return { success: true, order: data.order }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Order failed' }
    } finally {
      loading.value = false
    }
  }

  async function fetchOrders(params = {}) {
    loading.value = true
    try {
      const data = await api.get('/trading/orders', { params })
      orders.value = data.orders
      return data
    } catch (error) {
      console.error('Failed to fetch orders:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function fetchPositions() {
    loading.value = true
    try {
      const data = await api.get('/trading/positions')
      positions.value = data.positions
      return data.positions
    } catch (error) {
      console.error('Failed to fetch positions:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function cancelOrder(orderId) {
    try {
      await api.delete(`/trading/orders/${orderId}`)
      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Cancel failed' }
    }
  }

  return {
    orders,
    positions,
    loading,
    placeOrder,
    fetchOrders,
    fetchPositions,
    cancelOrder,
  }
})
