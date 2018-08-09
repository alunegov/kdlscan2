# kdlscan2

[![Build Status](https://travis-ci.org/alunegov/kdlscan2.svg?branch=master)](https://travis-ci.org/alunegov/kdlscan2)

Утилита для синхронизации файлов локализации для [Kryvich's Delphi Localizer](https://sites.google.com/site/kryvich/localizer). В отличие от оригинальных kdlscan/lngupdate не использует коды ресурсов (resourcestring) во временных файлах и поддерживает добавление строк для псевдо-gettext режима.

## Установка

### Вручную

Загрузить исполняемый файл со страницы [Releases](https://github.com/alunegov/kdlscan2/releases).

### Используя go get

Требуются установленные [Go 1.10+](https://golang.org/) и [dep 0.5.0+](https://github.com/golang/dep).

```
go get -d github.com/alunegov/kdlscan2
cd %GOPATH%/src/github.com/alunegov/kdlscan2
dep ensure
go install github.com/alunegov/kdlscan2/cmd/kdlscan2
```

## Использование

```
0. kdlscan exe out_path
1. kdlscan2 scan proto_lng lng [псевдо-gettext]
2. kdlscan2 update edit_lng proto_lng [-!] [-x]
3. kdlscan2 generate lng edit_lng drc [drc_encoding]
4. kdlscan2 sync edit_lng ref_edit_lng
```
