<template>
  <div class="min-h-screen">
    <!-- Hero Section -->
    <section class="bg-gradient-to-r from-primary-600 to-primary-800 text-white py-20">
      <div class="container mx-auto px-4 text-center">
        <h1 class="text-5xl font-bold mb-6">Welcome to Polygame</h1>
        <p class="text-xl mb-8 max-w-2xl mx-auto">
          A simulation prediction market platform. Trade with virtual currency and test your prediction skills!
        </p>
        <div class="flex justify-center space-x-4">
          <router-link v-if="!authStore.isAuthenticated" to="/register" class="btn btn-primary bg-white text-primary-600 hover:bg-gray-100">
            Get Started
          </router-link>
          <router-link v-else to="/markets" class="btn btn-primary bg-white text-primary-600 hover:bg-gray-100">
            Explore Markets
          </router-link>
        </div>
        
        <!-- Disclaimer -->
        <div class="mt-12 p-4 bg-white/10 rounded-lg max-w-2xl mx-auto">
          <p class="text-sm">
            ⚠️ <strong>Important:</strong> This is a simulation platform using virtual currency only. 
            No real money involved. For entertainment and educational purposes only.
          </p>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-16 bg-white">
      <div class="container mx-auto px-4">
        <h2 class="text-3xl font-bold text-center mb-12">How It Works</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div class="text-center">
            <div class="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <span class="text-3xl">📊</span>
            </div>
            <h3 class="text-xl font-semibold mb-2">Browse Markets</h3>
            <p class="text-gray-600">Explore prediction markets on sports, esports, entertainment, and more</p>
          </div>
          <div class="text-center">
            <div class="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <span class="text-3xl">💰</span>
            </div>
            <h3 class="text-xl font-semibold mb-2">Trade Shares</h3>
            <p class="text-gray-600">Buy and sell outcome shares with virtual currency</p>
          </div>
          <div class="text-center">
            <div class="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <span class="text-3xl">🎯</span>
            </div>
            <h3 class="text-xl font-semibold mb-2">Win Rewards</h3>
            <p class="text-gray-600">Earn virtual credits when your predictions are correct</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Trending Markets -->
    <section v-if="authStore.isAuthenticated" class="py-16 bg-gray-50">
      <div class="container mx-auto px-4">
        <h2 class="text-3xl font-bold mb-8">Trending Markets</h2>
        <div v-if="loading" class="text-center py-8">
          <p class="text-gray-600">Loading...</p>
        </div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <MarketCard v-for="market in trendingMarkets" :key="market.id" :market="market" />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useMarketStore } from '@/stores/market'
import MarketCard from '@/components/MarketCard.vue'

const authStore = useAuthStore()
const marketStore = useMarketStore()
const trendingMarkets = ref([])
const loading = ref(false)

onMounted(async () => {
  if (authStore.isAuthenticated) {
    loading.value = true
    try {
      trendingMarkets.value = await marketStore.fetchTrendingMarkets(6)
    } catch (error) {
      console.error('Failed to fetch trending markets:', error)
    } finally {
      loading.value = false
    }
  }
})
</script>
