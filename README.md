# Идентификация компонентов go

Инсталляция

``` sh
go install github.com/ilya2049/gocomponent/cmd/gocompvis
```

## Визуализация компонентов

``` sh
gocompvis -project-dir=/project/dir -root-namespace=internal > components.dot

dot -Tsvg -O components.dot
```