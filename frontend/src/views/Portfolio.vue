<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">My Portfolio</h1>

    <!-- Balance Card -->
    <div class="card mb-8">
      <h2 class="text-xl font-semibold mb-4">Virtual Balance</h2>
      <p class="text-4xl font-bold text-primary-600">
        {{ authStore.user?.virtual_balance?.toFixed(2) || '0.00' }}
      </p>
    </div>

    <!-- Positions -->
    <div class="card mb-8">
      <h2 class="text-xl font-semibold mb-4">My Positions</h2>
      
      <div v-if="loading" class="text-center py-8">
        <p class="text-gray-600">Loading...</p>
      </div>

      <div v-else-if="positions.length === 0" class="text-center py-8">
        <p class="text-gray-600">No positions yet</p>
      </div>

      <div v-else class="space-y-4">
        <div v-for="position in positions" :key="position.id" class="border rounded-lg p-4">
          <div class="flex items-start justify-between mb-2">
            <div>
              <h3 class="font-semibold">{{ position.market?.title }}</h3>
              <p class="text-sm text-gray-600">{{ position.outcome?.outcome_name }}</p>
            </div>
            <span class="px-2 py-1 bg-primary-100 text-primary-800 rounded text-sm">
              {{ position.shares }} shares
            </span>
          </div>
          <div class="text-sm text-gray-600">
            Avg Price: {{ position.avg_price?.toFixed(4) }}
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Orders -->
    <div class="card">
      <h2 class="text-xl font-semibold mb-4">Recent Orders</h2>
      
      <div v-if="ordersLoading" class="text-center py-8">
        <p class="text-gray-600">Loading...</p>
      </div>

      <div v-else-if="orders.length === 0" class="text-center py-8">
        <p class="text-gray-600">No orders yet</p>
      </div>

      <div v-else class="space-y-4">
        <div v-for="order in orders" :key="order.id" class="border rounded-lg p-4">
          <div class="flex items-start justify-between">
            <div>
              <h3 class="font-semibold">{{ order.market?.title }}</h3>
              <p class="text-sm text-gray-600">{{ order.outcome?.outcome_name }}</p>
            </div>
            <div class="text-right">
              <span class="px-2 py-1 rounded text-sm" :class="orderTypeClass(order.order_type)">
                {{ order.order_type }}
              </span>
              <p class="text-sm text-gray-600 mt-1">{{ order.shares }} @ {{ order.price }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useTradingStore } from '@/stores/trading'

const authStore = useAuthStore()
const tradingStore = useTradingStore()

const positions = ref([])
const orders = ref([])
const loading = ref(false)
const ordersLoading = ref(false)

const orderTypeClass = (type) => {
  return type === 'buy' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
}

onMounted(async () => {
  loading.value = true
  ordersLoading.value = true

  try {
    positions.value = await tradingStore.fetchPositions()
  } catch (error) {
    console.error('Failed to fetch positions:', error)
  } finally {
    loading.value = false
  }

  try {
    const data = await tradingStore.fetchOrders({ page: 1, page_size: 10 })
    orders.value = data.orders
  } catch (error) {
    console.error('Failed to fetch orders:', error)
  } finally {
    ordersLoading.value = false
  }
})
</script>
