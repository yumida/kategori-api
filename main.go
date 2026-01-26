package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Kategori struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
}

var kategori = []Kategori{
	{ID: 1, Nama: "Makanan", Deskripsi: "Makanan utama"},
	{ID: 2, Nama: "Minuman", Deskripsi: "Minuman segar"},
	{ID: 3, Nama: "Bumbu", Deskripsi: "Bumbu dapur"},
}

func getAllKategori(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Get all kategori")
	json.NewEncoder(w).Encode(kategori)
}

func getKategoriByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Get kategori by ID")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
		return
	}

	for _, p := range kategori {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	// Kalau tidak found
	http.Error(w, "Kategori belum ada", http.StatusNotFound)
}

func addKategori(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Add new kategori")

	var newKategori Kategori
	err := json.NewDecoder(r.Body).Decode(&newKategori)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	newKategori.ID = len(kategori) + 1
	kategori = append(kategori, newKategori)
	json.NewEncoder(w).Encode(newKategori)

}

func deleteKategori(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Delete kategori by ID")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
		return
	}

	for i, p := range kategori {
		if p.ID == id {
			kategori = append(kategori[:i], kategori[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}

	http.Error(w, "Kategori belum ada", http.StatusNotFound)
}

func updateKategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Update kategori by ID")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
		return
	}
	var updateKategori Kategori
	err = json.NewDecoder(r.Body).Decode(&updateKategori)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	for i := range kategori {
		if kategori[i].ID == id {
			updateKategori.ID = id
			kategori[i] = updateKategori
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateKategori)
			return
		}
	}

	http.Error(w, "Kategori belum ada", http.StatusNotFound)
}

func main() {

	http.HandleFunc("/api/kategori", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getAllKategori(w, r)
		case "POST":
			addKategori(w, r)
		}
	})

	http.HandleFunc("/api/kategori/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getKategoriByID(w, r)
		case "PUT":
			updateKategoryByID(w, r)
		case "DELETE":
			deleteKategori(w, r)
		}
	})

	fmt.Print("Server Running in localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Print("Server failed to run")
	}
}
