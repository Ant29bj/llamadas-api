package filtros

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"registro_llamadas/modelos"
)

func GetLlamadas(w http.ResponseWriter, r *http.Request) {
	// Obtener los parámetros de paginación de la consulta
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSizeStr := r.URL.Query().Get("page_size")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Calcular el offset en función de la página y el tamaño de la página
	offset := (page - 1) * pageSize

	// Abrir la conexión a la base de datos
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/llamadas")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Realizar la consulta paginada
	rows, err := db.Query("SELECT * FROM cdr ORDER BY calldate DESC LIMIT ?, ?", offset, pageSize)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Leer los resultados y almacenarlos en una lista de CDRs
	cdrs := make([]modelos.CDR, 0)
	for rows.Next() {
		var cdr modelos.CDR
		err := rows.Scan(
			&cdr.CallDate,
			&cdr.Clid,
			&cdr.Src,
			&cdr.Dst,
			&cdr.Dcontext,
			&cdr.Channel,
			&cdr.DstChannel,
			&cdr.LastApp,
			&cdr.LastData,
			&cdr.Duration,
			&cdr.Billsec,
			&cdr.Disposition,
			&cdr.Amaflags,
			&cdr.Accountcode,
			&cdr.UniqueID,
			&cdr.UserField,
			&cdr.Did,
			&cdr.RecordingFile,
			&cdr.Cnum,
			&cdr.Cnam,
			&cdr.OutboundCnum,
			&cdr.OutboundCnam,
			&cdr.DstCnam,
			&cdr.LinkedID,
			&cdr.PeerAccount,
			&cdr.Sequence,
		)
		if err != nil {
			log.Fatal(err)
		}
		cdrs = append(cdrs, cdr)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Obtener el total de registros en la tabla
	var totalRecords int
	err = db.QueryRow("SELECT COUNT(*) FROM cdr").Scan(&totalRecords)
	if err != nil {
		log.Fatal(err)
	}

	totalPages := totalRecords / pageSize
	if totalPages%pageSize != 0 {
		totalPages++
	}

	// Crear una estructura para la respuesta paginada
	type PaginatedResponse struct {
		TotalRecords int           `json:"total_records"`
		Page         int           `json:"page"`
		TotalPages   int           `json:"total_pages"`
		PageSize     int           `json:"page_size"`
		Records      []modelos.CDR `json:"records"`
	}

	// Construir la respuesta paginada
	response := PaginatedResponse{
		TotalRecords: totalRecords,
		Page:         page,
		TotalPages:   totalPages,
		PageSize:     pageSize,
		Records:      cdrs,
	}

	// Convertir la respuesta a JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	// Establecer las cabeceras de la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Enviar la respuesta JSON
	w.Write(jsonResponse)
}

func GetLlamadasPorDisposition(w http.ResponseWriter, r *http.Request) {
	// Leer el cuerpo de la solicitud
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	// Decodificar el cuerpo JSON en una estructura
	type RequestBody struct {
		Disposition string `json:"disposition"`
		Page        int    `json:"page"`
		PageSize    int    `json:"page_size"`
	}

	var requestBody RequestBody
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		log.Fatal(err)
	}

	// Obtener los valores del cuerpo de la solicitud
	disposition := requestBody.Disposition
	page := requestBody.Page
	pageSize := requestBody.PageSize

	// Calcular el offset en función de la página y el tamaño de la página
	offset := (page - 1) * pageSize

	// Abrir la conexión a la base de datos
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/llamadas")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Realizar la consulta filtrada por disposition y paginada
	rows, err := db.Query("SELECT * FROM cdr WHERE disposition = ? ORDER BY calldate DESC LIMIT ?, ?", disposition, offset, pageSize)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Leer los resultados y almacenarlos en una lista de CDRs
	cdrs := make([]modelos.CDR, 0)
	for rows.Next() {
		var cdr modelos.CDR
		err := rows.Scan(
			&cdr.CallDate,
			&cdr.Clid,
			&cdr.Src,
			&cdr.Dst,
			&cdr.Dcontext,
			&cdr.Channel,
			&cdr.DstChannel,
			&cdr.LastApp,
			&cdr.LastData,
			&cdr.Duration,
			&cdr.Billsec,
			&cdr.Disposition,
			&cdr.Amaflags,
			&cdr.Accountcode,
			&cdr.UniqueID,
			&cdr.UserField,
			&cdr.Did,
			&cdr.RecordingFile,
			&cdr.Cnum,
			&cdr.Cnam,
			&cdr.OutboundCnum,
			&cdr.OutboundCnam,
			&cdr.DstCnam,
			&cdr.LinkedID,
			&cdr.PeerAccount,
			&cdr.Sequence,
		)
		if err != nil {
			log.Fatal(err)
		}
		cdrs = append(cdrs, cdr)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Obtener el total de registros en la tabla (sin filtro aplicado)
	var totalRecords int
	err = db.QueryRow("SELECT COUNT(*) FROM cdr WHERE disposition = ?", disposition).Scan(&totalRecords)
	if err != nil {
		log.Fatal(err)
	}

	totalPages := totalRecords / pageSize
	if totalPages%pageSize != 0 {
		totalPages++
	}

	// Crear una estructura para la respuesta paginada
	type PaginatedResponse struct {
		TotalRecords int           `json:"total_records"`
		Page         int           `json:"page"`
		TotalPages   int           `json:"total_pages"`
		PageSize     int           `json:"page_size"`
		Records      []modelos.CDR `json:"records"`
	}

	// Construir la respuesta paginada
	response := PaginatedResponse{
		TotalRecords: totalRecords,
		Page:         page,
		TotalPages:   totalPages,
		PageSize:     pageSize,
		Records:      cdrs,
	}

	// Convertir la respuesta a JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	// Establecer las cabeceras de la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Enviar la respuesta JSON
	w.Write(jsonResponse)
}

func GetLlamadasPorRangoFecha(w http.ResponseWriter, r *http.Request) {
	// Leer el cuerpo de la solicitud
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	// Decodificar el cuerpo JSON en una estructura
	type RequestBody struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Page      int    `json:"page"`
		PageSize  int    `json:"page_size"`
	}

	var requestBody RequestBody
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		log.Fatal(err)
	}

	// Obtener los valores del cuerpo de la solicitud
	startDate := requestBody.StartDate
	endDate := requestBody.EndDate
	page := requestBody.Page
	pageSize := requestBody.PageSize

	// Calcular el offset en función de la página y el tamaño de la página
	offset := (page - 1) * pageSize

	// Abrir la conexión a la base de datos
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/llamadas")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Realizar la consulta filtrada por rango de fechas y paginada
	rows, err := db.Query("SELECT * FROM cdr WHERE calldate >= ? AND calldate <= ? ORDER BY calldate DESC LIMIT ?, ?", startDate, endDate, offset, pageSize)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Leer los resultados y almacenarlos en una lista de CDRs
	cdrs := make([]modelos.CDR, 0)
	for rows.Next() {
		var cdr modelos.CDR
		err := rows.Scan(
			&cdr.CallDate,
			&cdr.Clid,
			&cdr.Src,
			&cdr.Dst,
			&cdr.Dcontext,
			&cdr.Channel,
			&cdr.DstChannel,
			&cdr.LastApp,
			&cdr.LastData,
			&cdr.Duration,
			&cdr.Billsec,
			&cdr.Disposition,
			&cdr.Amaflags,
			&cdr.Accountcode,
			&cdr.UniqueID,
			&cdr.UserField,
			&cdr.Did,
			&cdr.RecordingFile,
			&cdr.Cnum,
			&cdr.Cnam,
			&cdr.OutboundCnum,
			&cdr.OutboundCnam,
			&cdr.DstCnam,
			&cdr.LinkedID,
			&cdr.PeerAccount,
			&cdr.Sequence,
		)
		if err != nil {
			log.Fatal(err)
		}
		cdrs = append(cdrs, cdr)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Obtener el total de registros en la tabla (sin filtro aplicado)
	var totalRecords int
	err = db.QueryRow("SELECT COUNT(*) FROM cdr WHERE calldate >= ? AND calldate <= ?", startDate, endDate).Scan(&totalRecords)
	if err != nil {
		log.Fatal(err)
	}

	totalPages := totalRecords / pageSize
	if totalPages%pageSize != 0 {
		totalPages++
	}

	// Crear una estructura para la respuesta paginada
	type PaginatedResponse struct {
		TotalRecords int           `json:"total_records"`
		Page         int           `json:"page"`
		TotalPages   int           `json:"total_pages"`
		PageSize     int           `json:"page_size"`
		Records      []modelos.CDR `json:"records"`
	}

	// Construir la respuesta paginada
	response := PaginatedResponse{
		TotalRecords: totalRecords,
		Page:         page,
		TotalPages:   totalPages,
		PageSize:     pageSize,
		Records:      cdrs,
	}

	// Convertir la respuesta a JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	// Establecer las cabeceras de la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Enviar la respuesta JSON
	w.Write(jsonResponse)
}

func GetLlamadasPorSrc(w http.ResponseWriter, r *http.Request) {
	// Obtener los parámetros de paginación y filtro de la consulta
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSizeStr := r.URL.Query().Get("page_size")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	srcFilter := r.URL.Query().Get("src")

	// Calcular el offset en función de la página y el tamaño de la página
	offset := (page - 1) * pageSize

	// Abrir la conexión a la base de datos
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/llamadas")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Construir la consulta SQL con el filtro
	query := "SELECT * FROM cdr"
	args := make([]interface{}, 0)

	if srcFilter != "" {
		query += " WHERE src = ?"
		args = append(args, srcFilter)
	}

	query += " ORDER BY calldate DESC LIMIT ?, ?"
	args = append(args, offset, pageSize)

	// Realizar la consulta paginada con filtro
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Leer los resultados y almacenarlos en una lista de CDRs
	cdrs := make([]modelos.CDR, 0)
	for rows.Next() {
		var cdr modelos.CDR
		err := rows.Scan(
			&cdr.CallDate,
			&cdr.Clid,
			&cdr.Src,
			&cdr.Dst,
			&cdr.Dcontext,
			&cdr.Channel,
			&cdr.DstChannel,
			&cdr.LastApp,
			&cdr.LastData,
			&cdr.Duration,
			&cdr.Billsec,
			&cdr.Disposition,
			&cdr.Amaflags,
			&cdr.Accountcode,
			&cdr.UniqueID,
			&cdr.UserField,
			&cdr.Did,
			&cdr.RecordingFile,
			&cdr.Cnum,
			&cdr.Cnam,
			&cdr.OutboundCnum,
			&cdr.OutboundCnam,
			&cdr.DstCnam,
			&cdr.LinkedID,
			&cdr.PeerAccount,
			&cdr.Sequence,
		)
		if err != nil {
			log.Fatal(err)
		}
		cdrs = append(cdrs, cdr)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Obtener el total de registros en la tabla (con el filtro aplicado)
	var totalRecords int
	if srcFilter != "" {
		err = db.QueryRow("SELECT COUNT(*) FROM cdr WHERE src = ?", srcFilter).Scan(&totalRecords)
	} else {
		err = db.QueryRow("SELECT COUNT(*) FROM cdr").Scan(&totalRecords)
	}
	if err != nil {
		log.Fatal(err)
	}

	totalPages := totalRecords / pageSize
	if totalPages%pageSize != 0 {
		totalPages++
	}

	// Crear una estructura para la respuesta paginada
	type PaginatedResponse struct {
		TotalRecords int           `json:"total_records"`
		Page         int           `json:"page"`
		TotalPages   int           `json:"total_pages"`
		PageSize     int           `json:"page_size"`
		Records      []modelos.CDR `json:"records"`
	}

	// Construir la respuesta paginada
	response := PaginatedResponse{
		TotalRecords: totalRecords,
		Page:         page,
		TotalPages:   totalPages,
		PageSize:     pageSize,
		Records:      cdrs,
	}

	// Convertir la respuesta a JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	// Establecer las cabeceras de la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Enviar la respuesta JSON
	w.Write(jsonResponse)
}
