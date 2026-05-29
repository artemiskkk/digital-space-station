<template>
  <div class="page">
    <!-- Top bar -->
    <header class="page-top">
      <div class="breadcrumb">
        <span class="crumb-dim">DSS</span>
        <span class="crumb-sep">/</span>
        <span>碎片记录</span>
      </div>
      <button class="btn-create" @click="showForm = !showForm">
        <span class="btn-icon">{{ showForm ? '×' : '+' }}</span>
        {{ showForm ? '取消' : '记一条' }}
      </button>
    </header>

    <!-- Compose -->
    <Transition name="slide">
      <div v-if="showForm" class="compose">
        <div class="compose-header">
          <span class="compose-label">NEW ENTRY</span>
          <span class="char-count">{{ form.text.length }}</span>
        </div>
        <div class="mood-bar">
          <span
            v-for="m in moodOptions" :key="m"
            :class="['mood', { active: form.mood === m }]"
            @click="form.mood = m"
          >{{ m }}</span>
        </div>
        <textarea
          v-model="form.text"
          placeholder="记录此刻..."
          rows="3"
          autofocus
        ></textarea>
        <div class="compose-foot">
          <button class="btn-publish" :disabled="saving || !form.text.trim()" @click="handleCreate">
            {{ saving ? 'TRANSMITTING...' : 'PUBLISH →' }}
          </button>
        </div>
      </div>
    </Transition>

    <!-- Stats row -->
    <div class="stat-row">
      <span class="stat-total">{{ total }} ENTRIES</span>
    </div>

    <!-- Timeline -->
    <div class="timeline">
      <div v-if="loading" class="empty-state"><div class="spinner"></div></div>
      <div v-else-if="!moments.length" class="empty-state">
        <p>NO ENTRIES LOGGED</p>
      </div>

      <div v-for="(m, i) in moments" :key="m.id" class="tl-item">
        <!-- Connector -->
        <div class="tl-line-col">
          <div class="tl-dot-ring">
            <span class="tl-mood">{{ m.mood || '💬' }}</span>
          </div>
          <div v-if="i < moments.length - 1" class="tl-connector"></div>
        </div>

        <!-- Card -->
        <div class="tl-card">
          <div class="tl-card-head">
            <time>{{ formatDate(m.created_at!) }}</time>
            <span class="tl-id">#{{ String(m.id).padStart(3, '0') }}</span>
            <button class="tl-del" @click="handleDelete(m.id!)">
              <svg width="11" height="11" viewBox="0 0 11 11" fill="none"><path d="M1.5 2.5h8M4 2.5V1.5h3v1M9 2.5l-.6 7H2.6L2 2.5" stroke="currentColor" stroke-width="1" stroke-linecap="round"/></svg>
            </button>
          </div>
          <p>{{ m.text }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getMoments, createMoment, deleteMoment, type Moment } from '../api'

const emit = defineEmits<{ 'update-counts': [counts: { moments: number }] }>()

const moments = ref<Moment[]>([])
const total = ref(0)
const loading = ref(false)
const saving = ref(false)
const showForm = ref(false)
const form = ref({ text: '', mood: '⚡' })
const moodOptions = ['⚡', '🛸', '💡', '🔧', '📡', '☕', '🌙', '📦', '🔒', '〰️']

async function load() {
  loading.value = true
  try {
    const res = await getMoments()
    moments.value = res.data.items || []
    total.value = res.data.meta?.total || 0
    emit('update-counts', { moments: total.value })
  } finally { loading.value = false }
}

async function handleCreate() {
  if (!form.value.text.trim()) return
  saving.value = true
  try {
    await createMoment({ text: form.value.text, mood: form.value.mood, images: [] })
    form.value = { text: '', mood: '⚡' }
    showForm.value = false
    await load()
  } finally { saving.value = false }
}

async function handleDelete(id: number) {
  if (!confirm('确认删除？')) return
  await deleteMoment(id)
  await load()
}

function formatDate(d: string) {
  return new Date(d).toLocaleString('zh-CN', {
    month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false
  })
}

onMounted(load)
</script>

<style scoped>
.page { position: relative; z-index: 1; }

/* ── Top bar ─────────────────── */
.page-top {
  display: flex; align-items: center; justify-content: space-between;
  padding: 20px 36px;
  border-bottom: 1px solid var(--border);
  background: rgba(8,12,20,0.6);
  backdrop-filter: blur(12px);
  position: sticky; top: 0; z-index: 10;
}
.breadcrumb {
  font-family: 'DM Mono', monospace;
  font-size: 12px; letter-spacing: 0.1em; color: var(--text-h);
}
.crumb-dim { color: var(--text2); }
.crumb-sep { color: var(--text2); margin: 0 6px; }
.btn-create {
  display: flex; align-items: center; gap: 6px;
  background: var(--accent); color: var(--bg);
  border: none; border-radius: 6px;
  padding: 7px 16px; font-size: 12px; font-weight: 600;
  cursor: pointer; transition: box-shadow 0.2s;
}
.btn-create:hover { box-shadow: 0 0 20px var(--accent-glow); }
.btn-icon { font-size: 15px; font-weight: 300; }

/* ── Compose ─────────────────── */
.compose {
  margin: 20px 36px;
  background: rgba(34,211,238,0.03);
  border: 1px solid rgba(34,211,238,0.1);
  border-radius: 10px;
  padding: 16px 18px;
  display: flex; flex-direction: column; gap: 10px;
}
.compose-header {
  display: flex; align-items: center; justify-content: space-between;
}
.compose-label {
  font-family: 'DM Mono', monospace;
  font-size: 9px; letter-spacing: 0.2em; color: var(--accent);
}
.char-count {
  font-family: 'DM Mono', monospace;
  font-size: 10px; color: var(--text2);
}
.mood-bar { display: flex; gap: 4px; }
.mood {
  font-size: 17px; padding: 3px 5px; border-radius: 5px;
  cursor: pointer; opacity: 0.3; transition: all 0.12s;
}
.mood:hover, .mood.active { opacity: 1; background: rgba(255,255,255,0.06); }
.compose textarea {
  background: rgba(255,255,255,0.03);
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 10px 14px;
  color: var(--text-h); font-size: 14px;
  resize: vertical; outline: none;
  line-height: 1.7;
}
.compose textarea:focus { border-color: rgba(34,211,238,0.25); }
.compose-foot { display: flex; justify-content: flex-end; }
.btn-publish {
  background: var(--accent); color: var(--bg);
  border: none; border-radius: 6px;
  padding: 7px 16px;
  font-family: 'DM Mono', monospace;
  font-size: 11px; font-weight: 600; letter-spacing: 0.1em;
  cursor: pointer; transition: box-shadow 0.2s;
}
.btn-publish:hover:not(:disabled) { box-shadow: 0 0 16px var(--accent-glow); }
.btn-publish:disabled { opacity: 0.4; cursor: not-allowed; }

/* ── Stats ─────────────────── */
.stat-row {
  padding: 14px 36px;
  border-bottom: 1px solid var(--border);
  font-family: 'DM Mono', monospace;
  font-size: 10px; letter-spacing: 0.15em; color: var(--text2);
}

/* ── Timeline ─────────────────── */
.timeline { padding: 24px 36px 60px; max-width: 700px; }
.tl-item { display: flex; gap: 16px; }
.tl-line-col {
  display: flex; flex-direction: column; align-items: center;
  flex-shrink: 0; width: 36px;
}
.tl-dot-ring {
  width: 36px; height: 36px;
  display: flex; align-items: center; justify-content: center;
  border-radius: 50%;
  background: var(--bg);
  border: 1px solid var(--border);
  font-size: 14px;
  flex-shrink: 0;
  transition: border-color 0.2s;
}
.tl-item:hover .tl-dot-ring {
  border-color: rgba(34,211,238,0.3);
  box-shadow: 0 0 10px rgba(34,211,238,0.1);
}
.tl-connector {
  flex: 1; width: 1px;
  background: linear-gradient(to bottom, rgba(34,211,238,0.15), var(--border));
  margin: 6px 0;
}

.tl-card {
  flex: 1; min-width: 0;
  background: rgba(255,255,255,0.015);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 14px 16px;
  margin-bottom: 16px;
  transition: all 0.15s;
}
.tl-card:hover { border-color: rgba(255,255,255,0.1); background: rgba(255,255,255,0.025); }

.tl-card-head {
  display: flex; align-items: center; gap: 8px; margin-bottom: 10px;
}
time {
  font-family: 'DM Mono', monospace;
  font-size: 11px; color: var(--text2); letter-spacing: 0.05em;
}
.tl-id {
  font-family: 'DM Mono', monospace;
  font-size: 10px; color: rgba(255,255,255,0.06);
}
.tl-del {
  margin-left: auto;
  background: none; border: none;
  color: var(--text2);
  padding: 4px; border-radius: 4px;
  cursor: pointer; opacity: 0; transition: all 0.15s;
}
.tl-card:hover .tl-del { opacity: 1; }
.tl-del:hover { color: var(--danger); background: rgba(248,81,73,0.08); }

.tl-card p { font-size: 14px; color: var(--text); line-height: 1.8; }

/* ── Shared ─────────────────── */
.empty-state {
  display: flex; flex-direction: column; align-items: center;
  padding: 80px 0; color: var(--text2);
  font-family: 'DM Mono', monospace; font-size: 12px; letter-spacing: 0.1em;
}
.spinner {
  width: 18px; height: 18px;
  border: 2px solid var(--accent-dim); border-top-color: var(--accent);
  border-radius: 50%; animation: spin .6s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

.slide-enter-active, .slide-leave-active { transition: all .2s ease; }
.slide-enter-from, .slide-leave-to { opacity: 0; transform: translateY(-6px); }
</style>
