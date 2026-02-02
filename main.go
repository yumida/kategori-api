package main

import (
	"fmt"
	"kategori-api/database"
	"kategori-api/handlers"
	"kategori-api/repositories"
	"kategori-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

// type Kategori struct {
// 	ID        int    `json:"id"`
// 	Nama      string `json:"nama"`
// 	Deskripsi string `json:"deskripsi"`
// }

// var kategori = []Kategori{
// 	{ID: 1, Nama: "Makanan", Deskripsi: "Makanan utama"},
// 	{ID: 2, Nama: "Minuman", Deskripsi: "Minuman segar"},
// 	{ID: 3, Nama: "Bumbu", Deskripsi: "Bumbu dapur"},
// }

// func getAllKategori(w http.ResponseWriter, _ *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Println("Get all kategori")
// 	json.NewEncoder(w).Encode(kategori)
// }

// func getKategoriByID(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Println("Get kategori by ID")

// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
// 		return
// 	}

// 	for _, p := range kategori {
// 		if p.ID == id {
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(p)
// 			return
// 		}
// 	}

// 	// Kalau tidak found
// 	http.Error(w, "Kategori belum ada", http.StatusNotFound)
// }

// func addKategori(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Println("Add new kategori")

// 	var newKategori Kategori
// 	err := json.NewDecoder(r.Body).Decode(&newKategori)
// 	if err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}
// 	newKategori.ID = len(kategori) + 1
// 	kategori = append(kategori, newKategori)
// 	json.NewEncoder(w).Encode(newKategori)

// }

// func deleteKategori(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Println("Delete kategori by ID")

// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
// 		return
// 	}

// 	for i, p := range kategori {
// 		if p.ID == id {
// 			kategori = append(kategori[:i], kategori[i+1:]...)
// 			json.NewEncoder(w).Encode(map[string]string{
// 				"message": "sukses delete",
// 			})
// 			return
// 		}
// 	}

// 	http.Error(w, "Kategori belum ada", http.StatusNotFound)
// }

// func updateKategoryByID(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Println("Update kategori by ID")

// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
// 		return
// 	}
// 	var updateKategori Kategori
// 	err = json.NewDecoder(r.Body).Decode(&updateKategori)
// 	if err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}
// 	for i := range kategori {
// 		if kategori[i].ID == id {
// 			updateKategori.ID = id
// 			kategori[i] = updateKategori
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(updateKategori)
// 			return
// 		}
// 	}

// 	http.Error(w, "Kategori belum ada", http.StatusNotFound)
// }

func main() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	http.HandleFunc("/api/kategori", func(w http.ResponseWriter, r *http.Request) {
		categoryHandler.HandleCategory(w, r)
	})

	http.HandleFunc("/api/kategori/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			categoryHandler.GetByID(w, r)
		case "PUT":
			categoryHandler.UpdateByID(w, r)
		case "DELETE":
			categoryHandler.Delete(w, r)
		}
	})

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	fmt.Print("Server Running in localhost:8080")
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Print("Server failed to run")
	}
}
