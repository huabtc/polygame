<template>
  <nav class="bg-white shadow-md">
    <div class="container mx-auto px-4">
      <div class="flex justify-between items-center h-16">
        <!-- Logo -->
        <router-link to="/" class="flex items-center space-x-2">
          <span class="text-2xl font-bold text-primary-600">Polygame</span>
        </router-link>

        <!-- Desktop Navigation -->
        <div class="hidden md:flex items-center space-x-6">
          <router-link to="/" class="nav-link">Home</router-link>
          <router-link to="/markets" class="nav-link">Markets</router-link>
          <router-link to="/portfolio" class="nav-link">Portfolio</router-link>
        </div>

        <!-- User Menu -->
        <div class="flex items-center space-x-4">
          <div class="text-sm">
            <span class="text-gray-600">Balance:</span>
            <span class="font-semibold text-primary-600 ml-1">
              {{ formatBalance(authStore.user?.virtual_balance) }}
            </span>
          </div>
          <button @click="showUserMenu = !showUserMenu" class="flex items-center space-x-2">
            <div class="w-8 h-8 bg-primary-500 rounded-full flex items-center justify-center text-white">
              {{ authStore.user?.username?.[0]?.toUpperCase() }}
            </div>
          </button>
        </div>

        <!-- Mobile Menu Button -->
        <button @click="showMobileMenu = !showMobileMenu" class="md:hidden">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
      </div>

      <!-- User Dropdown -->
      <div v-if="showUserMenu" class="absolute right-4 mt-2 w-48 bg-white rounded-lg shadow-lg py-2 z-50">
        <router-link to="/profile" class="block px-4 py-2 hover:bg-gray-100">Profile</router-link>
        <button @click="handleLogout" class="block w-full text-left px-4 py-2 hover:bg-gray-100 text-red-600">
          Logout
        </button>
      </div>

      <!-- Mobile Menu -->
      <div v-if="showMobileMenu" class="md:hidden py-4">
        <router-link to="/" class="block py-2">Home</router-link>
        <router-link to="/markets" class="block py-2">Markets</router-link>
        <router-link to="/portfolio" class="block py-2">Portfolio</router-link>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const showUserMenu = ref(false)
const showMobileMenu = ref(false)

const formatBalance = (balance) => {
  return balance?.toFixed(2) || '0.00'
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.nav-link {
  @apply text-gray-700 hover:text-primary-600 transition-colors;
}

.router-link-active {
  @apply text-primary-600 font-semibold;
}
</style>
