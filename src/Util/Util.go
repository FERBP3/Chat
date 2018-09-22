package Util

type Comando int 
var Comandos = []string{
	"IDENTIFY", 
	"STATUS",
	"USERS",
	"MESSAGE",
	"PUBLICMESSAGE",
	"CREATEROOM",
	"INVITE",
	"JOINROOM",
	"ROOMESSAGE",
	"DISCONNECT",
	}

func EsComando(c string) bool {
	for _, comando := range Comandos {
		if c == comando {
			return true
		}
	}
	return false
}