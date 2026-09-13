# Обзор и поддержка

Мини-апка обращается к приложению через SDK `@indigico/tsa-bridge` версии 2.3.0.

```bash
npm install @indigico/tsa-bridge@2.3.0
```

**Каждый вызов проходит проверку origin.** Мост отвечает только странице, загруженной с адреса, который мы вам зарегистрировали. Мини-апка, открытая с другого адреса, ответов от моста не получит.

## Работают

| Метод | Что делает |
|---|---|
| [`ready`](/bridge/lifecycle) | Снимает сплеш. Необязателен: сплеш уходит сам после загрузки, таймаут 10 секунд |
| [`closeMiniApp`](/bridge/lifecycle) | Закрывает мини-апку, с учётом подтверждения ниже |
| [`enableClosingConfirmation` / `disableClosingConfirmation`](/bridge/lifecycle) | Диалог подтверждения при закрытии |
| [`enableSwipeBack` / `disableSwipeBack`](/bridge/lifecycle) | Разрешить или запретить закрытие жестом вниз |
| [`openExternalUrl(url)`](/bridge/system) | Открыть ссылку во внешнем браузере |
| [`share(text)`](/bridge/system) | Системное меню «Поделиться» |
| [`copyToClipboard(text)`](/bridge/system) | Копирование в буфер |
| [`storage.setItem` / `getItem` / `clear`](/bridge/storage) | Локальное хранилище, своё у каждой мини-апки |
| [`getPhone()`](/bridge/get-phone) | Подписанный конверт с номером абонента. Требует права `phone:read` |

Высота, которую перекрывают нативный футер и системная полоса, лежит в CSS-переменной `--tsa-bottom-inset` и обновляется на лету — нижний отступ страницы считайте от неё. Самим футером мини-апка не управляет: [жизненный цикл](/bridge/lifecycle).

## Объявлено в SDK, в приложении не реализовано

Эти методы есть в типах SDK, но действие не выполняется. Часть отвечает пустым значением или строкой `"not_implemented"` — промис при этом разрешается, и успешным этот ответ не является. Часть отвечает ошибкой, то есть отклоняет промис: такие вызовы оборачивайте в `try`. Ни один из ответов ниже не подтверждает, что что-то произошло.

| Метод | Что отвечает |
|---|---|
| `getMe` | Пустой профиль: `{"name":"","lastname":"","id":"", …}` |
| `getContacts` | `{"contacts":[],"sign":""}` |
| `getGeo` | Ошибка `getGeo not implemented` |
| `getQr` | `"not_implemented"` |
| `getSMSCode` | `"not_implemented"` |
| `getUserProfile` | `{"name":"","lastname":""}` |
| `selectContact` | `"not_implemented"` |
| `setTitle` | `"success"`, заголовок не меняется |
| `setHeaderMenuItems` | `"success"`, меню не появляется |
| `shareFile` | `"not_implemented"` |
| `vibrate` | `null` |
| `openPayment` | `"not_implemented"` |
| `checkBiometry` | `"unavailable"` |
| `readNFCData` | Ошибка `readNFCData not implemented` |
| `readNFCPassport` | Ошибка `readNFCPassport not implemented` |
| `isESimSupported` | `"not_implemented"` |
| `activateESim` | `"not_implemented"` |
| `getUserStepInfo` | `{"steps":[]}` |
| `subscribeUserStepInfo` | `"not_implemented"` |
| `unsubscribeUserStepInfo` | `"not_implemented"` |
| `setNavigationItemMode` | `null` |
| `getNavigationItemMode` | `"NoItem"` |
| `setCustomBackArrowMode` | `"success"`, режим не меняется |
| `getCustomBackArrowMode` | `false` |
| `setCustomBackArrowVisible` | `"success"`, стрелка не появляется |
| `setTabBarVisible` | `"success"`, футер контейнера не меняется |
| `enableNotifications` | `{}` |
| `disableNotifications` | `{}` |
| `enableScreenCapture` | `null` |
| `disableScreenCapture` | `null` |
| `openSettings` | `"not_implemented"` |
| `closeApplication` | `"not_implemented"` |
| `openUserProfile` | `"not_implemented"` |
| `openMiniApp` | Ответа нет: канал обрабатывается только в главном окне приложения |

Модуль `auth` (`getState`, `getToken`, `requestLogin`) — внутренний канал приложения для его собственного кабинета. Из мини-апки каждый его метод отвечает `PERMISSION_DENIED`. Сессию мини-апка поднимает по [контексту запуска](/launch-context), а не через `auth`.

::: warning supports() не говорит, работает ли метод
`bridge.supports('getPhone')` возвращает `false`, хотя метод работает: он уходит по низкоуровневому каналу, а `supports` смотрит на список именованных функций. При этом `bridge.supports('getGeo')` возвращает `true`, хотя `getGeo` — заглушка.

Проверяйте `bridge.isSupported()` — это ответ на вопрос «мы внутри приложения». Исход конкретного вызова разбирайте по коду ошибки, как на странице [`getPhone`](/bridge/get-phone). Единственный источник того, что реализовано, — таблицы выше.
:::
