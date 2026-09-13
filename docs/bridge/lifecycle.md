# Жизненный цикл

Методы, которыми мини-апка управляет своим контейнером: сплешем, закрытием и нативным футером.

## ready

```ts
ready: () => Promise<ActionResponse>
```

Сообщает приложению, что мини-апка отрисовалась, и снимает сплеш. Необязателен: сплеш уходит сам после загрузки, таймаут 10 секунд. Вызывайте, когда первый экран готов, чтобы пользователь не смотрел на сплеш лишнее время.

```js
await bridge.ready();
```

## closeMiniApp

```ts
closeMiniApp: () => Promise<ActionResponse>
```

Закрывает мини-апку и возвращает пользователя в приложение. Если включено подтверждение закрытия, сначала показывается диалог.

```js
await bridge.closeMiniApp();
```

## enableClosingConfirmation / disableClosingConfirmation

```ts
enableClosingConfirmation:  () => Promise<ActionResponse>
disableClosingConfirmation: () => Promise<ActionResponse>
```

Включает и выключает диалог подтверждения при закрытии. Включайте его на экранах с незавершённым вводом и выключайте, когда терять нечего.

```js
await bridge.enableClosingConfirmation();
```

## enableSwipeBack / disableSwipeBack

```ts
enableSwipeBack:  () => Promise<ActionResponse>
disableSwipeBack: () => Promise<ActionResponse>
```

Разрешает и запрещает закрытие мини-апки жестом. Запрещайте на экранах, где жест конфликтует с вашими собственными свайпами — каруселью или списком со смахиванием.

```js
await bridge.disableSwipeBack();
```

## setTabBarVisible

```ts
setTabBarVisible: (visible: boolean) => Promise<ActionResponse>
```

Скрывает и показывает нативный футер контейнера. Приложение не возвращает футер само при переходах внутри страницы: уходя с экрана, который его скрыл, вызовите `setTabBarVisible(true)` явно.

```js
await bridge.setTabBarVisible(false);
```

Высота, которую перекрывает футер, лежит в CSS-переменной `--tsa-bottom-inset`: считайте нижний отступ страницы от неё, а не от нуля.
