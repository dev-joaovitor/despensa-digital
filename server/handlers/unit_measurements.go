package handlers

import (
	"fmt"
	"net/http"

	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
)

func (e *Env) UnitMeasurementsHandler(r chi.Router) {
	r.Get("/", e.ListMeasurementsHandler)
}

func (e *Env) ListMeasurementsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := e.DB.Query(
		r.Context(),
		`
		SELECT id, name, acronym
		FROM unit_measurements
		`,
	)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer rows.Close()

	measurements := []models.UnitMeasurement{}

	for rows.Next() {
		measurement := models.UnitMeasurement{}
		rows.Scan(
			&measurement.ID,
			&measurement.Name,
			&measurement.Acronym,
		)
		measurements = append(measurements, measurement)
	}

	WriteJSON(w, http.StatusOK, measurements, "Unidades de medida listadas com sucesso")
}

