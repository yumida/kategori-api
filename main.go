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

	http.HandleFunc("/api/category", func(w http.ResponseWriter, r *http.Request) {
		categoryHandler.HandleCategory(w, r)
	})

	http.HandleFunc("/api/category/", func(w http.ResponseWriter, r *http.Request) {
		categoryHandler.HandleCategoryById(w, r)
	})

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	fmt.Print("Server Running in localhost:8080")
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Print("Server failed to run")
	}
}
