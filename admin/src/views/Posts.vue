<template>
  <div class="page">
    <!-- Top bar -->
    <header class="page-top">
      <div class="breadcrumb">
        <span class="crumb-dim">DSS</span>
        <span class="crumb-sep">/</span>
        <span>博客文章</span>
      </div>
      <button class="btn-create" @click="openCreate">
        <span class="btn-icon">+</span>
        写文章
      </button>
    </header>

    <!-- Stats -->
    <div class="stats">
      <div class="stat-card">
        <div class="stat-value">{{ posts.length }}</div>
        <div class="stat-label">TOTAL</div>
      </div>
      <div class="stat-card">
        <div class="stat-value accent">{{ posts.filter(p => p.status === 'published').length }}</div>
        <div class="stat-label">PUBLISHED</div>
      </div>
      <div class="stat-card">
        <div class="stat-value dim">{{ posts.filter(p => p.status === 'draft').length }}</div>
        <div class="stat-label">DRAFTS</div>
      </div>
    </div>

    <!-- Table -->
    <div class="table-wrap">
      <div class="table-head">
        <span class="th-status"></span>
        <span class="th-title">TITLE</span>
        <span class="th-tags">TAGS</span>
        <span class="th-date">DATE</span>
        <span class="th-actions"></span>
      </div>

      <div v-if="loading" class="empty-state">
        <div class="spinner"></div>
      </div>
      <div v-else-if="!posts.length" class="empty-state">
        <p>NO RECORDS FOUND</p>
        <p class="empty-hint">点击右上角「写文章」开始创作</p>
      </div>

      <TransitionGroup name="list" tag="div">
        <div
          v-for="p in posts"
          :key="p.id"
          class="table-row"
          @click="openEdit(p)"
        >
          <span class="row-status">
            <span :class="['dot', p.status]"></span>
          </span>
          <div class="row-title">
            <span class="title-text">{{ p.title }}</span>
            <span class="title-slug">/{{ p.slug }}</span>
          </div>
          <div class="row-tags">
            <span v-for="tag in (p.tags || []).slice(0, 3)" :key="tag" class="tag">{{ tag }}</span>
          </div>
          <span class="row-date">{{ formatDate(p.created_at!) }}</span>
          <div class="row-actions" @click.stop>
            <button class="act-btn" @click="openEdit(p)" title="编辑">
              <svg width="13" height="13" viewBox="0 0 13 13" fill="none"><path d="M9.5 1.5l2 2-7 7H2.5v-2l7-7Z" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round"/></svg>
            </button>
            <button class="act-btn danger" @click="handleDelete(p.id!)" title="删除">
              <svg width="13" height="13" viewBox="0 0 13 13" fill="none"><path d="M2 3.5h9M5 3.5V2h3v1.5M10.5 3.5l-.7 7.5H3.2L2.5 3.5" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round"/></svg>
            </button>
          </div>
        </div>
      </TransitionGroup>
    </div>

    <!-- Editor Modal -->
    <Transition name="modal">
      <div v-if="modal" class="modal-mask" @click.self="modal = false">
        <div class="editor">
          <!-- Editor top bar -->
          <div class="editor-bar">
            <div class="editor-bar-left">
              <span class="editor-dot"></span>
              <span>{{ editing ? 'EDIT POST' : 'NEW POST' }}</span>
            </div>
            <button class="editor-close" @click="modal = false">ESC</button>
          </div>

          <div class="editor-body">
            <div class="editor-fields">
              <div class="field-group">
                <div class="field flex-1">
                  <label>TITLE</label>
                  <input v-model="form.title" placeholder="文章标题" @input="autoSlug" />
                </div>
                <div class="field w-140">
                  <label>STATUS</label>
                  <select v-model="form.status">
                    <option value="draft">Draft</option>
                    <option value="published">Published</option>
                  </select>
                </div>
              </div>

              <div class="field-group">
                <div class="field flex-1">
                  <label>SLUG</label>
                  <input v-model="form.slug" placeholder="url-slug" />
                </div>
                <div class="field w-140">
                  <label>READ TIME</label>
                  <input v-model="form.read_time" placeholder="5 min" />
                </div>
              </div>

              <div class="field">
                <label>EXCERPT</label>
                <input v-model="form.excerpt" placeholder="一句话描述..." />
              </div>

              <div class="field">
                <label>TAGS</label>
                <input v-model="tagsInput" placeholder="嵌入式, GD32, PLC" />
              </div>
            </div>

            <div class="editor-content">
              <label>CONTENT · MARKDOWN</label>
              <textarea v-model="form.content" placeholder="# 在这里写 Markdown..."></textarea>
            </div>
          </div>

          <div class="editor-footer">
            <button class="btn-cancel" @click="modal = false">取消</button>
            <button class="btn-save" :disabled="saving || !form.title" @click="handleSave">
              <span v-if="saving" class="spinner-sm"></span>
              {{ saving ? '保存中...' : '保存' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getPosts, createPost, updatePost, deletePost, type Post } from '../api'

const emit = defineEmits<{ 'update-counts': [counts: { posts: number }] }>()

const posts = ref<Post[]>([])
const loading = ref(false)
const saving = ref(false)
const modal = ref(false)
const editing = ref<Post | null>(null)
const tagsInput = ref('')

const emptyForm = (): Post => ({
  title: '', slug: '', excerpt: '', content: '',
  cover_url: '', tags: [], status: 'draft', read_time: '',
})
const form = ref(emptyForm())

async function load() {
  loading.value = true
  try {
    const res = await getPosts()
    posts.value = res.data.items || []
    emit('update-counts', { posts: posts.value.length })
  } finally {
    loading.value = false
  }
}

function autoSlug() {
  if (!editing.value) {
    form.value.slug = form.value.title
      .toLowerCase().replace(/[\s一-鿿]+/g, '-')
      .replace(/[^a-z0-9-]/g, '').replace(/-+/g, '-').replace(/^-|-$/g, '')
  }
}

function openCreate() {
  editing.value = null
  form.value = emptyForm()
  tagsInput.value = ''
  modal.value = true
}
function openEdit(p: Post) {
  editing.value = p
  form.value = { ...p }
  tagsInput.value = p.tags?.join(', ') || ''
  modal.value = true
}

async function handleSave() {
  form.value.tags = tagsInput.value.split(',').map(t => t.trim()).filter(Boolean)
  saving.value = true
  try {
    if (editing.value?.id) await updatePost(editing.value.id, form.value)
    else await createPost(form.value)
    modal.value = false
    await load()
  } finally { saving.value = false }
}

async function handleDelete(id: number) {
  if (!confirm('确认删除？')) return
  await deletePost(id)
  await load()
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
}

onMounted(load)
</script>

<style scoped>
.page { padding: 0; position: relative; z-index: 1; }

/* ── Top bar ─────────────────────────── */
.page-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 36px;
  border-bottom: 1px solid var(--border);
  background: rgba(8,12,20,0.6);
  backdrop-filter: blur(12px);
  position: sticky; top: 0; z-index: 10;
}
.breadcrumb {
  font-family: 'DM Mono', monospace;
  font-size: 12px;
  letter-spacing: 0.1em;
  color: var(--text-h);
}
.crumb-dim { color: var(--text2); }
.crumb-sep { color: var(--text2); margin: 0 6px; }

.btn-create {
  display: flex; align-items: center; gap: 6px;
  background: var(--accent);
  color: var(--bg);
  border: none;
  border-radius: 6px;
  padding: 7px 16px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: box-shadow 0.2s;
}
.btn-create:hover { box-shadow: 0 0 20px var(--accent-glow); }
.btn-icon { font-size: 15px; font-weight: 300; }

/* ── Stats ─────────────────────────── */
.stats {
  display: flex;
  gap: 1px;
  background: var(--border);
  border-bottom: 1px solid var(--border);
}
.stat-card {
  flex: 1;
  background: var(--surface);
  padding: 20px 28px;
}
.stat-value {
  font-family: 'Orbitron', sans-serif;
  font-size: 28px;
  font-weight: 700;
  color: var(--text-h);
  margin-bottom: 4px;
}
.stat-value.accent { color: var(--accent); }
.stat-value.dim { color: var(--text2); }
.stat-label {
  font-family: 'DM Mono', monospace;
  font-size: 9px;
  letter-spacing: 0.2em;
  color: var(--text2);
}

/* ── Table ─────────────────────────── */
.table-wrap { padding: 0; }
.table-head {
  display: flex;
  align-items: center;
  padding: 10px 36px;
  font-family: 'DM Mono', monospace;
  font-size: 9px;
  letter-spacing: 0.18em;
  color: rgba(255,255,255,0.08);
  border-bottom: 1px solid var(--border);
  background: rgba(255,255,255,0.01);
}
.th-status { width: 24px; }
.th-title { flex: 1; }
.th-tags { width: 200px; }
.th-date { width: 70px; }
.th-actions { width: 70px; }

.table-row {
  display: flex;
  align-items: center;
  padding: 14px 36px;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  transition: background 0.12s;
}
.table-row:hover { background: rgba(255,255,255,0.02); }

.row-status { width: 24px; display: flex; align-items: center; }
.dot { width: 7px; height: 7px; border-radius: 50%; }
.dot.published { background: var(--accent); box-shadow: 0 0 6px var(--accent-glow); }
.dot.draft { background: var(--text2); }

.row-title { flex: 1; min-width: 0; }
.title-text {
  display: block;
  font-size: 14px;
  color: var(--text-h);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.title-slug {
  font-family: 'DM Mono', monospace;
  font-size: 11px;
  color: var(--text2);
}

.row-tags { width: 200px; display: flex; gap: 4px; flex-wrap: wrap; }
.tag {
  font-size: 10px;
  color: var(--text);
  background: rgba(255,255,255,0.04);
  padding: 2px 8px;
  border-radius: 3px;
  border: 1px solid var(--border);
}

.row-date {
  width: 70px;
  font-family: 'DM Mono', monospace;
  font-size: 11px;
  color: var(--text2);
}

.row-actions { width: 70px; display: flex; gap: 4px; justify-content: flex-end; opacity: 0; transition: opacity 0.15s; }
.table-row:hover .row-actions { opacity: 1; }
.act-btn {
  background: none;
  border: 1px solid var(--border);
  border-radius: 5px;
  color: var(--text2);
  width: 28px; height: 28px;
  display: flex; align-items: center; justify-content: center;
  cursor: pointer;
  transition: all 0.12s;
}
.act-btn:hover { color: var(--accent); border-color: rgba(34,211,238,0.25); background: var(--accent-dim); }
.act-btn.danger:hover { color: var(--danger); border-color: rgba(248,81,73,0.25); background: rgba(248,81,73,0.06); }

.empty-state {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  padding: 80px 0; color: var(--text2); font-family: 'DM Mono', monospace; font-size: 12px; letter-spacing: 0.1em; gap: 6px;
}
.empty-hint { font-size: 12px; color: var(--text2); letter-spacing: 0; font-family: 'Inter', sans-serif; }
.spinner { width: 20px; height: 20px; border: 2px solid var(--accent-dim); border-top-color: var(--accent); border-radius: 50%; animation: spin .6s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* ── Editor Modal ─────────────────── */
.modal-mask {
  position: fixed; inset: 0;
  background: rgba(1,4,9,0.88);
  backdrop-filter: blur(6px);
  display: flex; align-items: center; justify-content: center;
  z-index: 200; padding: 24px;
}
.editor {
  width: 760px; max-width: 100%;
  max-height: 92vh;
  background: var(--surface2);
  border: 1px solid var(--border2);
  border-radius: 12px;
  display: flex; flex-direction: column;
  box-shadow: 0 0 0 1px rgba(34,211,238,0.04), 0 40px 100px rgba(0,0,0,0.6);
  overflow: hidden;
}
.editor-bar {
  display: flex; align-items: center; justify-content: space-between;
  padding: 14px 20px;
  border-bottom: 1px solid var(--border);
  font-family: 'DM Mono', monospace;
  font-size: 10px; letter-spacing: 0.15em; color: var(--text2);
}
.editor-bar-left { display: flex; align-items: center; gap: 8px; }
.editor-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--accent); }
.editor-close {
  background: rgba(255,255,255,0.04);
  border: 1px solid var(--border);
  border-radius: 4px;
  color: var(--text2);
  font-family: 'DM Mono', monospace;
  font-size: 10px; letter-spacing: 0.1em;
  padding: 3px 10px; cursor: pointer;
  transition: all 0.12s;
}
.editor-close:hover { color: var(--text-h); border-color: var(--border2); }

.editor-body {
  flex: 1; overflow-y: auto; padding: 20px;
  display: flex; flex-direction: column; gap: 14px;
}
.editor-fields { display: flex; flex-direction: column; gap: 12px; }
.field-group { display: flex; gap: 12px; }
.flex-1 { flex: 1; }
.w-140 { width: 140px; flex-shrink: 0; }
.field { display: flex; flex-direction: column; gap: 5px; }
.field label {
  font-family: 'DM Mono', monospace;
  font-size: 9px; letter-spacing: 0.18em; color: var(--text2);
}
.field input, .field select, .editor-content textarea {
  background: rgba(255,255,255,0.03);
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 9px 12px;
  color: var(--text-h);
  font-size: 13px;
  outline: none;
  transition: border-color 0.15s;
}
.field input:focus, .editor-content textarea:focus { border-color: rgba(34,211,238,0.3); }
.field select option { background: var(--surface2); }

.editor-content { flex: 1; display: flex; flex-direction: column; gap: 5px; }
.editor-content label {
  font-family: 'DM Mono', monospace;
  font-size: 9px; letter-spacing: 0.18em; color: var(--text2);
}
.editor-content textarea {
  flex: 1; min-height: 200px;
  resize: vertical;
  font-family: 'DM Mono', monospace;
  font-size: 13px;
  line-height: 1.8;
}

.editor-footer {
  display: flex; justify-content: flex-end; gap: 8px;
  padding: 14px 20px;
  border-top: 1px solid var(--border);
}
.btn-cancel {
  background: none;
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text);
  padding: 8px 18px; font-size: 12px;
  cursor: pointer; transition: all 0.12s;
}
.btn-cancel:hover { color: var(--text-h); border-color: var(--border2); }
.btn-save {
  display: flex; align-items: center; gap: 6px;
  background: var(--accent);
  color: var(--bg);
  border: none; border-radius: 6px;
  padding: 8px 20px; font-size: 12px; font-weight: 600;
  cursor: pointer; transition: box-shadow 0.2s;
}
.btn-save:hover:not(:disabled) { box-shadow: 0 0 18px var(--accent-glow); }
.btn-save:disabled { opacity: 0.4; cursor: not-allowed; }
.spinner-sm { width: 12px; height: 12px; border: 2px solid rgba(8,12,20,.3); border-top-color: var(--bg); border-radius: 50%; animation: spin .6s linear infinite; }

/* Transitions */
.modal-enter-active, .modal-leave-active { transition: opacity .15s; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
.list-enter-active, .list-leave-active { transition: all .2s ease; }
.list-enter-from { opacity: 0; transform: translateX(-8px); }
.list-leave-to { opacity: 0; }
</style>
