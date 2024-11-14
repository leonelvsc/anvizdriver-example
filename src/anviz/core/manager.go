package core

import (
	"math"
	"anviz/core/comandos"
	"errors"
)

const CantidadPorLote = 25

type Manager struct {
	conn Conexion
}

func (m *Manager) Conectar(direccion string) (err error) {
	m.conn = Conexion{}

	err = m.conn.Establecer(direccion)
	return
}

func (m *Manager) ObtenerRegistros() (records []comandos.Record, err error) {
	defer m.conn.Cerrar()

	deviceInfo := &comandos.GetDeviceInfo2{}
	recordInfo := &comandos.GetRecordInfo{}

	cmdErr := m.conn.EnviarComando(deviceInfo)
	if cmdErr != nil {
		err = errors.New("No es posible obtener la configuracion del dispositivo")
		return
	}

	// Si el modo en tiempo real esta activado
	if deviceInfo.Respuesta.RespuestaOriginal.DATA[4] == 1 {
		dataAux := deviceInfo.Respuesta.RespuestaOriginal.DATA
		dataAux[4] = 0
		setInfo := &comandos.SetDeviceInfo2{Data: dataAux}

		cmdErr = m.conn.EnviarComando(setInfo)
		if cmdErr != nil {
			err = errors.New("No es posible actualizar la configuracion del dispositivo")
			return
		}
	}


	cmdErr = m.conn.EnviarComando(recordInfo)
	if cmdErr != nil {
		err = errors.New("No es posible obtener la informacion del dispositivo")
		return
	}

	iteraciones := int(math.Ceil(float64(recordInfo.Respuesta.NewRecordAmount) / float64(CantidadPorLote)))

	restantes := int(recordInfo.Respuesta.NewRecordAmount)

	// Si no hay nada que leer
	if restantes == 0 {
		return
	}

	i := 0

	var cantidadDescargar int

	if restantes < 25 {
		cantidadDescargar = restantes
	} else {
		cantidadDescargar = CantidadPorLote
	}

	// Primer byte en 2 indica solo nuevos, 1 indica todos
	cmd := &comandos.DownloadRecords{Data: []byte{2, uint8(cantidadDescargar)}}

	for i < iteraciones {
		cmdErr := m.conn.EnviarComando(cmd)

		// Intentamos una vez mas para ver si falla
		if cmdErr != nil {
			cmd = &comandos.DownloadRecords{Data: []byte{0x10, uint8(cantidadDescargar)}}
			cmdErr = m.conn.EnviarComando(cmd)

			if cmdErr != nil {
				err = errors.New("No es posible realizar la descarga, intente mas tarde")
				return
			}
		}

		// La unica forma de que pase esto es que haya mas de una iteracion, cada iteracion descuenta 25 obviamente
		restantes = restantes - CantidadPorLote

		if restantes < 25 {
			cantidadDescargar = restantes
		} else {
			cantidadDescargar = CantidadPorLote
		}

		records = append(records, cmd.Respuesta.Records...)

		i++

		if i < iteraciones {
			cmd = &comandos.DownloadRecords{Data: []byte{0, uint8(cantidadDescargar)}}
		}
	}

	clearCommand := &comandos.ClearRecord{Cantidad: uint32(len(records))}
	m.conn.EnviarComando(clearCommand)

	return
}