// The only file that knows the page runs inside the Telecom app.
// Everything else talks to our own backend.
import bridge from '@indigico/tsa-bridge';
import type { GetPhoneResponse } from '@indigico/tsa-bridge';

const FRAGMENT_KEY = 'tsaWebAppData';

declare global {
  interface Window {
    tsaWebApp?: { initData?: string };
  }
}

/**
 * The launch context exactly as the app issued it. Prefer window.tsaWebApp:
 * the fragment is lost on navigation and reload. Empty string means the page
 * was opened outside the app.
 */
export function readInitData(): string {
  const fromHost = window.tsaWebApp?.initData;
  if (fromHost) return fromHost;
  return new URLSearchParams(location.hash.slice(1)).get(FRAGMENT_KEY) ?? '';
}

/**
 * Remove the context from the address bar, keep the rest of the fragment
 * (your own route). Filtered by hand: URLSearchParams.toString() would turn
 * "#/catalog" into "#/catalog=".
 */
export function stripFragment(): void {
  const rest = location.hash
    .slice(1)
    .split('&')
    .filter((part) => !part.startsWith(FRAGMENT_KEY + '='))
    .join('&');
  history.replaceState(null, '', location.pathname + location.search + (rest ? '#' + rest : ''));
}

/** True inside the Telecom app. In a desktop browser every bridge call fails. */
export function inShell(): boolean {
  return bridge.isSupported();
}

// #region phone
export type PhoneOutcome =
  | { kind: 'envelope'; envelope: GetPhoneResponse }
  | { kind: 'unavailable' } // subscriber has no confirmed number: ask for it manually
  | { kind: 'denied' } // phone:read not granted, or a guest
  | { kind: 'error'; message: string };

function errorCode(e: unknown): string | undefined {
  if (typeof e === 'object' && e !== null && 'code' in e) {
    return String((e as { code: unknown }).code);
  }
  return undefined;
}

/** Ask the app for the signed phone envelope. Hand the envelope to the backend whole. */
export async function requestPhone(): Promise<PhoneOutcome> {
  try {
    return { kind: 'envelope', envelope: await bridge.getPhone() };
  } catch (e) {
    switch (errorCode(e)) {
      case 'PHONE_UNAVAILABLE':
        return { kind: 'unavailable' };
      case 'PERMISSION_DENIED':
        return { kind: 'denied' };
      default:
        return { kind: 'error', message: e instanceof Error ? e.message : String(e) };
    }
  }
}
// #endregion phone
