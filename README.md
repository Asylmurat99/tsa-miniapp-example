# tsa-miniapp-example

Документация и рабочий пример мини-апки для приложения Казахтелекома.

    git clone https://github.com/Asylmurat99/tsa-miniapp-example && cd tsa-miniapp-example
    npm ci && npm run build
    cp .env.example .env && set -a && source .env && set +a && go run ./server

После этого `http://localhost:8080/` - документация, `http://localhost:8080/app/` - пример мини-апки.
Документацию можно править на лету: `npm run docs:dev`.

## Что внутри

- `docs/` - сайт документации (VitePress). Начните с него.
- `web/` - мини-апка на Vue. `src/tsa.ts` - единственный файл, который знает о платформе.
- `server/` - сервер на Go без зависимостей: проверяет подпись контекста и конверта `getPhone`, поднимает сессию.
- `verify/` - проверка подписи на Go, Node, PHP, Python и Java на одном тестовом векторе. Эти же файлы показаны в документации.

## Проверить подпись без приложения

    go run ./server -sign -sign-user 3f1a9c40-77ad-4e0f-9c3d-11c2a0f5e881 -sign-scope phone:read

печатает свежеподписанный контекст. Откройте `http://localhost:8080/app/#tsaWebAppData=<контекст, закодированный процентами>`.
