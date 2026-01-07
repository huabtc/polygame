import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/axios'

export const useMarketStore = defineStore('market', () => {
  const markets = ref([])
  const currentMarket = ref(null)
  const loading = ref(false)

  async function fetchMarkets(params = {}) {
    loading.value = true
    try {
      const data = await api.get('/markets', { params })
      markets.value = data.markets
      return data
    } catch (error) {
      console.error('Failed to fetch markets:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function fetchMarket(id) {
    loading.value = true
    try {
      const data = await api.get(`/markets/${id}`)
      currentMarket.value = data.market
      return data.market
    } catch (error) {
      console.error('Failed to fetch market:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function fetchTrendingMarkets(limit = 10) {
    try {
      const data = await api.get('/markets/trending', { params: { limit } })
      return data.markets
    } catch (error) {
      console.error('Failed to fetch trending markets:', error)
      throw error
    }
  }

  async function searchMarkets(keyword, params = {}) {
    loading.value = true
    try {
      const data = await api.get('/markets/search', { params: { q: keyword, ...params } })
      markets.value = data.markets
      return data
    } catch (error) {
      console.error('Failed to search markets:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  return {
    markets,
    currentMarket,
    loading,
    fetchMarkets,
    fetchMarket,
    fetchTrendingMarkets,
    searchMarkets,
  }
})
