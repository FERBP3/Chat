# CHAT
_El chat hasta ahora solo funciona en la terminal_
## INSTALACIÓN
_Se obtiene con la ayuda del comando go get_
```
$ go get github.com/FERBP3/Chat

```
## Se compila y se generan los binarios del Cliente y el Servidor
_Estos estarán en $GOPATH/bin_
```
$ go install github.com/FERBP3/Proyecto-Chat/src/Cliente/
$ go install github.com/FERBP3/Proyecto-Chat/src/Servidor/
```

## Luego para levantar el servidor ejecutamos
```
$ $GOPATH/bin/Servidor :1234
```
_Como argumentos en la línea de comandos se debe pasar la ip y el puerto._
_Se puede no espeficar la ip poniendo solo el puerto como en el ejemplo de arriba (:1234)._
_Luego se ejecuta el cliente de la misma manera:_
```
$ $GOPATH/bin/Cliente :1234
```
## PRUEBAS UNITARIAS
_Para correr las pruebas unitarias se ejecuta_
```
$ go test $GOPATH/src/github.com/FERBP3/Proyecto-Chat/src/Servidor/
$ go test $GOPATH/src/github.com/FERBP3/Proyecto-Chat/src/Cliente/
```
## Autor
* **Brigido Pablo José Fernando**

