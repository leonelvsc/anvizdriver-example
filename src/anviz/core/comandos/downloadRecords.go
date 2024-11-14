package comandos

import (
	"log"
	"encoding/binary"
	"time"
)

const (
	Entrada = 0
 	Salida = 1
 	Break = 2

	Huella1 = 1
	Huella2 = 2
	Pass = 4
	Tarjeta = 8

	MASK_ATTENDANCE_STATE = 0x0F // 0000 1111 para usar con un & ( AND )
	MASK_CAN_OPEN_DOOR = 0x40 // 0100 0000
)
const TimeStampAdder = 946684800 + 10800 + 86400

type DownloadRecords struct {
	Cmd Comando
	Respuesta DownloadRecordsRespuesta
	RespuestaLongitud uint
	Data []byte
}

type DownloadRecordsRespuesta struct {
	RespuestaOriginal Respuesta
	ValidRecords int8
	Records []Record
}

type Record struct {
	UserCode uint64 `json:"user_code"`
	DateTime time.Time `json:"date_time"`
	BackUpCode int8 `json:"backup_code"`
	AttendanceMode int8 `json:"attendance_mode"`
	CanOpenDoor int8 `json:"can_open_door"`
	WorkTypes uint32 `json:"work_types"`
	Tipo string `json:"tipo"`
	ComoMarco string `json:"como_marco"`
}

func (c *DownloadRecords) Build() (comando, respuesta []byte) {
	c.RespuestaLongitud = uint(12 + int(c.Data[1]) * 14)
	c.Cmd = Comando{0x40, c.Data, c.RespuestaLongitud }
	c.Respuesta = DownloadRecordsRespuesta{}
	return c.Cmd.Build()
}

func (c *DownloadRecords) Decode(buffer []byte) (err error) {
	decoded, cmdErr := c.Cmd.Decode(buffer)

	err = cmdErr
	if err != nil {
		log.Println("error download records: ", err)
		return
	}

	c.Respuesta = DownloadRecordsRespuesta{RespuestaOriginal: decoded}
	c.Respuesta.ValidRecords = int8(decoded.DATA[0])

	i := 1

	// Por cada registro lo vamos agregando al slice de records de la respuesta
	// Vamos saltenado de 14 en 14 bytes .. registro a registro
	for i <= int(c.Respuesta.ValidRecords) {
		record := decoded.DATA[ 14*(i-1) + 1 : 14*i + 1 ]

		fichada := &Record{
			UserCode: int40ToInt64(record[0:5]),
			DateTime: toDate(record[5:9]),
			BackUpCode: int8(record[9]),
			AttendanceMode: int8(record[10] & MASK_ATTENDANCE_STATE),
			CanOpenDoor: int8(record[10] & MASK_CAN_OPEN_DOOR),
			WorkTypes: int24ToInt32(record[11:14]),
		}

		switch int8(record[10] & MASK_ATTENDANCE_STATE) {
		case Salida:
			fichada.Tipo = "Salida"
		case Entrada:
			fichada.Tipo = "Entrada"
		case Break:
			fichada.Tipo = "Break"
		}

		switch int8(record[9]) {
		case Huella1:
			fichada.ComoMarco = "Huella1"
		case Huella2:
			fichada.ComoMarco = "Huella2"
		case Pass:
			fichada.ComoMarco = "Pass"
		case Tarjeta:
			fichada.ComoMarco = "Tarjeta"
		}

		c.Respuesta.Records = append(c.Respuesta.Records, *fichada)

		i++
	}
	return
}

func toDate(bytes []byte) time.Time {
	return time.Unix(int64(binary.BigEndian.Uint32(bytes) + TimeStampAdder), 0)
}