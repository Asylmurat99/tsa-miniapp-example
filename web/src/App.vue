<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { inShell, readInitData, requestPhone, stripFragment } from './tsa';

type Session = { auth: 'customer' | 'guest'; user_id: string | null; scope: string[] };

const session = ref<Session | null>(null);
const phone = ref<string | null>(null);
const message = ref('');

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
      ? `The app rejected the launch context: ${(e as Error).message}. Reopen the mini app.`
      : 'Open this page from the Telecom app to sign in.';
  }
});

async function getPhone() {
  message.value = '';
  if (!inShell()) {
    message.value = 'getPhone works only inside the Telecom app.';
    return;
  }
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
        message.value = `Backend rejected the envelope: ${(e as Error).message}`;
      }
      break;
    case 'unavailable':
      message.value = 'No confirmed number for this subscriber. Ask for it in your own form.';
      break;
    case 'denied':
      message.value = 'phone:read is not granted to this mini app, or the visitor is a guest.';
      break;
    case 'error':
      message.value = outcome.message;
  }
}
</script>

<template>
  <main>
    <h1>Sample MiniApp</h1>
    <section v-if="session">
      <p v-if="session.auth === 'customer'">Signed in as <code>{{ session.user_id }}</code></p>
      <p v-else>Guest visit. Nothing is known about the visitor.</p>
      <p>Scope: <code>{{ session.scope.join(', ') || '—' }}</code></p>
      <button v-if="session.auth === 'customer' && session.scope.includes('phone:read')" @click="getPhone">
        Get phone number
      </button>
      <p v-if="phone">Phone: <code>{{ phone }}</code></p>
    </section>
    <p v-if="message">{{ message }}</p>
  </main>
</template>
