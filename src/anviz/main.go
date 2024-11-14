package main

import (
	"github.com/xfxdev/xlog"
	"anviz/core"
	"anviz/core/comandos"
	"sync"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"encoding/json"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "thePassWord"
	dbname   = "theDB"
)

type Dispositivo struct {
	Id int
	Direccion string
}

type ResultadoCanal struct {
	Dispositivo_id int `json:"dispositivo_id"`
	Records []comandos.Record `json:"records"`
}

type StoreParam struct {

}

func main() {
	//psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	//	"password=%s dbname=%s sslmode=disable",
	//	host, port, user, password, dbname)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	dispositivos := obtenerDispositivos(db)

	lenDisp := len(dispositivos)

	c := make(chan ResultadoCanal)

	var semaforo sync.Mutex

	managers := make([]core.Manager, lenDisp)

	cantidadDispositivosLeidos := 0
	cantidadRegistrosLeidos := 0

	for i:=0; i < lenDisp; i++ {
		go obtenerRegistros(dispositivos[i], managers[i], &cantidadDispositivosLeidos, &lenDisp, c, &semaforo)
	}

	for resultado := range c {

		param, err := json.Marshal(resultado)
		if err != nil {
			xlog.Error("Ha ocurrido un error construyendo el json")
		}

		cantidad := len(resultado.Records)

		if cantidad > 0 {
			db.Query("select * from grabar_fichada($1)", param)
		}

		cantidadRegistrosLeidos += len(resultado.Records)
	}

	xlog.Info("=========================== ")
	xlog.Info("Cantidad total de registros leidos: ", cantidadRegistrosLeidos)
}

func obtenerDispositivos(db *sql.DB) (dispositivos []Dispositivo) {
	rows, err := db.Query("select id, ip || ':' || puerto as direccion from dispositivo where habilitado")
	if err != nil {
		xlog.Error("Ha ocurrido un error leyendo los dispositivos desde la base")
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var direccion string

		err = rows.Scan(&id, &direccion)
		if err != nil {
			xlog.Error("Ha ocurrido un error leyendo los dispositivos desde la base, tratando de decodificar las rows")
			panic(err)
		}

		dispositivos = append(dispositivos, Dispositivo{Id: id, Direccion: direccion})
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		xlog.Error("Ha ocurrido un error leyendo los dispositivos desde la base, iterando las rows")
		panic(err)
	}
	return
}

func obtenerRegistros(dispositivo Dispositivo, manager core.Manager, cantidadDispositivosLeidos *int, cantidadDispositivos *int, resultados chan ResultadoCanal, semaforo *sync.Mutex) {

	err := manager.Conectar(dispositivo.Direccion)
	registros := make([]comandos.Record, 0)

	if err != nil {
		xlog.Error(err)
	} else {
		registros, err = manager.ObtenerRegistros()
		if err != nil {
			xlog.Error(err)
		}
	}

	resultados <- ResultadoCanal{Dispositivo_id: dispositivo.Id, Records: registros}

	semaforo.Lock()

	*cantidadDispositivosLeidos++
	if *cantidadDispositivosLeidos == *cantidadDispositivos {
		close(resultados)
	}

	semaforo.Unlock()

	xlog.Info(dispositivo.Direccion, ", cantidad: ", len(registros))
}