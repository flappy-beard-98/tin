# tinkoff

Система для работы с рынком, использующая [tinkoff API](https://russianinvestments.github.io/investAPI/) как источник
данных. Анализирует данные рынка и строит прогнозы. Позволяет слать уведомления в телеграм бот.

## contents

- [analyzer](/analyzer/readme.md) - анализаторы, которые собранные данные
- [collector](/collector/readme.md) - коллекторы, которые собирают данные
- [core](/core/readme.md) - общая библиотека со вспомогательными типами и функциями
- [tgbot](/tgbot/readme.md) - телеграм бот
- [trader](/trader/readme.md) - поручения и их исполнение, то есть торговля

## todo

- [ ] задокументировать
- [x] получение аккаунтов
- [x] получение портфеля
- [x] получение котировок
- [ ] дивиденды
    - [x] загрузка данных о дивидендах
    - [x] поиск наилучшей последовательности покупка-получение-продажа-покупка для дивидендных акций
    - [ ] анализ котировок на дивидендные акции
        - [ ] загрузка котировок
        - [ ] найти минимальные котировки сразу после до даты возможности последней покупки
        - [ ] найти все котировки до даты возможности последней покупки
        - [ ] найти все котировки после даты выплат
    - [ ] телеграм бот, который будет слать уведомления о том что нужно купить-продать
- [ ] автоматический трейдинг
    - [ ] автоматические поручения на покупку-продажу по результатам работы
- [ ] добавить парсинг командной строки куда передавать токен
    - [x] токен можно читать из файла под .gitignore
- [ ] рефакторинг
    - [x] отказ от стейтфул объектов
    - [ ] использование [cobra](https://github.com/spf13/cobra) для создания cli
    - [ ] использование [viper](https://github.com/spf13/viper) для конфигов
    - [ ] запуск сценариев через шелл
- [ ] =/=
    
