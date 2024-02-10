# Идентификация компонентов go

Инсталляция

``` sh
go install github.com/ilya2049/gocomponent/cmd/gocompvis
```

## Визуализация компонентов

``` sh
gocompvis -project-dir=/project/dir > components.dot

dot -Tsvg -O components.dot
```

TODO:
- добавить цвета узлов [color=black, fillcolor=cyan, style=filled]
- добавить http сервер
- определять по go module одинаковые модули