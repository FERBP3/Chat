package main

import (
	"testing"
	//"fmt"
    //"io/ioutil"
    "log"
    "net"
    "bytes"
    "strings"
    "Util"
)

var servidor *TCPServidor

func init() {
	//Creamos el servidor
	servidor, err := NewServer("tcp", ":1234")
	if err != nil {
		log.Println("Error al iniciar el servidor")
		return
	}
	go servidor.Run()
}

func TestServidorRun(t *testing.T){
	conn, err := net.Dial("tcp",":1234")
	if err != nil {
		t.Error("No se pudo conectar con el servidor: ", err)
	}
	defer conn.Close()
}

func TestLeeMensaje(t *testing.T){
	TestLeeMensajeInvalido(t)
	TestLeeMensajeSinIdentificar(t)

}

func TestLeeMensajeInvalido(t *testing.T) {
	mensaje := []byte("Cualquier cosa\n")

	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		t.Error("No se pudo conectar al servidor: ", err)
	}
	defer conn.Close()
	if _, err := conn.Write(mensaje); err != nil {
		t.Error("No se pudo escribir al servidor")
	}
	salida := make([]byte,1024)
	n, err := conn.Read(salida)
	salida = salida[:n]
	if err != nil {
			t.Error("No se pudo leer del servidor ")
	}
	if bytes.Compare(salida, []byte(MensajeComandos)) != 0 {
		t.Error("El servidor no regresó lo esperado")
	}
}

func TestLeeMensajeSinIdentificar(t *testing.T) {
	mensajes := [...]string{
		"STATUS status\n",
		"USERS\n",
		"MESSAGE username messageContent\n",
		"PUBLICMESSAGE messageContent\n",
		"CREATEROOM roomname\n",
		"INVITE roomname user1\n",
		"JOINROOM roomname\n",
		"ROOMESSAGE roomname messageContent\n",
		"DISCONNECT\n",
	}

	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		t.Error("No se pudo conectar al servidor: ", err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(mensajes[0])); err != nil {
		t.Error("No se pudo escribir al servidor")
	}
	salida := make([]byte,1024)
	n, err := conn.Read(salida)
	if err != nil {
			t.Error("No se pudo leer del servidor ")
	}
	salida = salida[:n]
	comando := strings.Fields(strings.TrimSpace(string(salida)))
	if 	Util.EsComando(comando[0]){
		if bytes.Compare(salida, []byte(MensajeIdentidad)) != 0 {
			t.Error("El servidor no regresó lo esperado")
		}
	}

}


func TestMandaMensajes(t *testing.T){
	if false {
		t.Fail()
	}
}

func TestHandleConnection(t *testing.T){
	TestNumeroClientes(t)
}

func TestNumeroClientes(t *testing.T){
	if false {
		t.Fail()
	}
}