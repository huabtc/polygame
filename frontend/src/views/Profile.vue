<template>
  <div class="container mx-auto px-4 py-8 max-w-2xl">
    <h1 class="text-3xl font-bold mb-8">Profile</h1>

    <div class="card">
      <div class="flex items-center space-x-4 mb-6">
        <div class="w-20 h-20 bg-primary-500 rounded-full flex items-center justify-center text-white text-3xl">
          {{ authStore.user?.username?.[0]?.toUpperCase() }}
        </div>
        <div>
          <h2 class="text-2xl font-bold">{{ authStore.user?.username }}</h2>
          <p class="text-gray-600">{{ authStore.user?.email }}</p>
        </div>
      </div>

      <div class="border-t pt-6 space-y-4">
        <div class="flex justify-between">
          <span class="text-gray-600">Virtual Balance</span>
          <span class="font-bold text-primary-600">
            {{ authStore.user?.virtual_balance?.toFixed(2) }}
          </span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-600">Member Since</span>
          <span class="font-medium">
            {{ formatDate(authStore.user?.created_at) }}
          </span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-600">Account Type</span>
          <span class="font-medium">
            {{ authStore.user?.is_admin ? 'Admin' : 'User' }}
          </span>
        </div>
      </div>

      <div class="border-t pt-6 mt-6">
        <div class="p-4 bg-yellow-50 rounded-lg">
          <p class="text-sm text-yellow-800">
            <strong>⚠️ Reminder:</strong> This is a simulation platform. All balances are virtual currency for entertainment purposes only.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleDateString()
}
</script>
