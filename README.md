# Go CLI Calculator

Go CLI calculator totally written from scratch using basic Golang packages. 

* Basic operations
* Expression checking and logging
* Totally written in GO from scratch

| Operation |  Symbol  |
|:-----|:--------:|
| Add   | _+_ |
| Substraction   |  _-_  |
| Multiplication   | _*_ |
| Division   |  _/_  |
| Exponentiation   | _^_ |
| Root   | _$_ |

## Examples
```console
test@mac goCalculator % go run main.go "(((50 + 23) - 234) ^ 2)"          
25921
```

```console
test@mac goCalculator % go run main.go "(((((234 * 5) - 234) ^2) / 100) $ 3) * 2.234.21"
undefined token at position 41
```

```console
test@mac goCalculator % go run main.go "((((234 * 5) - 234) ^2) / 0)"
division by zero
```

```console
test@mac goCalculator % go run main.go "((((234.90 * 5) - 234) ^2) / UNDEFINED SYMBOL)"
undefined token at position 30
```

## Lessons learned

* Basic Golang syntax
* Basic Golang packages
* Shutting yard algorithm
