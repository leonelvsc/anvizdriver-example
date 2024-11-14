package comandos

import (
	"log"
)

type SetDeviceInfo2 struct {
	Cmd Comando
	Data []byte
	Respuesta SetDeviceInfo2Respuesta
}

type SetDeviceInfo2Respuesta struct {
	RespuestaOriginal Respuesta
}

func (c *SetDeviceInfo2) Build() (comando, respuesta []byte) {
	c.Cmd = Comando{0x33, c.Data, 11 }
	c.Respuesta = SetDeviceInfo2Respuesta{}
	return c.Cmd.Build()
}

func (c *SetDeviceInfo2) Decode(buffer []byte) (err error) {
	decoded, cmdErr := c.Cmd.Decode(buffer)

	err = cmdErr
	if err != nil {
		log.Println("error set device 2 : ", err)
		return
	}

	c.Respuesta = SetDeviceInfo2Respuesta{RespuestaOriginal: decoded}
	return
}