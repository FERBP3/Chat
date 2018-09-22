package main 

import (
	"fmt"
	"net"
	"bufio"
	"os"
	"io"
	//"./Util"
)

func main(){
	direccion := os.Args[1]

	entrada := bufio.NewReader(os.Stdin)
	var nombreUsuario string

	fmt.Printf("Conectando a %s ...\n", direccion)
	conn, err := net.Dial("tcp", direccion)
	if err != nil {
		fmt.Println("Hubo un error en la conexión")
		os.Exit(1)
	}
	defer conn.Close()

	go leeMensaje(conn)
	mandaMensaje(conn, entrada, nombreUsuario)
}

func mandaMensaje(conn net.Conn, entrada *bufio.Reader, nombreUsuario string){
	for {
		buffer, _ := entrada.ReadString('\n')
		if len(buffer) > 0 {
			conn.Write([]byte(buffer))
		}
	}
}

func leeMensaje(conn net.Conn){
	var data []byte
	for {
		n, err := conn.Read(data)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Hubo un error al leer en la conexión")
				os.Exit(0)
				}
		}
		data = data[:n]
		fmt.Print(string(data))
		data = make([]byte, 256)
	}
}
