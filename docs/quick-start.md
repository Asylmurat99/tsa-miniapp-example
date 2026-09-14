# Быстрый старт

Несколько команд поднимают этот сайт и рабочую мини-апку на вашей машине.

```bash
git clone https://github.com/Asylmurat99/tsa-miniapp-example
cd tsa-miniapp-example
npm ci && npm run build
cp .env.example .env
set -a && source .env && set +a
go run ./server
```

Откройте `http://localhost:8080/` - это тот же сайт, который вы читаете. Мини-апка живёт на `http://localhost:8080/app/`.

## Первый экран без приложения

Приложение открывает мини-апку с подписанным контекстом во фрагменте адреса. Пока у вас нет секрета, контекст можно подписать локально тестовым:

```bash
go run ./server -sign \
  -sign-user 3f1a9c40-77ad-4e0f-9c3d-11c2a0f5e881 \
  -sign-scope phone:read
```

```bash
go run ./server -sign \
  -sign-user 3f1a9c40-77ad-4e0f-9c3d-11c2a0f5e881 \
  -sign-scope phone:read \
  | python3 -c 'import sys, urllib.parse
print("http://localhost:8080/app/#tsaWebAppData="
  + urllib.parse.quote(sys.stdin.read().strip(), safe=""))'
```

Вторая команда сразу печатает готовый адрес - откройте его в браузере. Мини-апка отправит контекст своему серверу, тот проверит подпись и покажет, кто пришёл.

## Первый экран из приложения

1. Опубликуйте локальный сервер по https: `cloudflared tunnel --url http://localhost:8080`.
2. Сообщите менеджеру адрес `https://<subdomain>.trycloudflare.com/app/` как адрес мини-апки для разработки и получите `app_id` и секрет - [что вы присылаете и что получаете](/onboarding).
3. Запишите их в `.env`, перезапустите сервер.
4. Откройте мини-апку из приложения с dev-окружением - [отладка](/debugging).

::: warning Адрес туннеля меняется
Быстрый туннель получает новый адрес при каждом запуске. У нас адрес обновляется с задержкой до 60 секунд.
:::

## Что дальше

- [Как это работает](/how-it-works) - одна схема.
- [Контекст запуска](/launch-context) - что приходит и что с этим делать; файл `web/src/tsa.ts` целиком.
- [Проверка конверта](/verify) - алгоритм, тестовый вектор и код на пяти языках.
