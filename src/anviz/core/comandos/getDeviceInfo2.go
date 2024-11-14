package comandos

import (
	"log"
)

type GetDeviceInfo2 struct {
	Cmd Comando
	Respuesta GetDeviceInfo2Respuesta
}

type GetDeviceInfo2Respuesta struct {
	RespuestaOriginal Respuesta
}

func (c *GetDeviceInfo2) Build() (comando, respuesta []byte) {
	c.Cmd = Comando{0x32, make([]byte, 0), 26 }
	c.Respuesta = GetDeviceInfo2Respuesta{}
	return c.Cmd.Build()
}

func (c *GetDeviceInfo2) Decode(buffer []byte) (err error) {
	decoded, cmdErr := c.Cmd.Decode(buffer)

	err = cmdErr
	if err != nil {
		log.Println("error get device info 2: ", err)
		return
	}

	c.Respuesta = GetDeviceInfo2Respuesta{RespuestaOriginal: decoded}
	return
}