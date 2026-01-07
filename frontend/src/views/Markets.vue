<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">Markets</h1>

    <!-- Filters -->
    <div class="card mb-8">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="block text-sm font-medium mb-2">Category</label>
          <select v-model="filters.category" class="input">
            <option value="">All Categories</option>
            <option value="sports">Sports</option>
            <option value="esports">Esports</option>
            <option value="entertainment">Entertainment</option>
            <option value="tech">Tech</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium mb-2">Status</label>
          <select v-model="filters.status" class="input">
            <option value="">All Status</option>
            <option value="active">Active</option>
            <option value="closed">Closed</option>
            <option value="resolved">Resolved</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium mb-2">Search</label>
          <input v-model="searchKeyword" type="text" class="input" placeholder="Search markets..." />
        </div>
      </div>
    </div>

    <!-- Markets Grid -->
    <div v-if="loading" class="text-center py-12">
      <p class="text-gray-600">Loading markets...</p>
    </div>

    <div v-else-if="markets.length === 0" class="text-center py-12">
      <p class="text-gray-600">No markets found</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <MarketCard v-for="market in markets" :key="market.id" :market="market" />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useMarketStore } from '@/stores/market'
import MarketCard from '@/components/MarketCard.vue'

const marketStore = useMarketStore()
const markets = ref([])
const loading = ref(false)
const searchKeyword = ref('')
const filters = ref({
  category: '',
  status: '',
})

const fetchMarkets = async () => {
  loading.value = true
  try {
    const data = await marketStore.fetchMarkets(filters.value)
    markets.value = data.markets
  } catch (error) {
    console.error('Failed to fetch markets:', error)
  } finally {
    loading.value = false
  }
}

watch(filters, fetchMarkets, { deep: true })

onMounted(fetchMarkets)
</script>
