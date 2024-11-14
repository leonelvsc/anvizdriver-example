package comandos

import (
	"log"
)

type GetRecordInfo struct {
	Cmd Comando
	Respuesta GetRecordInfoRespuesta
}

type GetRecordInfoRespuesta struct {
	RespuestaOriginal Respuesta
	UserAmount uint32
	FpAmount uint32
	PasswordAmount uint32
	CardAmount uint32
	AllRecordAmount uint32
	NewRecordAmount uint32
}

func (c *GetRecordInfo) Build() (comando, respuesta []byte) {
	c.Cmd = Comando{0x3c, make([]byte, 0), 29 }
	c.Respuesta = GetRecordInfoRespuesta{}
	return c.Cmd.Build()
}

func (c *GetRecordInfo) Decode(buffer []byte) (err error) {
	decoded, cmdErr := c.Cmd.Decode(buffer)

	err = cmdErr
	if err != nil {
		log.Println("error get record info: ", err)
		return
	}

	c.Respuesta = GetRecordInfoRespuesta{RespuestaOriginal: decoded}

	c.Respuesta.UserAmount = int24ToInt32(decoded.DATA[0:3])
	c.Respuesta.FpAmount = int24ToInt32(decoded.DATA[3:6])
	c.Respuesta.PasswordAmount = int24ToInt32(decoded.DATA[6:9])
	c.Respuesta.CardAmount = int24ToInt32(decoded.DATA[9:12])
	c.Respuesta.AllRecordAmount = int24ToInt32(decoded.DATA[12:15])
	c.Respuesta.NewRecordAmount = int24ToInt32(decoded.DATA[15:18])
	return
}