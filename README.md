### Зачем?
Я очень привык в гитхабе использовать [hub](https://hub.github.com) чтобы создавать пулл реквесты в консоли.
Основная используемая команда `hub pull-request`.

Но после переезда основных репозиториев на гитлаб пришлось что-то придумывать.

### Как использовать?

- В профиле на гитлабе получить ключ для API и добавить его в переменную окружения `GITLAB_API_TOKEN`
- Собрать приложение (если Go установлен) go build или ./build.sh (он сразу скопирует бинарник в `/usr/local/opt/bin/lab`)

в папке с кодом выполнить `lab pull-request [title]`
Если `title` не задан, то программа спросит

### WARNING

Это пре-пре-пре альфа которая умеет делать только пулл реквесты, все остальные команды проксируются на системный гит
Программа поставляется as is, и если она вам что-то сломала, я не виноват. Умысла никакого нет, исходники все открыты :)

Contributions are always welcome!