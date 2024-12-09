package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers" // Asegúrate de incluir esto
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Actividad struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

var dataSource string

func main() {
	// Cargar configuración desde .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Conexión a MySQL

	dataSource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	r := mux.NewRouter()
	r.HandleFunc("/actividades", getActividades).Methods("GET")
	r.HandleFunc("/actividades", createActividad).Methods("POST")
	r.HandleFunc("/actividades/{id}", updateActividad).Methods("PUT")
	r.HandleFunc("/actividades/{id}", deleteActividad).Methods("DELETE")

	// Habilitar CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	// Registrar el router con CORS
	loggedRouter := handlers.CORS(headersOk, originsOk, methodsOk)(r)

	log.Println("API ejecutándose en http://localhost:8080")

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", loggedRouter))

}

func getActividades(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, nombre, descripcion FROM actividades")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	actividades := []Actividad{}
	for rows.Next() {
		var actividad Actividad
		if err := rows.Scan(&actividad.ID, &actividad.Nombre, &actividad.Descripcion); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		actividades = append(actividades, actividad)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actividades)
}

func createActividad(w http.ResponseWriter, r *http.Request) {
	var nuevaActividad Actividad
	if err := json.NewDecoder(r.Body).Decode(&nuevaActividad); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := "INSERT INTO actividades (nombre, descripcion) VALUES (?, ?)"
	result, err := db.Exec(query, nuevaActividad.Nombre, nuevaActividad.Descripcion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nuevaActividad.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevaActividad)
}

func updateActividad(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var actividad Actividad
	if err := json.NewDecoder(r.Body).Decode(&actividad); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := "UPDATE actividades SET nombre = ?, descripcion = ? WHERE id = ?"
	_, err = db.Exec(query, actividad.Nombre, actividad.Descripcion, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actividad)
}

func deleteActividad(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := "DELETE FROM actividades WHERE id = ?"
	_, err = db.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
