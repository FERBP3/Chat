# CHAT
_ El chat hasta solo funciona en la terminal
## INSTALACIÓN
_ Se clona el repositorio
```
$ https://github.com/FERBP3/Proyecto-Chat.git
```
## En este punto hay dos formas de correr el programa
_ Se redirecciona el GOPATH al directorio donde se clonó el reposotorio
```
$ export $GOPATH=$RutaActual/
```
## Después se generan los binarios. Estando dentro de la carpeta src, escribimos:
```
$ go install Servidor/ Cliente/ Util/
```
## Luego para levantar el servidor ejecutamos
```
$ $GOPATH/bin/Servidor :1234
```
_ Como argumentos en la línea de comandos se debe pasar la ip y el puerto.
_ Se puede no espeficar la ip poniendo solo el puerto como en el ejemplo de arriba (:1234).
_ Luego se ejecuta el cliente de la misma manera:
```
$ $GOPATH/bin/Cliente :1234
```
