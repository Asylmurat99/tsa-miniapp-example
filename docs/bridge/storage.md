# storage

Ключ-значение, своё пространство у каждой мини-апки: другая мини-апка ваших ключей не видит и не перезаписывает.

```ts
storage.setItem: (keyName: string, keyValue: string) => Promise<void>
storage.getItem: (keyName: string) => Promise<string | null>
storage.clear:   () => Promise<void>
```

- `setItem` - записывает значение. Значение - строка; объекты сериализуйте сами.
- `getItem` - читает значение. `null`, если ключа нет.
- `clear` - удаляет все ключи вашей мини-апки.

```js
await bridge.storage.setItem('theme', 'dark');
const theme = await bridge.storage.getItem('theme'); // 'dark'
```

::: warning Значения лежат в открытом виде
Хранилище не шифруется. Секреты, токены сессии и персональные данные туда не класть. Место для настроек интерфейса и черновиков, не для доступа.
:::
