# Идентификация компонентов go

## Визуализация компонентов

``` sh
gocompvis -project-dir=/project/dir -root-namespace=internal > components.dot

dot -Tsvg -O components.dot
```