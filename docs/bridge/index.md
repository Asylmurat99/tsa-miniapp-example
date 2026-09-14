# Обзор и поддержка

Мини-апка обращается к приложению через SDK `@indigico/tsa-bridge` версии 2.3.0.

```bash
npm install @indigico/tsa-bridge@2.3.0
```

**Каждый вызов проходит проверку origin.** Мост отвечает только странице, загруженной с адреса, который мы вам зарегистрировали. Мини-апка, открытая с другого адреса, ответов от моста не получит.

Таблицы ниже сверены с кодом приложения. Метода нет ни в одной из них - значит, из мини-апки он не работает.

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
| [`getLanguage()`](/bridge/system) | Язык интерфейса приложения и список поддерживаемых |
| [`checkBiometry()`](/bridge/system) | Системный диалог биометрии. Результат - только для интерфейса |
| [`storage.setItem` / `getItem` / `clear`](/bridge/storage) | Локальное хранилище, своё у каждой мини-апки |
| [`getPhone()`](/bridge/get-phone) | Подписанный конверт с номером абонента. Требует права `phone:read` |

Нативный футер под мини-апкой есть не у всех - только там, где мы его дали; это не часть стандартного контейнера, и сама мини-апка им не управляет: [жизненный цикл](/bridge/lifecycle).

## Объявлено в SDK, в приложении не реализовано

Эти методы есть в типах SDK, но действие не выполняется. Часть отвечает пустым значением или строкой `"not_implemented"` - промис при этом разрешается, и успешным этот ответ не является. Часть отвечает ошибкой, то есть отклоняет промис: такие вызовы оборачивайте в `try`. Ни один из ответов ниже не подтверждает, что что-то произошло.

| Метод | Что отвечает |
|---|---|
| `getMe` | Пустой профиль: `{"name":"","lastname":"","id":"", ...}` |
| `getContacts` | `{"contacts":[],"sign":""}` |
| `getGeo` | Отказ `PERMISSION_DENIED`: координаты абонента мини-апкам не выдаются |
| `getInitData` | Обработчика нет, промис не разрешается никогда, в консоли `--getInitData-isUnknown`. Контекст лежит в `window.tsaWebApp.initData` |
| `setLanguage` | Отказ `PERMISSION_DENIED`: язык приложения мини-апка не меняет |
| `enablePrivateMessaging` / `disablePrivateMessaging` | `"success"`, ничего не происходит |
| `getQr` | `"not_implemented"` |
| `getSMSCode` | `"not_implemented"` |
| `getUserProfile` | `{"name":"","lastname":""}` |
| `selectContact` | `"not_implemented"` |
| `setTitle` | `"success"`, заголовок не меняется |
| `setHeaderMenuItems` | `"success"`, меню не появляется |
| `shareFile` / `shareImage` | `"not_implemented"` |
| `openSettings` | `"not_implemented"` |
| `vibrate` | `null` |
| `openPayment` | `"not_implemented"` |
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
| `closeApplication` | `"not_implemented"` |
| `openUserProfile` | `"not_implemented"` |
| `openMiniApp` | Ответа нет: канал обрабатывается только в главном окне приложения |

Модуль `auth` (`getState`, `getToken`, `requestLogin`) - внутренний канал приложения для его собственного кабинета. Из мини-апки каждый его метод отвечает `PERMISSION_DENIED`. Сессию мини-апка поднимает по [контексту запуска](/launch-context), а не через `auth`.

::: warning supports() не говорит, работает ли метод
`bridge.supports('getPhone')` возвращает `false`, хотя метод работает: он уходит по низкоуровневому каналу, а `supports` смотрит на список именованных функций. При этом `bridge.supports('getGeo')` возвращает `true`, хотя `getGeo` мини-апке отказывает.

Проверяйте `bridge.isSupported()` - это ответ на вопрос «мы внутри приложения». Исход конкретного вызова разбирайте по коду ошибки, как на странице [`getPhone`](/bridge/get-phone). Единственный источник того, что реализовано, - таблицы выше.
:::
