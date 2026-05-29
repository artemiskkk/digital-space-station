<template>
  <div class="layout">
    <!-- Sidebar -->
    <aside class="sidebar">
      <!-- Logo -->
      <a href="/" class="logo-link">
        <svg width="20" height="20" viewBox="0 0 22 22" fill="none">
          <circle cx="11" cy="11" r="9.5" stroke="rgba(255,255,255,0.12)" stroke-width="0.8"/>
          <circle cx="11" cy="11" r="5.5" stroke="rgba(255,255,255,0.06)" stroke-width="0.8" stroke-dasharray="2 3"/>
          <circle cx="11" cy="11" r="2" fill="#22d3ee"/>
          <ellipse cx="11" cy="11" rx="9.5" ry="3.5" stroke="#22d3ee" stroke-width="0.8" stroke-opacity="0.45" transform="rotate(-30 11 11)"/>
        </svg>
        <div class="logo-text">
          <span class="logo-name">DSS</span>
          <span class="logo-sub">COMMAND CENTER</span>
        </div>
      </a>

      <!-- Nav -->
      <div class="nav-section">
        <span class="nav-label">CONTENT</span>
        <nav>
          <RouterLink to="/posts" class="nav-item">
            <div class="nav-icon">
              <svg width="15" height="15" viewBox="0 0 15 15" fill="none"><rect x="2" y="2.5" width="11" height="1.4" rx=".7" fill="currentColor" opacity=".9"/><rect x="2" y="6" width="7.5" height="1.4" rx=".7" fill="currentColor" opacity=".5"/><rect x="2" y="9.5" width="9.5" height="1.4" rx=".7" fill="currentColor" opacity=".3"/></svg>
            </div>
            <span>博客文章</span>
            <span class="nav-count">{{ postCount }}</span>
          </RouterLink>
          <RouterLink to="/moments" class="nav-item">
            <div class="nav-icon">
              <svg width="15" height="15" viewBox="0 0 15 15" fill="none"><circle cx="7.5" cy="7.5" r="5.5" stroke="currentColor" stroke-width="1.1" opacity=".7"/><path d="M7.5 4.5v3.5l2 1.5" stroke="currentColor" stroke-width="1.1" stroke-linecap="round" stroke-linejoin="round"/></svg>
            </div>
            <span>碎片记录</span>
            <span class="nav-count">{{ momentCount }}</span>
          </RouterLink>
        </nav>
      </div>

      <!-- Footer -->
      <div class="sidebar-bottom">
        <div class="sys-status">
          <span class="sys-dot"></span>
          <span>ALL SYSTEMS NOMINAL</span>
        </div>
      </div>
    </aside>

    <!-- Main -->
    <main class="main">
      <div class="main-noise" aria-hidden="true"></div>
      <RouterView @update-counts="updateCounts" />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, RouterView } from 'vue-router'

const postCount = ref(0)
const momentCount = ref(0)

function updateCounts(counts: { posts?: number; moments?: number }) {
  if (counts.posts !== undefined) postCount.value = counts.posts
  if (counts.moments !== undefined) momentCount.value = counts.moments
}
</script>

<style scoped>
.layout {
  display: flex;
  min-height: 100vh;
}

/* ── Sidebar ─────────────────────────────────── */
.sidebar {
  width: 230px;
  flex-shrink: 0;
  background: rgba(1,4,9,0.95);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  padding: 0;
  position: fixed;
  top: 0; left: 0; bottom: 0;
  z-index: 50;
}

.logo-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 20px 18px;
  border-bottom: 1px solid var(--border);
  transition: background 0.2s;
}
.logo-link:hover { background: rgba(255,255,255,0.02); }
.logo-text { display: flex; flex-direction: column; gap: 1px; }
.logo-name {
  font-family: 'Orbitron', sans-serif;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.15em;
  color: var(--text-h);
}
.logo-sub {
  font-family: 'DM Mono', monospace;
  font-size: 8px;
  letter-spacing: 0.18em;
  color: var(--text2);
}

.nav-section { flex: 1; padding: 20px 10px 10px; }
.nav-label {
  font-family: 'DM Mono', monospace;
  font-size: 9px;
  letter-spacing: 0.2em;
  color: rgba(255,255,255,0.08);
  padding: 0 10px;
  display: block;
  margin-bottom: 10px;
}
nav { display: flex; flex-direction: column; gap: 2px; }

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 10px;
  border-radius: 7px;
  font-size: 13px;
  color: var(--text2);
  transition: all 0.15s;
  border: 1px solid transparent;
}
.nav-item:hover { color: var(--text); background: rgba(255,255,255,0.03); }
.nav-item.router-link-active {
  color: var(--accent);
  background: var(--accent-dim);
  border-color: rgba(34,211,238,0.08);
}
.nav-icon { width: 18px; display: flex; align-items: center; justify-content: center; }
.nav-count {
  margin-left: auto;
  font-family: 'DM Mono', monospace;
  font-size: 10px;
  color: var(--text2);
  background: rgba(255,255,255,0.04);
  padding: 1px 7px;
  border-radius: 10px;
}

.sidebar-bottom {
  padding: 14px;
  border-top: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.sys-status {
  display: flex;
  align-items: center;
  gap: 7px;
  font-family: 'DM Mono', monospace;
  font-size: 9px;
  letter-spacing: 0.12em;
  color: rgba(255,255,255,0.08);
  padding: 0 6px;
}
.sys-dot {
  width: 5px; height: 5px;
  border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 6px var(--accent-glow);
  animation: pulse 2.5s infinite;
}
@keyframes pulse { 0%,100%{opacity:1} 50%{opacity:.25} }


/* ── Main ─────────────────────────────────── */
.main {
  flex: 1;
  margin-left: 230px;
  position: relative;
  min-height: 100vh;
}
.main-noise {
  position: fixed;
  top: 0; left: 230px; right: 0; bottom: 0;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='200' height='200'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='200' height='200' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E");
  opacity: 0.015;
  pointer-events: none;
  z-index: 0;
}
</style>
