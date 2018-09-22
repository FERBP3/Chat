package main 

import (
	"fmt"
	"net"
	"os"
	"log"
	"bufio"
	"strings"
	"errors"
	//"./Util"
)

var Comandos = [...]string{
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
const (
	

	MensajeComandos = "...INVALID MESSAGE\n"+
			"...VALID MENSSAGES ARE:\n"+
			"...IDENTIFY username\n"+
			"...STATUS userStatus = {ACTIVE, AWAY, BUSY}\n"+
			"...USERS\n"+
			"...MESSAGE username messageContent\n"+
			"...PUBLICMESSAGE messageContent\n"+
			"...CREATEROOM roomname\n"+
			"...INVITE roomname user1 user2 ...\n"+
			"...JOINROOM roomname\n"+
			"...ROOMESSAGE roomname messageContent\n"+
			"...DISCONNECT\n"
	MensajeIdentidad = "...DEBES IDENTIFICARTE PRIMERO\n"+
						"...PARA IDENTIFICARTE: IDENTIFY USERNAME\n"
	mensajeDisponible = "EL NOMBRE DE USUARIO NO ESTÁ DISPONIBLE\n"
	mensajeEstado = "...ESTADO INVALIDO\n"+
					 "...LOS ESTADOS PUEDEN SER: ACTIVE, AWAY, BUSY\n"

)

type Cliente struct{
	Nombre string 
	Conn net.Conn
	Estado string
}

type Sala struct {
	Nombre string
	Creador *Cliente
	Miembros []Cliente
	Invitados []Cliente
}

type TCPServidor struct {
	Direccion string
	Servidor net.Listener
	Clientes map[string]*Cliente
	NuevasConexiones chan net.Conn
	ConexionesCerradas chan net.Conn
	Mensajes chan string
	Salas map[string]*Sala
}

func NewServer(protocolo, direccion string) (*TCPServidor, error) {
	if protocolo == "tcp" {
		return &TCPServidor {
			Direccion: direccion,
			Clientes: make(map[string]*Cliente),
			NuevasConexiones: make(chan net.Conn),
			ConexionesCerradas: make(chan net.Conn),
			Mensajes: make(chan string),
			Salas: make(map[string]*Sala),
		}, nil
	}
	return nil, errors.New("EL PROTOCOLO ES INVÁLIDO")
}

func (t *TCPServidor) Run() (err error) {
	t.Servidor, err = net.Listen("tcp", t.Direccion)
	if err != nil {
		return err
	}
	defer t.Close()

	go t.AceptaConexiones()
	for {
		select {
		case conn := <- t.NuevasConexiones:
			go manejaConexion(conn, t.Clientes, t.Mensajes, t.Salas, t.ConexionesCerradas)
		case mensaje := <- t.Mensajes:
			for _, cliente := range t.Clientes {
				go mandaMensaje(cliente.Conn, mensaje, t.ConexionesCerradas)
			}
			log.Printf("Nuevo Mensaje: %s", mensaje)
			log.Printf("Transmitido a %d clientes", len(t.Clientes))

		case conn := <- t.ConexionesCerradas:
			usuario := buscaUsuario(conn, t.Clientes)
			delete(t.Clientes, usuario)
			conn.Close()
			log.Printf("Cliente %v desconectado", usuario)
		}
	}
	return
}


func (t *TCPServidor) Close() (err error) {
	return t.Servidor.Close()
}

func (t *TCPServidor) AceptaConexiones(){
	for {
		conn, err := t.Servidor.Accept()
		if err != nil {
			err = errors.New("No se pudo aceptar la conexion")
			break
		}
		if conn == nil {
			err = errors.New("No se pudo crear la conexion")
			break
		}
		t.NuevasConexiones <- conn
	}
}

func main(){
	direccion := os.Args[1]
	fmt.Println("Cargando servidor...")
	servidor, err := NewServer("tcp", direccion)
	if err != nil {
		fmt.Println("Error en la conexión del servidor")
		os.Exit(1)
	}
	fmt.Println("Servidor activo")
	servidor.Run()

}

func leeMensaje(conn net.Conn, mensajes chan string, 
				clientes map[string]*Cliente, salas map[string]*Sala, 
				conexionesCerradas chan net.Conn){
	lector := bufio.NewReader(conn)
	var nombreUsuario string
	var evento string
	var comando []string
	var nombreSala string

	// MensajeIdentidad := "...DEBES IDENTIFICARTE PRIMERO\n"+
	// 					"...PARA IDENTIFICARTE: IDENTIFY USERNAME\n"
	mensajeDisponible := "EL NOMBRE DE USUARIO NO ESTÁ DISPONIBLE\n"
	mensajeEstado := "...ESTADO INVALIDO\n"+
					 "...LOS ESTADOS PUEDEN SER: ACTIVE, AWAY, BUSY\n"

	for {
		entrada, err := lector.ReadString('\n')
		if err != nil {
			break
		}
		entrada = strings.TrimSpace(entrada)
		comando = strings.Fields(entrada)
		if len(comando) < 1 {
			mandaMensaje(conn, MensajeComandos, conexionesCerradas)
			continue
		}

		switch evento = comando[0]; evento {
		case "IDENTIFY":
			if len(comando) < 2{
				mandaMensaje(conn, MensajeComandos, conexionesCerradas)
				break
			}
			if nombreUsuario == "" {
				nombreUsuario = comando[1]
			_, ok := clientes[nombreUsuario]
			if ok {
				mandaMensaje(conn, mensajeDisponible, conexionesCerradas)
				nombreUsuario = ""
				break
			}else {
				mandaMensaje(conn,"Conectado!\n", conexionesCerradas)
				clientes[nombreUsuario] = &Cliente{
					Nombre: nombreUsuario,
					Conn: conn,
					Estado: "ACTIVE",
				}
	 		log.Printf("Nuevo Cliente aceptado %v. No. clientes: %d", nombreUsuario, len(clientes))
				break
			}
			}else {

			nuevoNombre := comando[1]
			_, ok := clientes[nuevoNombre]
			if ok {
				mandaMensaje(conn, mensajeDisponible, conexionesCerradas)
			}else {
				clientes[nuevoNombre] = &Cliente{
					Nombre: nuevoNombre,
					Conn: conn,
					Estado: "ACTIVE",
				}
				delete(clientes, nombreUsuario)
				mandaMensaje(conn,"...Cambiaste tu nombre a "+nuevoNombre+"\n", conexionesCerradas)
				nombreUsuario = nuevoNombre
			}
			}

		case "STATUS":
			if nombreUsuario == "" {
				mandaMensaje(conn, MensajeIdentidad, conexionesCerradas)
				break
			}
			if len(comando) < 2{
				mandaMensaje(conn, MensajeIdentidad, conexionesCerradas)
				break
			}
			estado := comando[1]
			if (estado == "ACTIVE") || (estado == "AWAY") || (estado == "BUSY") {
				clientes[nombreUsuario].Estado = estado
				break
			}
			mandaMensaje(conn, mensajeEstado, conexionesCerradas)

		case "USERS":
			if nombreUsuario == "" {
				mandaMensaje(conn, MensajeIdentidad, conexionesCerradas)
				break
			}
			mensaje := ""
			for _,usuario := range clientes{
				mensaje += "["+usuario.Estado+"]"+usuario.Nombre+"\n"
			}
			mandaMensaje(conn, mensaje, conexionesCerradas)

		case "PUBLICMESSAGE":
			if len(comando) < 2 {
				mandaMensaje(conn, MensajeComandos, conexionesCerradas)
				break
			}
			if nombreUsuario == "" {
				mandaMensaje(conn, MensajeIdentidad, conexionesCerradas)
				break
			}
			mensaje := "...PUBLIC-"+nombreUsuario+":"+strings.Join(comando[1:]," ")+"\n"
			mensajes <- mensaje

		case "MESSAGE":
			if len(comando) < 3 {
				mandaMensaje(conn, MensajeComandos, conexionesCerradas)
				break
			}
			if nombreUsuario == "" {
				mandaMensaje(conn, MensajeIdentidad, conexionesCerradas)
				break
			}
			receptor,ok := clientes[comando[1]]
			if !ok {
				mandaMensaje(conn, "El usuario "+comando[1]+" no existe.\n", conexionesCerradas)
				break
			}
			mensaje := "..."+nombreUsuario+":"+strings.Join(comando[2:]," ")+"\n"
			mandaMensaje(receptor.Conn, mensaje, conexionesCerradas)

		case "CREATEROOM":
			if len(comando) < 2 {
				mandaMensaje(conn, MensajeComandos, conexionesCerradas)
				break
			}
			if nombreUsuario == "" {
				mandaMensaje(conn, MensajeIdentidad, conexionesCerradas)
				break
			}
			nombreSala = comando[1]
			salas[nombreSala] = &Sala{ 
				Nombre: nombreSala, 
				Creador: clientes[nombreUsuario], 
				Miembros: make([]Cliente,0), 
				Invitados: make([]Cliente,0), }
			salas[nombreSala].Miembros = append(salas[nombreSala].Miembros, *clientes[nombreUsuario])
			mandaMensaje(conn, "...Se creó la sala "+nombreSala+"\n", conexionesCerradas)

		case "INVITE":
			if len(comando) < 3 {
				mandaMensaje(conn, MensajeComandos, conexionesCerradas)
				break
			}
			if nombreUsuario == "" {
				mandaMensaje(conn, MensajeIdentidad, conexionesCerradas)
				break
			}
			nombreSala = comando[1]
			sala,ok := salas[nombreSala]	
			if !ok {
				mandaMensaje(conn, "...Esa sala no existe.\n", conexionesCerradas)
				break
			} 
			if clientes[nombreUsuario] != sala.Creador {
				mandaMensaje(conn, "...TÚ NO ERES EL PROPIETARIO DE LA SALA\n", conexionesCerradas)
				break
			}

			for _, invitado := range comando[2:]{
				cliente, ok := clientes[invitado]
				if !ok {
					mandaMensaje(conn, "...EL USUARIO "+invitado+" NO SE ENCONTRÓ\n", conexionesCerradas)
				}else {
					mensajeInvitacion := "...Invitación para unirse a la sala "+
										 sala.Nombre+" de "+nombreUsuario+"\n"
					mandaMensaje(cliente.Conn, mensajeInvitacion, conexionesCerradas)
					mandaMensaje(conn, "...Invitación enviada a "+invitado+"\n", conexionesCerradas)
					sala.Invitados = append(sala.Invitados, *cliente)
				}
			}

		case "ROOMESSAGE":
			if len(comando) < 3 {
				mandaMensaje(conn, MensajeIdentidad, conexionesCerradas)
				break
			}
			if nombreUsuario == "" {
				mandaMensaje(conn, MensajeIdentidad, conexionesCerradas)
				break
			}
			nombreSala = comando[1]
			sala,ok := salas[nombreSala]	
			if !ok {
				mandaMensaje(conn, "....ESA SALA NO EXISTE\n", conexionesCerradas)
				break
			}
			if !contiene(nombreUsuario, sala.Miembros) {
				mandaMensaje(conn, "...TÚ NO ERES PARTE DE ESTA SALA\n", conexionesCerradas)
				break
			}
			mensaje := "..."+nombreSala+"-"+nombreUsuario+":"+strings.Join(comando[2:], " ")+"\n"
			for _, miembro := range sala.Miembros {
				mandaMensaje(miembro.Conn, mensaje, conexionesCerradas)
			}

		case "JOINROOM":
			if len(comando) < 2 {
				mandaMensaje(conn, MensajeIdentidad, conexionesCerradas)
				break
			}
			if nombreUsuario == "" {
				mandaMensaje(conn, MensajeIdentidad, conexionesCerradas)
				break
			}
			nombreSala := comando[1]
			sala,ok := salas[nombreSala]	
			if !ok {
				mandaMensaje(conn, "...ESA SALA NO EXISTE\n", conexionesCerradas)
				break
			}
			if !contiene(nombreUsuario, sala.Invitados) {
				mandaMensaje(conn, "...NO ERES INVITADO EN ESTA SALA\n", conexionesCerradas)
				break
			}
			if contiene(nombreUsuario, sala.Miembros) {
				mandaMensaje(conn, "...YA ERES PARTE DE ESTA SALA\n", conexionesCerradas)
				break
			}
			sala.Miembros = append(sala.Miembros, *clientes[nombreUsuario])
			mandaMensaje(conn, "...AHORA ERES PARTE DE LA SALA\n", conexionesCerradas)

		case "DISCONNECT":
			conexionesCerradas <- conn
		default:
			mandaMensaje(conn, MensajeComandos, conexionesCerradas)
		}
	}
}

func mandaMensaje(conn net.Conn, mensaje string, conexionesCerradas chan net.Conn){
	_, err := conn.Write([]byte(mensaje))
	if err != nil {
		conexionesCerradas <- conn
	}
}

func contiene(nombre string, miembros []Cliente) bool{
	for _, miembro := range miembros {
		if miembro.Nombre == nombre {
			return true
		}
	}
	return false
}

func buscaUsuario(conn net.Conn, clientes map[string]*Cliente) string {
	for nombre, cliente := range clientes{
		if conn == cliente.Conn{
			return nombre
		}
	}
	return ""
}

func manejaConexion(conn net.Conn, clientes map[string]*Cliente, mensajes chan string, 
					salas map[string]*Sala, conexionesCerradas chan net.Conn){
	leeMensaje(conn, mensajes, clientes, salas, conexionesCerradas)
}
