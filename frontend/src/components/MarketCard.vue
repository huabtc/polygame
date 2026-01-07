<template>
  <div class="card hover:shadow-lg transition-shadow cursor-pointer" @click="goToMarket">
    <div v-if="market.image_url" class="mb-4">
      <img :src="market.image_url" :alt="market.title" class="w-full h-48 object-cover rounded-lg" />
    </div>
    
    <div class="flex items-start justify-between mb-2">
      <h3 class="text-lg font-semibold flex-1">{{ market.title }}</h3>
      <span class="px-2 py-1 text-xs rounded-full" :class="statusClass">
        {{ market.status }}
      </span>
    </div>
    
    <p class="text-gray-600 text-sm mb-4 line-clamp-2">{{ market.description }}</p>
    
    <div class="flex items-center justify-between text-sm text-gray-500 mb-4">
      <span class="px-2 py-1 bg-gray-100 rounded">{{ market.category }}</span>
      <span>Volume: {{ formatVolume(market.total_volume) }}</span>
    </div>
    
    <div v-if="market.outcomes && market.outcomes.length > 0" class="space-y-2">
      <div v-for="outcome in market.outcomes.slice(0, 2)" :key="outcome.id" 
           class="flex items-center justify-between p-2 bg-gray-50 rounded">
        <span class="text-sm font-medium">{{ outcome.outcome_name }}</span>
        <span class="text-sm font-bold" :class="getPriceColor(outcome.current_price)">
          {{ (outcome.current_price * 100).toFixed(1) }}%
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  market: {
    type: Object,
    required: true,
  },
})

const router = useRouter()

const statusClass = computed(() => {
  const status = props.market.status
  return {
    'bg-green-100 text-green-800': status === 'active',
    'bg-gray-100 text-gray-800': status === 'pending',
    'bg-blue-100 text-blue-800': status === 'closed',
    'bg-purple-100 text-purple-800': status === 'resolved',
  }
})

const formatVolume = (volume) => {
  if (volume >= 1000000) return `${(volume / 1000000).toFixed(1)}M`
  if (volume >= 1000) return `${(volume / 1000).toFixed(1)}K`
  return volume?.toFixed(0) || '0'
}

const getPriceColor = (price) => {
  if (price >= 0.6) return 'text-green-600'
  if (price >= 0.4) return 'text-yellow-600'
  return 'text-red-600'
}

const goToMarket = () => {
  router.push(`/markets/${props.market.id}`)
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
