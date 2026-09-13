import { defineConfig } from 'vitepress';

export default defineConfig({
  lang: 'ru-RU',
  title: 'Мини-апки Telecom',
  description: 'Интеграция мини-апки в приложение Казахтелекома: контекст запуска, проверка подписи, методы моста.',
  // Served by the Go server from docs/.vitepress/dist at /; .html links stay
  // so a plain file server resolves them.
  cleanUrls: false,
  themeConfig: {
    nav: [
      { text: 'Руководство', link: '/' },
      { text: 'Методы моста', link: '/bridge/' },
      { text: 'Пример', link: 'https://github.com/Asylmurat99/tsa-miniapp-example' },
    ],
    sidebar: [
      {
        text: 'Начало',
        items: [
          { text: 'Быстрый старт', link: '/' },
          { text: 'Как это работает', link: '/how-it-works' },
          { text: 'Что вы присылаете и что получаете', link: '/onboarding' },
        ],
      },
      {
        text: 'Интеграция',
        items: [
          { text: 'Контекст запуска', link: '/launch-context' },
          { text: 'Проверка конверта', link: '/verify' },
          { text: 'Абонент и гость', link: '/customer-and-guest' },
          { text: 'Смена секрета', link: '/secret-rotation' },
        ],
      },
      {
        text: 'Методы моста',
        items: [
          { text: 'Обзор и поддержка', link: '/bridge/' },
          { text: 'getPhone', link: '/bridge/get-phone' },
          { text: 'storage', link: '/bridge/storage' },
          { text: 'Жизненный цикл', link: '/bridge/lifecycle' },
          { text: 'Системные действия', link: '/bridge/system' },
        ],
      },
      {
        text: 'Разработка',
        items: [
          { text: 'Отладка', link: '/debugging' },
          { text: 'Частые ошибки', link: '/troubleshooting' },
          { text: 'Версия контракта', link: '/changelog' },
        ],
      },
    ],
    outline: { label: 'На этой странице' },
    docFooter: { prev: 'Назад', next: 'Дальше' },
  },
});
