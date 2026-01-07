<template>
  <div class="container mx-auto px-4 py-8">
    <div v-if="loading" class="text-center py-12">
      <p class="text-gray-600">Loading...</p>
    </div>

    <div v-else-if="market" class="max-w-4xl mx-auto">
      <div class="card mb-6">
        <h1 class="text-3xl font-bold mb-4">{{ market.title }}</h1>
        <p class="text-gray-600 mb-4">{{ market.description }}</p>
        <div class="flex items-center space-x-4 text-sm">
          <span class="px-3 py-1 bg-gray-100 rounded">{{ market.category }}</span>
          <span class="px-3 py-1 rounded" :class="statusClass">{{ market.status }}</span>
        </div>
      </div>

      <!-- Outcomes -->
      <div class="card">
        <h2 class="text-2xl font-bold mb-6">Outcomes</h2>
        <div class="space-y-4">
          <div v-for="outcome in market.outcomes" :key="outcome.id" class="border rounded-lg p-4">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-semibold">{{ outcome.outcome_name }}</h3>
              <span class="text-2xl font-bold text-primary-600">
                {{ (outcome.current_price * 100).toFixed(1) }}%
              </span>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <button @click="openTradeModal(outcome, 'buy')" class="btn btn-success">
                Buy
              </button>
              <button @click="openTradeModal(outcome, 'sell')" class="btn btn-danger">
                Sell
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Trade Modal -->
      <div v-if="showTradeModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="closeTradeModal">
        <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
          <h3 class="text-xl font-bold mb-4">
            {{ tradeType === 'buy' ? 'Buy' : 'Sell' }} {{ selectedOutcome?.outcome_name }}
          </h3>

          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Shares</label>
              <input v-model.number="tradeForm.shares" type="number" min="1" step="1" class="input" />
            </div>

            <div>
              <label class="block text-sm font-medium mb-2">Price per share</label>
              <input v-model.number="tradeForm.price" type="number" min="0.01" max="1" step="0.01" class="input" />
            </div>

            <div class="p-3 bg-gray-50 rounded">
              <p class="text-sm text-gray-600">Total Cost: <span class="font-bold">{{ totalCost.toFixed(2) }}</span></p>
            </div>

            <div v-if="tradeError" class="p-3 bg-red-50 text-red-600 rounded text-sm">
              {{ tradeError }}
            </div>

            <div class="flex space-x-4">
              <button @click="handleTrade" :disabled="tradeLoading" class="flex-1 btn btn-primary">
                {{ tradeLoading ? 'Processing...' : 'Confirm' }}
              </button>
              <button @click="closeTradeModal" class="flex-1 btn btn-secondary">
                Cancel
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useMarketStore } from '@/stores/market'
import { useTradingStore } from '@/stores/trading'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const marketStore = useMarketStore()
const tradingStore = useTradingStore()
const authStore = useAuthStore()

const market = ref(null)
const loading = ref(false)
const showTradeModal = ref(false)
const selectedOutcome = ref(null)
const tradeType = ref('buy')
const tradeForm = ref({
  shares: 1,
  price: 0.5,
})
const tradeLoading = ref(false)
const tradeError = ref('')

const statusClass = computed(() => {
  const status = market.value?.status
  return {
    'bg-green-100 text-green-800': status === 'active',
    'bg-blue-100 text-blue-800': status === 'closed',
    'bg-purple-100 text-purple-800': status === 'resolved',
  }
})

const totalCost = computed(() => {
  return tradeForm.value.shares * tradeForm.value.price
})

const openTradeModal = (outcome, type) => {
  selectedOutcome.value = outcome
  tradeType.value = type
  tradeForm.value.price = outcome.current_price
  showTradeModal.value = true
  tradeError.value = ''
}

const closeTradeModal = () => {
  showTradeModal.value = false
  selectedOutcome.value = null
  tradeForm.value = { shares: 1, price: 0.5 }
}

const handleTrade = async () => {
  tradeLoading.value = true
  tradeError.value = ''

  const result = await tradingStore.placeOrder({
    market_id: market.value.id,
    outcome_id: selectedOutcome.value.id,
    order_type: tradeType.value,
    shares: tradeForm.value.shares,
    price: tradeForm.value.price,
  })

  if (result.success) {
    closeTradeModal()
    await authStore.fetchProfile()
    alert('Order placed successfully!')
  } else {
    tradeError.value = result.error
  }

  tradeLoading.value = false
}

onMounted(async () => {
  loading.value = true
  try {
    market.value = await marketStore.fetchMarket(route.params.id)
  } catch (error) {
    console.error('Failed to fetch market:', error)
  } finally {
    loading.value = false
  }
})
</script>
