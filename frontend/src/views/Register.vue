<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4">
    <div class="max-w-md w-full">
      <div class="text-center mb-8">
        <h2 class="text-3xl font-bold text-gray-900">Create Account</h2>
        <p class="mt-2 text-gray-600">Join Polygame and start trading</p>
      </div>

      <div class="card">
        <form @submit.prevent="handleRegister" class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Username</label>
            <input
              v-model="form.username"
              type="text"
              required
              minlength="3"
              class="input"
              placeholder="Choose a username"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Email</label>
            <input
              v-model="form.email"
              type="email"
              required
              class="input"
              placeholder="Enter your email"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Password</label>
            <input
              v-model="form.password"
              type="password"
              required
              minlength="6"
              class="input"
              placeholder="Create a password (min 6 characters)"
            />
          </div>

          <div class="p-3 bg-yellow-50 text-yellow-800 rounded-lg text-sm">
            <p class="font-semibold mb-1">⚠️ Important Notice:</p>
            <ul class="text-xs space-y-1">
              <li>• You will receive 10,000 virtual credits upon registration</li>
              <li>• This is a simulation platform - no real money involved</li>
              <li>• For entertainment purposes only</li>
            </ul>
          </div>

          <div v-if="error" class="p-3 bg-red-50 text-red-600 rounded-lg text-sm">
            {{ error }}
          </div>

          <button type="submit" :disabled="loading" class="w-full btn btn-primary">
            {{ loading ? 'Creating account...' : 'Create Account' }}
          </button>
        </form>

        <div class="mt-6 text-center">
          <p class="text-sm text-gray-600">
            Already have an account?
            <router-link to="/login" class="text-primary-600 hover:text-primary-700 font-medium">
              Sign in
            </router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  username: '',
  email: '',
  password: '',
})

const loading = ref(false)
const error = ref('')

const handleRegister = async () => {
  loading.value = true
  error.value = ''

  const result = await authStore.register(
    form.value.username,
    form.value.email,
    form.value.password
  )

  if (result.success) {
    router.push('/')
  } else {
    error.value = result.error
  }

  loading.value = false
}
</script>
