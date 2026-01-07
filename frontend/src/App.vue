<template>
  <div id="app" class="min-h-screen flex flex-col">
    <Navbar v-if="authStore.isAuthenticated" />
    <main class="flex-1">
      <router-view />
    </main>
    <Footer />
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import Navbar from '@/components/Navbar.vue'
import Footer from '@/components/Footer.vue'

const authStore = useAuthStore()
</script>
