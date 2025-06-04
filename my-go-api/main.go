package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"

    _ "github.com/lib/pq"
    "github.com/gorilla/mux"
)

var db *sql.DB

type Drug struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Quantity int     `json:"quantity"`
    Price    float64 `json:"price"`
}

func main() {
    var err error
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

    db, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    router := mux.NewRouter()

    router.HandleFunc("/drugs", getAllDrugs).Methods("GET")
    router.HandleFunc("/drugs", createDrug).Methods("POST")
    router.HandleFunc("/drugs/{id}", updateDrug).Methods("PUT")
    router.HandleFunc("/drugs/{id}", deleteDrug).Methods("DELETE")

    fmt.Println("Server running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}

// Get all drugs
func getAllDrugs(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name, quantity, price FROM drugs")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var drugs []Drug
    for rows.Next() {
        var d Drug
        if err := rows.Scan(&d.ID, &d.Name, &d.Quantity, &d.Price); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        drugs = append(drugs, d)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(drugs)
}

// Create new drug
func createDrug(w http.ResponseWriter, r *http.Request) {
    var d Drug
    if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := db.QueryRow(
        "INSERT INTO drugs (name, quantity, price) VALUES ($1, $2, $3) RETURNING id",
        d.Name, d.Quantity, d.Price).Scan(&d.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(d)
}

// Update drug by ID
func updateDrug(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid drug ID", http.StatusBadRequest)
        return
    }

    var d Drug
    if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    res, err := db.Exec(
        "UPDATE drugs SET name=$1, quantity=$2, price=$3 WHERE id=$4",
        d.Name, d.Quantity, d.Price, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if rowsAffected == 0 {
        http.Error(w, "Drug not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// Delete drug by ID
func deleteDrug(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid drug ID", http.StatusBadRequest)
        return
    }

    res, err := db.Exec("DELETE FROM drugs WHERE id=$1", id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if rowsAffected == 0 {
        http.Error(w, "Drug not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent) // 204 No Content
}
