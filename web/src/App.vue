<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { lang, languages, setLang, t } from './i18n';
import { inShell, readInitData, requestPhone, stripFragment } from './tsa';

type Session = { auth: 'customer' | 'guest'; user_id: string | null; scope: string[] };
type MessageKey = 'rejected' | 'openFromApp' | 'onlyInsideApp' | 'backendRejected' | 'unavailable' | 'denied';
// A message is stored as a key plus the raw reason, so switching the language re-renders it.
type Message = { kind: 'hint' | 'error'; key: MessageKey; reason?: string } | { kind: 'error'; text: string };

const session = ref<Session | null>(null);
const phone = ref<string | null>(null);
const busy = ref(false);
const message = ref<Message | null>(null);

const isCustomer = computed(() => session.value?.auth === 'customer');
const canRequestPhone = computed(() => isCustomer.value && session.value!.scope.includes('phone:read'));

const messageText = computed(() => {
  const m = message.value;
  if (!m) return '';
  if ('text' in m) return m.text;
  const entry = t.value[m.key];
  return typeof entry === 'function' ? entry(m.reason ?? '') : entry;
});

async function api<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, init);
  const body = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(body.error ?? `HTTP ${res.status}`);
  return body as T;
}

onMounted(async () => {
  const initData = readInitData();
  stripFragment();
  try {
    if (initData) {
      session.value = await api<Session>('/api/session', {
        method: 'POST',
        headers: { 'Content-Type': 'text/plain' },
        body: initData,
      });
    } else {
      session.value = await api<Session>('/api/me');
    }
  } catch (e) {
    message.value = initData
      ? { kind: 'error', key: 'rejected', reason: (e as Error).message }
      : { kind: 'hint', key: 'openFromApp' };
  }
});

async function getPhone() {
  message.value = null;
  if (!inShell()) {
    message.value = { kind: 'hint', key: 'onlyInsideApp' };
    return;
  }
  busy.value = true;
  try {
    const outcome = await requestPhone();
    switch (outcome.kind) {
      case 'envelope':
        try {
          const res = await api<{ phone: string }>('/api/phone', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(outcome.envelope),
          });
          phone.value = res.phone;
        } catch (e) {
          message.value = { kind: 'error', key: 'backendRejected', reason: (e as Error).message };
        }
        break;
      case 'unavailable':
        message.value = { kind: 'hint', key: 'unavailable' };
        break;
      case 'denied':
        message.value = { kind: 'error', key: 'denied' };
        break;
      case 'error':
        message.value = { kind: 'error', text: outcome.message };
    }
  } finally {
    busy.value = false;
  }
}
</script>

<template>
  <div class="page">
    <header class="hero">
      <div class="hero-inner">
        <div class="brand">
          <span class="mark" aria-hidden="true">S</span>
          <div>
            <h1>{{ t.title }}</h1>
            <p class="tagline">{{ t.tagline }}</p>
          </div>
        </div>
        <div class="hero-side">
          <div class="lang-switch" role="group" aria-label="Language">
            <button
              v-for="l in languages"
              :key="l.code"
              type="button"
              :class="{ active: l.code === lang }"
              :aria-pressed="l.code === lang"
              @click="setLang(l.code)"
            >
              {{ l.label }}
            </button>
          </div>
          <span v-if="session" class="badge" :class="session.auth">
            {{ isCustomer ? t.subscriber : t.guest }}
          </span>
        </div>
      </div>
    </header>

    <main class="content">
      <section v-if="session" class="card visitor fade">
        <div class="avatar" :class="{ guest: !isCustomer }" aria-hidden="true">
          <svg viewBox="0 0 24 24"><circle cx="12" cy="8" r="4" /><path d="M4 20c0-4 3.6-6 8-6s8 2 8 6" /></svg>
        </div>
        <div class="visitor-text">
          <p class="visitor-title">{{ isCustomer ? t.signedIn : t.browsingAsGuest }}</p>
          <code v-if="isCustomer" class="pseudonym">{{ session.user_id }}</code>
          <p v-else class="muted">{{ t.noIdentity }}</p>
          <div class="chips">
            <span v-for="s in session.scope" :key="s" class="chip">
              <svg viewBox="0 0 16 16" aria-hidden="true"><path d="M3 8.5l3 3 7-7" /></svg>
              {{ s }}
            </span>
            <span v-if="session.scope.length === 0" class="chip empty">{{ t.noPermissions }}</span>
          </div>
        </div>
      </section>

      <section v-if="session" class="card action fade">
        <h2>{{ t.tryBridge }}</h2>
        <template v-if="canRequestPhone">
          <p class="muted">{{ t.askPhone }}</p>
          <button class="primary" :disabled="busy" @click="getPhone">
            <span v-if="busy" class="spinner" aria-hidden="true"></span>
            {{ busy ? t.waiting : t.getPhone }}
          </button>
          <Transition name="pop">
            <div v-if="phone" class="phone-card">
              <span class="phone-label">
                <svg viewBox="0 0 16 16" aria-hidden="true"><path d="M3 8.5l3 3 7-7" /></svg>
                {{ t.confirmed }}
              </span>
              <code class="phone">{{ phone }}</code>
            </div>
          </Transition>
        </template>
        <p v-else-if="isCustomer" class="muted">{{ t.needsGrant }}</p>
        <p v-else class="muted">{{ t.signInToUnlock }}</p>
      </section>

      <p v-if="message" class="notice fade" :class="message.kind">{{ messageText }}</p>

      <section class="card fade">
        <h2>{{ t.showsTitle }}</h2>
        <ul class="highlights">
          <li v-for="h in t.highlights" :key="h.title">
            <span class="dot" aria-hidden="true"></span>
            <div>
              <p class="highlight-title">{{ h.title }}</p>
              <p class="muted">{{ h.text }}</p>
            </div>
          </li>
        </ul>
      </section>

      <footer class="footer">
        {{ t.footer }} <a href="/" target="_blank" rel="noopener">{{ t.readGuide }}</a>
      </footer>
    </main>
  </div>
</template>

<style>
:root {
  --bg: #f2f4f8;
  --card: #ffffff;
  --text: #15181e;
  --muted: #667085;
  --line: #e4e7ec;
  --accent: #2f6df6;
  --accent-strong: #1d4fd8;
  --accent-text: #ffffff;
  --hero-from: #0f1f4d;
  --hero-to: #2f6df6;
  --ok-bg: #e8f7ee;
  --ok-text: #14733a;
  --error-bg: #fdecec;
  --error-text: #9f1d1d;
  --hint-bg: #eef1f5;
  --radius: 16px;
  --shadow: 0 8px 24px rgba(15, 31, 77, 0.08);
  color-scheme: light dark;
}

@media (prefers-color-scheme: dark) {
  :root {
    --bg: #0f1218;
    --card: #181c25;
    --text: #eef1f6;
    --muted: #98a2b3;
    --line: #262c38;
    --hero-from: #0b1533;
    --hero-to: #2456c9;
    --ok-bg: #12301d;
    --ok-text: #7ad39b;
    --error-bg: #3a1717;
    --error-text: #f4a6a6;
    --hint-bg: #222836;
    --shadow: 0 8px 24px rgba(0, 0, 0, 0.35);
  }
}

* {
  box-sizing: border-box;
}

html,
body {
  margin: 0;
  background: var(--bg);
  color: var(--text);
  font: 16px/1.45 -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
}

h1,
h2,
p {
  margin: 0;
}

.page {
  min-height: 100vh;
}

.hero {
  background: linear-gradient(135deg, var(--hero-from), var(--hero-to));
  color: #ffffff;
  padding: calc(20px + env(safe-area-inset-top)) 16px 56px;
}

.hero-inner {
  max-width: 480px;
  margin: 0 auto;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.mark {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: grid;
  place-items: center;
  background: rgba(255, 255, 255, 0.16);
  border: 1px solid rgba(255, 255, 255, 0.28);
  font-size: 22px;
  font-weight: 800;
  flex: none;
}

.hero h1 {
  font-size: 22px;
  font-weight: 700;
  line-height: 1.2;
}

.tagline {
  font-size: 13px;
  opacity: 0.8;
  margin-top: 2px;
}

.hero-side {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
  flex: none;
}

.lang-switch {
  display: inline-flex;
  padding: 3px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.14);
  border: 1px solid rgba(255, 255, 255, 0.24);
}

.lang-switch button {
  border: 0;
  background: transparent;
  color: rgba(255, 255, 255, 0.75);
  font: inherit;
  font-size: 12px;
  font-weight: 700;
  padding: 4px 9px;
  border-radius: 999px;
  cursor: pointer;
}

.lang-switch button.active {
  background: #ffffff;
  color: var(--hero-from);
}

.badge {
  flex: none;
  padding: 5px 12px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 600;
  background: rgba(255, 255, 255, 0.18);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.content {
  max-width: 480px;
  margin: -36px auto 0;
  padding: 0 16px calc(24px + env(safe-area-inset-bottom));
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.card {
  background: var(--card);
  border: 1px solid var(--line);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 16px;
}

.card h2 {
  font-size: 16px;
  font-weight: 700;
  margin-bottom: 8px;
}

.visitor {
  display: flex;
  gap: 14px;
  align-items: flex-start;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: var(--accent);
  color: var(--accent-text);
  font-size: 20px;
  font-weight: 700;
  flex: none;
}

.avatar.guest {
  background: var(--hint-bg);
  color: var(--muted);
}

.avatar svg {
  width: 24px;
  height: 24px;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.visitor-text {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.visitor-title {
  font-weight: 600;
}

.pseudonym {
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--muted);
  word-break: break-all;
}

.muted {
  color: var(--muted);
  font-size: 14px;
}

.chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 4px;
}

.chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 10px;
  border-radius: 999px;
  background: var(--ok-bg);
  color: var(--ok-text);
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 12px;
}

.chip.empty {
  background: var(--hint-bg);
  color: var(--muted);
  font-family: inherit;
}

.chip svg,
.phone-label svg {
  width: 12px;
  height: 12px;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.action .muted {
  margin-bottom: 12px;
}

button.primary {
  width: 100%;
  padding: 14px;
  border: 0;
  border-radius: 12px;
  background: var(--accent);
  color: var(--accent-text);
  font: inherit;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  box-shadow: 0 6px 16px rgba(47, 109, 246, 0.35);
  transition: transform 0.1s ease, background 0.15s ease;
}

button.primary:active {
  transform: scale(0.98);
  background: var(--accent-strong);
}

button.primary:disabled {
  opacity: 0.7;
  cursor: default;
  transform: none;
}

.spinner {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.4);
  border-top-color: #ffffff;
  animation: spin 0.8s linear infinite;
}

.phone-card {
  margin-top: 12px;
  padding: 12px 14px;
  border-radius: 12px;
  background: var(--ok-bg);
  color: var(--ok-text);
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.phone-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 600;
}

.phone {
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 24px;
  font-weight: 600;
  letter-spacing: 0.04em;
  color: var(--text);
}

.notice {
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 14px;
}

.notice.hint {
  background: var(--hint-bg);
  color: var(--text);
}

.notice.error {
  background: var(--error-bg);
  color: var(--error-text);
}

.highlights {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.highlights li {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--accent);
  margin-top: 6px;
  flex: none;
}

.highlight-title {
  font-weight: 600;
  font-size: 15px;
}

.footer {
  text-align: center;
  font-size: 13px;
  color: var(--muted);
}

.footer a {
  color: var(--accent);
  text-decoration: none;
  font-weight: 600;
}

.fade {
  animation: rise 0.35s ease both;
}

.pop-enter-active {
  animation: rise 0.3s ease both;
}

@keyframes rise {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: none;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (prefers-reduced-motion: reduce) {
  .fade,
  .pop-enter-active,
  .spinner {
    animation: none;
  }
}
</style>
