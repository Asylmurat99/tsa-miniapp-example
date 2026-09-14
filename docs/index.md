---
layout: home

hero:
  name: Мини-апки Telecom
  text: Ваш сервис внутри приложения Казахтелекома
  tagline: Контекст запуска, проверка подписи, методы моста. Рабочий пример на Go и Vue, который можно поднять за пять минут.
  actions:
    - theme: brand
      text: Быстрый старт
      link: /quick-start
    - theme: alt
      text: Методы моста
      link: /bridge/
    - theme: alt
      text: Пример на GitHub
      link: https://github.com/Asylmurat99/tsa-miniapp-example

features:
  - title: Контекст запуска
    details: Приложение открывает вашу страницу с подписанным контекстом. Один файл разбирает его, ваш бэкенд проверяет подпись.
    link: /launch-context
    linkText: Как проверить
  - title: Абонент и гость
    details: Мини-апка работает и без входа в приложение. Контекст честно говорит, кто пришёл и что о нём известно.
    link: /customer-and-guest
    linkText: Два режима
  - title: Номер телефона
    details: С правом phone:read вы получаете подтверждённый номер абонента через мост, без собственной формы и SMS. Ключ учётной записи при этом user.id, а не номер.
    link: /bridge/get-phone
    linkText: getPhone
  - title: Проверка на пяти языках
    details: Один тестовый вектор и одинаковые реализации на Go, Node, Python, PHP и Java. Сверяйте свою с эталоном.
    link: /verify
    linkText: Проверка конверта
  - title: Отладка без приложения
    details: Подпишите контекст локально и откройте мини-апку в обычном браузере. Приложение и постоянный адрес нужны только в конце.
    link: /debugging
    linkText: Отладка
  - title: Частые ошибки
    details: invalid_signature, expired, app_mismatch. Что значит каждый отказ и с чего начать разбор.
    link: /troubleshooting
    linkText: Разобрать
---
