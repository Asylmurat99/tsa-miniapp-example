// Screen strings in three languages. The first visit follows the app language
// from the bridge; a choice made in the switcher lives in localStorage and
// wins afterwards. A dictionary is enough for a screen this size; swap it for
// vue-i18n when the mini app grows.
import { computed, ref } from 'vue';

export type Lang = 'ru' | 'kk' | 'en';

export const languages: { code: Lang; label: string }[] = [
  { code: 'ru', label: 'RU' },
  { code: 'kk', label: 'KZ' },
  { code: 'en', label: 'EN' },
];

const STORAGE_KEY = 'lang';
const DEFAULT_LANG: Lang = 'ru';

const strings = {
  ru: {
    title: 'Пример мини-приложения',
    tagline: 'Эталонная интеграция для мини-приложений Telecom',
    subscriber: 'Абонент',
    guest: 'Гость',
    signedIn: 'Вход через приложение Telecom',
    userIdHint: 'Постоянный идентификатор абонента для вашей мини-апки. Ключ учётной записи.',
    browsingAsGuest: 'Гостевой визит',
    noIdentity: 'Приложение не передало данные о посетителе.',
    noPermissions: 'без прав',
    tryBridge: 'Проверить мост',
    askPhone: 'Запросить у приложения подтверждённый номер телефона абонента.',
    getPhone: 'Получить номер телефона',
    waiting: 'Ждём ответ приложения',
    confirmed: 'Подтверждено Telecom',
    needsGrant: 'Для доступа к номеру нужно право phone:read. У этого мини-приложения его нет.',
    signInToUnlock: 'Войдите в приложение Telecom, чтобы делиться номером.',
    showsTitle: 'Что показывает пример',
    highlights: [
      {
        title: 'Проверенный контекст запуска',
        text: 'Приложение подписывает данные о том, кто открыл мини-приложение. Бэкенд проверяет подпись, прежде чем им доверять.',
      },
      {
        title: 'Гостевой режим',
        text: 'Посетители без аккаунта Telecom тоже попадают внутрь. Контекст просто не содержит личности.',
      },
      {
        title: 'Передача номера',
        text: 'С правом phone:read мини-приложение запрашивает у приложения подтверждённый номер. Это контактные данные, а не ключ: номер может измениться.',
      },
    ],
    footer: 'Работает на мосте Telecom.',
    readGuide: 'Открыть гайд',
    rejected: (reason: string) => `Приложение отклонило контекст запуска: ${reason}. Откройте мини-приложение заново.`,
    openFromApp: 'Откройте эту страницу из приложения Telecom, чтобы войти.',
    onlyInsideApp: 'getPhone работает только внутри приложения Telecom.',
    backendRejected: (reason: string) => `Бэкенд отклонил конверт: ${reason}`,
    unavailable: 'У абонента нет подтверждённого номера. Запросите его в своей форме.',
    denied: 'Право phone:read не выдано этому мини-приложению, либо посетитель гость.',
  },
  kk: {
    title: 'Мини-қосымша үлгісі',
    tagline: 'Telecom мини-қосымшаларына арналған эталондық интеграция',
    subscriber: 'Абонент',
    guest: 'Қонақ',
    signedIn: 'Telecom қосымшасы арқылы кірді',
    userIdHint: 'Абоненттің сіздің мини-қосымшаңызға арналған тұрақты идентификаторы. Есептік жазба кілті.',
    browsingAsGuest: 'Қонақ ретінде қарау',
    noIdentity: 'Қосымша келуші туралы дерек бермеді.',
    noPermissions: 'рұқсат жоқ',
    tryBridge: 'Көпірді тексеру',
    askPhone: 'Қосымшадан абоненттің расталған телефон нөмірін сұрау.',
    getPhone: 'Телефон нөмірін алу',
    waiting: 'Қосымшаның жауабын күтудеміз',
    confirmed: 'Telecom растады',
    needsGrant: 'Нөмірге қол жеткізу үшін phone:read рұқсаты қажет. Бұл мини-қосымшада ол жоқ.',
    signInToUnlock: 'Нөмірмен бөлісу үшін Telecom қосымшасына кіріңіз.',
    showsTitle: 'Үлгі нені көрсетеді',
    highlights: [
      {
        title: 'Тексерілген іске қосу контексті',
        text: 'Қосымша мини-қосымшаны кім ашқанын қолтаңбамен растайды. Бэкенд сенім артпас бұрын қолтаңбаны тексереді.',
      },
      {
        title: 'Қонақ режимі',
        text: 'Telecom аккаунты жоқ келушілер де кіре алады. Контексте жеке дерек болмайды.',
      },
      {
        title: 'Нөмірді беру',
        text: 'phone:read рұқсаты болса, мини-қосымша расталған нөмірді қосымшадан сұрайды. Бұл байланыс дерегі, кілт емес: нөмір өзгеруі мүмкін.',
      },
    ],
    footer: 'Telecom көпірінде жұмыс істейді.',
    readGuide: 'Нұсқаулықты ашу',
    rejected: (reason: string) => `Қосымша іске қосу контекстін қабылдамады: ${reason}. Мини-қосымшаны қайта ашыңыз.`,
    openFromApp: 'Кіру үшін бұл бетті Telecom қосымшасынан ашыңыз.',
    onlyInsideApp: 'getPhone тек Telecom қосымшасының ішінде жұмыс істейді.',
    backendRejected: (reason: string) => `Бэкенд конвертті қабылдамады: ${reason}`,
    unavailable: 'Абоненттің расталған нөмірі жоқ. Оны өз формаңызда сұраңыз.',
    denied: 'Бұл мини-қосымшаға phone:read рұқсаты берілмеген немесе келуші қонақ.',
  },
  en: {
    title: 'Sample MiniApp',
    tagline: 'Reference integration for Telecom mini apps',
    subscriber: 'Subscriber',
    guest: 'Guest',
    signedIn: 'Signed in via Telecom',
    userIdHint: 'Permanent subscriber id for your mini app. The account key.',
    browsingAsGuest: 'Browsing as guest',
    noIdentity: 'The app shared no identity for this visit.',
    noPermissions: 'no permissions',
    tryBridge: 'Try the bridge',
    askPhone: 'Ask the app for the confirmed phone number of this subscriber.',
    getPhone: 'Get phone number',
    waiting: 'Waiting for the app',
    confirmed: 'Confirmed by Telecom',
    needsGrant: 'Phone sharing needs the phone:read grant. This mini app does not have it.',
    signInToUnlock: 'Sign in to the Telecom app to unlock phone sharing.',
    showsTitle: 'What this sample shows',
    highlights: [
      {
        title: 'Verified launch context',
        text: 'The app signs who opened the mini app. The backend checks the signature before trusting it.',
      },
      {
        title: 'Guest mode',
        text: 'Visitors without a Telecom account still get in. The context just carries no identity.',
      },
      {
        title: 'Phone sharing',
        text: 'With the phone:read grant the mini app asks the app for the confirmed number. It is contact data, not a key: the number can change.',
      },
    ],
    footer: 'Built on the Telecom bridge.',
    readGuide: 'Read the guide',
    rejected: (reason: string) => `The app rejected the launch context: ${reason}. Reopen the mini app.`,
    openFromApp: 'Open this page from the Telecom app to sign in.',
    onlyInsideApp: 'getPhone works only inside the Telecom app.',
    backendRejected: (reason: string) => `Backend rejected the envelope: ${reason}`,
    unavailable: 'No confirmed number for this subscriber. Ask for it in your own form.',
    denied: 'phone:read is not granted to this mini app, or the visitor is a guest.',
  },
} as const;

export type Strings = (typeof strings)[Lang];

export function isLang(value: unknown): value is Lang {
  return value === 'ru' || value === 'kk' || value === 'en';
}

/** The language the visitor picked in the switcher, if any. */
export function storedLang(): Lang | null {
  try {
    const value = localStorage.getItem(STORAGE_KEY);
    if (isLang(value)) return value;
  } catch {
    // Storage may be blocked inside some web views; treat it as no choice.
  }
  return null;
}

export const lang = ref<Lang>(storedLang() ?? DEFAULT_LANG);
document.documentElement.lang = lang.value;

/** Switch the screen language. Only a choice made by the visitor is remembered. */
export function setLang(code: Lang, remember = true) {
  lang.value = code;
  document.documentElement.lang = code;
  if (!remember) return;
  try {
    localStorage.setItem(STORAGE_KEY, code);
  } catch {
    // Same as above: the choice just will not survive a reload.
  }
}

export const t = computed<Strings>(() => strings[lang.value]);
