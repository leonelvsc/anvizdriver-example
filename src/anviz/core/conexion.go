package core

import (
	"net"
	"time"
	"log"
	"io"
	"anviz/core/comandos"
)

const TimeOut = 30 * time.Second

// Clase encargada de hablar con el aparato
// Entiende Comandos del paquete comandos
type Conexion struct {
	conn net.Conn
	err error
	cmdBytes []byte
	buffer []byte
}

func (c *Conexion) Establecer(direccion string) (err error) {
	c.conn, c.err = net.DialTimeout("tcp", direccion, TimeOut)
	if c.err != nil {
		err = c.err
	}
	return
}

func (c *Conexion) Cerrar() {
	c.conn.Close()
}

//Recibe un comando que sabe construirse y decodificar su respuesta del aparato
func (c *Conexion) EnviarComando(cmd comandos.IComando) (err error){
	c.cmdBytes, c.buffer = cmd.Build()

	timer := c.conn.SetReadDeadline(time.Now().Add(TimeOut))
	if c.err != nil {
		log.Println("SetReadDeadline failed:", timer)
		err = c.err
		// do something else, for example create new conn
		return
	}

	c.conn.Write(c.cmdBytes)

	//n, err := io.ReadAtLeast(c.conn, c.buffer, 11)
	//Si no usamos readfull sucede que corta el buffer antes de que termine la respuesta
	n, err := io.ReadFull(c.conn, c.buffer)

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Println("read timeout:", err)
			log.Println("readed:", n)
			// time out
		} else {
			log.Println("read error:", err)
			log.Println("readed:", n)
			// some error else, do something else, for example create new conn
		}

		return
	}

	err = cmd.Decode(c.buffer)
	return
}