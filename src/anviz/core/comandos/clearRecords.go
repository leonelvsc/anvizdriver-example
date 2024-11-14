package comandos

type ClearRecord struct {
	Cmd Comando
	Cantidad uint32
	Respuesta ClearRecordRespuesta
}

type ClearRecordRespuesta struct {
	RespuestaOriginal Respuesta
	UserAmount uint32
	FpAmount uint32
	PasswordAmount uint32
	CardAmount uint32
	AllRecordAmount uint32
	NewRecordAmount uint32
}

func (c *ClearRecord) Build() (comando, respuesta []byte) {
	c.Cmd = Comando{0x4E, append([]byte{2}, PutUint24(c.Cantidad)...), 14 }
	c.Respuesta = ClearRecordRespuesta{}
	return c.Cmd.Build()
}

func (c *ClearRecord) Decode(buffer []byte) (err error) {
	//decoded, cmdErr := c.Cmd.Decode(buffer)
	//
	//err = cmdErr
	//if err != nil {
	//	log.Println("error en clear records: ", err)
	//	return
	//}
	//
	//c.Respuesta = ClearRecordRespuesta{RespuestaOriginal: decoded}
	return
}