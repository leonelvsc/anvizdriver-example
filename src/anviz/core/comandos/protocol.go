package comandos

import (
	"encoding/binary"
	"errors"
)

const (
	ACK_SUCCESS = 0x00
	ACK_FAIL = 0x01
	ACK_FULL = 0x04
	ACK_EMPTY = 0x05
	ACK_NO_USER = 0x06
	ACK_TIME_OUT = 0x08
	ACK_USER_OCCUPIED = 0x0A
	ACK_FINGER_OCUPPIED = 0x0B
)

type IComando interface {
	Build() (comando, respuesta []byte)
	Decode(buffer []byte) (err error)
}

type Comando struct {
	// Que comando es ej 0x3c
	cmd byte
	// Los datos del comando de 0-400 bytes
	data []byte

	respuestaLongitud uint
}

type Respuesta struct {
	STX byte
	CH uint32
	ACK byte
	RET byte
	LEN uint16
	DATA []byte
	CRC16 []byte
}

// Devuelve el array de bytes que representa al comando + el array de buffer para la respuesta
func (c *Comando) Build() (comando, respuesta []byte) {
	aux := []byte{0xa5, 0x00, 0x00, 0x00, 0x00, c.cmd}

	auxBytes := make([]byte, 2)

	dataLen := len(c.data)

	binary.BigEndian.PutUint16(auxBytes, uint16(dataLen))

	aux = append(aux, auxBytes...)

	if dataLen > 0 {
		aux = append(aux, c.data...)
	}

	Checksum(&aux)

	return aux, make([]byte, c.respuestaLongitud)
}

// Agarra la respuesta la decodifica y devuelve el array de datos con su longitud listo para parsear por el comando
func (c *Comando) Decode(buffer []byte) (r Respuesta, err error) {
	check, longitudBuffer := checkCrc16(buffer)

	if !check {
		err = errors.New("Error el checksum del paquete no es valido")
		return
	}

	longitudData := binary.BigEndian.Uint16(buffer[7:9])

	r = Respuesta{
		STX: buffer[0],
		CH: binary.BigEndian.Uint32(buffer[1:5]),
		ACK: buffer[5],
		RET: buffer[6],
		LEN: longitudData,
		DATA: buffer[9:9+longitudData],
		CRC16: buffer[longitudBuffer-2:longitudBuffer],
	}

	if r.RET != ACK_SUCCESS {
		var msg string

		switch r.RET {
		case ACK_FAIL:
			msg = "El dispositivo informa que no pudo realizar la operacion"
		case ACK_FULL:
			msg = "User full"
		case ACK_EMPTY:
			msg = "User empty"
		case ACK_NO_USER:
			msg = "El usuario no existe"
		case ACK_TIME_OUT:
			msg = "TimeOut"
		case ACK_USER_OCCUPIED:
			msg = "El usuario ya existe"
		case ACK_FINGER_OCUPPIED:
			msg = "La huella ya existe"
		}

		err = errors.New(msg)
	}

	return
}

func int24ToInt32(bytes []byte) uint32 {
	return uint32(bytes[2]) | uint32(bytes[1])<<8 | uint32(bytes[0])<<16
}

func int40ToInt64(bytes []byte) uint64 {
	return uint64(bytes[4]) | uint64(bytes[3])<<8 | uint64(bytes[2])<<16 | uint64(bytes[1])<<24 | uint64(bytes[0])<<32
}

func PutUint24(num uint32) (bytes []byte) {
	aux := make([]byte, 4)
	binary.BigEndian.PutUint32(aux, num)
	bytes = aux[1:4]
	return
}

// Recalcula el checksum para ver que los datos enviados esten bien
func checkCrc16(buffer []byte) (check bool, longitud int){
	longitud = len(buffer)

	crcEnviado := buffer[longitud-2:longitud]

	datosSinCrc := make([]byte, 0)

	datosSinCrc = append(datosSinCrc, buffer[0:longitud-2]...)

	Checksum(&datosSinCrc)

	check = binary.LittleEndian.Uint16(datosSinCrc[longitud-2:longitud]) == binary.LittleEndian.Uint16(crcEnviado)
	return
}