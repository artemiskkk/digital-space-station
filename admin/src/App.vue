<template>
  <RouterView v-if="auth.ready" />
  <div v-else class="boot">
    <div class="boot-spinner"></div>
    <span>INITIALIZING...</span>
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import { onMounted } from 'vue'
import { useAuthStore } from './stores/auth'

const auth = useAuthStore()
onMounted(() => auth.ensureAuth())
</script>

<style>
.boot {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;
  font-family: 'DM Mono', monospace;
  font-size: 10px;
  letter-spacing: 0.2em;
  color: rgba(255,255,255,0.1);
}
.boot-spinner {
  width: 20px; height: 20px;
  border: 2px solid rgba(34,211,238,0.15);
  border-top-color: #22d3ee;
  border-radius: 50%;
  animation: bootspin 0.6s linear infinite;
}
@keyframes bootspin { to { transform: rotate(360deg); } }
</style>
