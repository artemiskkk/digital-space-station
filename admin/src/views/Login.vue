<template>
  <div class="login-scene">
    <!-- Grid bg -->
    <div class="grid-bg" aria-hidden="true"></div>
    <div class="glow" aria-hidden="true"></div>

    <div class="terminal">
      <!-- Top bar -->
      <div class="terminal-bar">
        <span class="bar-dot"></span>
        <span class="bar-title">DSS://AUTH_GATEWAY</span>
      </div>

      <!-- Logo -->
      <div class="logo-area">
        <svg width="44" height="44" viewBox="0 0 22 22" fill="none">
          <circle cx="11" cy="11" r="9.5" stroke="rgba(255,255,255,0.12)" stroke-width="0.8"/>
          <circle cx="11" cy="11" r="5.5" stroke="rgba(255,255,255,0.06)" stroke-width="0.8" stroke-dasharray="2 3"/>
          <circle cx="11" cy="11" r="2" fill="#22d3ee"/>
          <ellipse cx="11" cy="11" rx="9.5" ry="3.5" stroke="#22d3ee" stroke-width="0.8" stroke-opacity="0.45" transform="rotate(-30 11 11)"/>
        </svg>
        <h1>DIGITAL SPACE STATION</h1>
        <p class="sub">COMMAND CENTER · AUTHORIZATION REQUIRED</p>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleLogin">
        <label>
          <span class="label-text">ACCESS KEY</span>
          <div class="input-wrap" :class="{ focus: inputFocus }">
            <svg width="14" height="14" viewBox="0 0 14 14" fill="none"><rect x="3" y="5" width="8" height="7" rx="1.5" stroke="currentColor" stroke-width="1.1"/><path d="M5 5V3.5a2 2 0 0 1 4 0V5" stroke="currentColor" stroke-width="1.1" stroke-linecap="round"/></svg>
            <input
              v-model="password"
              type="password"
              placeholder="输入管理员密码..."
              autofocus
              :disabled="loading"
              @focus="inputFocus = true"
              @blur="inputFocus = false"
            />
          </div>
        </label>

        <button type="submit" :disabled="loading || !password" class="submit-btn">
          <template v-if="loading">
            <span class="spinner"></span>
            AUTHENTICATING...
          </template>
          <template v-else>
            AUTHORIZE ACCESS
            <svg width="14" height="14" viewBox="0 0 14 14" fill="none"><path d="M2 7h10M8 3l4 4-4 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </template>
        </button>
      </form>

      <p v-if="error" class="error">
        <svg width="12" height="12" viewBox="0 0 12 12" fill="none"><circle cx="6" cy="6" r="5" stroke="currentColor" stroke-width="1.1"/><path d="M6 3.5v3M6 8h.005" stroke="currentColor" stroke-width="1.2" stroke-linecap="round"/></svg>
        {{ error }}
      </p>

      <!-- Footer -->
      <div class="terminal-foot">
        <span class="foot-dot"></span>
        SYSTEM ONLINE · ENCRYPTED CONNECTION
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const password = ref('')
const loading = ref(false)
const error = ref('')
const inputFocus = ref(false)

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(password.value)
    router.push('/posts')
  } catch {
    error.value = 'ACCESS DENIED — 密码错误'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-scene {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.grid-bg {
  position: fixed; inset: 0;
  background-image:
    linear-gradient(rgba(255,255,255,0.025) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255,255,255,0.025) 1px, transparent 1px);
  background-size: 48px 48px;
  mask-image: radial-gradient(ellipse 70% 60% at 50% 50%, black 20%, transparent 100%);
  -webkit-mask-image: radial-gradient(ellipse 70% 60% at 50% 50%, black 20%, transparent 100%);
}
.glow {
  position: fixed;
  top: 30%; left: 50%;
  transform: translate(-50%,-50%);
  width: 500px; height: 500px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(34,211,238,0.06) 0%, transparent 70%);
  pointer-events: none;
}

.terminal {
  position: relative;
  width: 420px;
  background: rgba(13,17,23,0.85);
  border: 1px solid var(--border2);
  border-radius: 14px;
  backdrop-filter: blur(20px);
  overflow: hidden;
  box-shadow:
    0 0 0 1px rgba(34,211,238,0.05),
    0 30px 80px rgba(0,0,0,0.5),
    0 0 120px rgba(34,211,238,0.03);
}

.terminal-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--border);
  font-family: 'DM Mono', monospace;
  font-size: 10px;
  letter-spacing: 0.15em;
  color: var(--text2);
}
.bar-dot {
  width: 6px; height: 6px;
  border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 6px var(--accent-glow);
  animation: pulse 2.5s infinite;
}
@keyframes pulse { 0%,100%{opacity:1} 50%{opacity:.3} }

.logo-area {
  text-align: center;
  padding: 32px 24px 24px;
}
.logo-area h1 {
  font-family: 'Orbitron', sans-serif;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.2em;
  color: var(--text-h);
  margin-top: 16px;
}
.sub {
  font-family: 'DM Mono', monospace;
  font-size: 9px;
  letter-spacing: 0.2em;
  color: var(--text2);
  margin-top: 8px;
}

form {
  padding: 0 24px 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
label { display: flex; flex-direction: column; gap: 8px; }
.label-text {
  font-family: 'DM Mono', monospace;
  font-size: 10px;
  letter-spacing: 0.18em;
  color: var(--text2);
}

.input-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(255,255,255,0.03);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 0 14px;
  color: var(--text2);
  transition: all 0.2s;
}
.input-wrap.focus {
  border-color: rgba(34,211,238,0.3);
  box-shadow: 0 0 0 3px rgba(34,211,238,0.06);
  color: var(--accent);
}
.input-wrap input {
  flex: 1;
  background: none;
  border: none;
  padding: 12px 0;
  color: var(--text-h);
  font-size: 14px;
  outline: none;
  letter-spacing: 0.08em;
}
.input-wrap input::placeholder { color: var(--text2); }

.submit-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 12px;
  background: var(--accent);
  color: var(--bg);
  border: none;
  border-radius: 8px;
  font-family: 'DM Mono', monospace;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.15em;
  cursor: pointer;
  transition: all 0.2s;
}
.submit-btn:hover:not(:disabled) {
  box-shadow: 0 0 24px var(--accent-glow);
}
.submit-btn:disabled { opacity: 0.4; cursor: not-allowed; }

.spinner {
  width: 14px; height: 14px;
  border: 2px solid rgba(8,12,20,0.3);
  border-top-color: var(--bg);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

.error {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: -4px 24px 16px;
  font-size: 12px;
  color: var(--danger);
  font-family: 'DM Mono', monospace;
  letter-spacing: 0.05em;
}

.terminal-foot {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 12px 20px;
  border-top: 1px solid var(--border);
  font-family: 'DM Mono', monospace;
  font-size: 9px;
  letter-spacing: 0.15em;
  color: rgba(255,255,255,0.08);
}
.foot-dot {
  width: 5px; height: 5px;
  border-radius: 50%;
  background: rgba(34,211,238,0.3);
}
</style>
