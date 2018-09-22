# CHAT
_El chat hasta ahora solo funciona en la terminal_
## INSTALACIÓN
_Se clona el repositorio_
```
$ https://github.com/FERBP3/Proyecto-Chat.git
```
## En este punto hay dos formas de correr el programa
_Se redirecciona el GOPATH al directorio donde se clonó el reposotorio (La ruta es absoluta)_
```
$ export GOPATH=$RutaActual/
```
## Después se generan los binarios. Estando dentro de la carpeta src, escribimos:
```
$ go install Servidor/ Cliente/ Util/
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
_Para correr las prubas unitarias se ejecuta estando dentro de la carpeta src/_
```
$ go test Servidor/
$ go test Cliente/
```
## Autor
* **Brigido Pablo José Fernando**
* **No. Cuenta:314146171**
