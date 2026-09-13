# Системные действия

Три метода, которыми мини-апка отдаёт данные наружу - в системное меню, в буфер обмена и во внешний браузер.

## share

```ts
share: (text: string) => Promise<ActionResponse>
```

Открывает системное меню «Поделиться» с текстом. Куда именно уйдёт текст, выбирает пользователь - приложение этого не знает и вам не сообщает.

```js
await bridge.share('Мой заказ №1024: https://miniapp.example.kz/orders/1024');
```

## copyToClipboard

```ts
copyToClipboard: (text: string) => Promise<ActionResponse>
```

Кладёт текст в системный буфер обмена. Системного уведомления о копировании нет - показывайте своё.

```js
await bridge.copyToClipboard('KZ123456789012345678');
```

## openExternalUrl

```ts
openExternalUrl: (url: string) => Promise<ActionResponse>
```

Открывает ссылку во внешнем браузере, за пределами контейнера мини-апки. Мини-апка при этом остаётся открытой.

```js
await bridge.openExternalUrl('https://example.kz/terms');
```

Контекст запуска во внешний браузер не переносится: страница, открытая так, о пользователе ничего не знает.
