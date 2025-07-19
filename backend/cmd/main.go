package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// global o'zgaruvchi

type User struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	Age     int    `json:"age"`
}
type Foods struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	ImagePath   string `json:"image_path"`
}

func main() {
	connStr := "user=postgres password=1234 host=localhost port=5432 dbname=mydatabase sslmode=disable"
	var err error
	userdb, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("PostgreSql ga ulanishda xatolik: %v", err)
	}
	err = userdb.Ping()
	if err != nil {
		log.Fatalf("Ma'lumotlar bazasiga ping qila olmadi: %v", err)
	}
	fmt.Println("✅ PostgreSql ga muvaffaqiyatli ulandi")
	defer userdb.Close()

	foodsdb, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Postgressga ulanishda xatolik: %v", err)
	}
	err = foodsdb.Ping()
	if err != nil {
		log.Fatalf("Malmuotlar bazasiga Ping qila olmadi: %v", err)
	}
	fmt.Println("Foods Malumotlar bazasiga muvafaqqiyatli ulandi✅✅✅")
	defer foodsdb.Close()
	// 	deleteTableSql :="DROP TABLE foods"
	// _, err = db.Exec(deleteTableSql)
	// if err != nil {
	// 	log.Fatalf("users jadvalini ochirishda xatolik: %v", err)
	// }
	// fmt.Println("✅ Foods jadvali ochirildi")

	createFoodsTable(foodsdb)

	http.HandleFunc("/foods", getFoodsHandler(foodsdb))
	http.HandleFunc("/add-foods",addFoodsHandler(foodsdb))
	fmt.Println("Server :8080 portida ishlayabti")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func createFoodsTable(foodsdb *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS foods (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	price INTEGER NOT NULL,
	image_path	VARCHAR(105) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err := foodsdb.Exec(query)
	if err != nil {
		log.Fatalf("Foods jadvalini yaratishda xatolik %v", err)
	}
	fmt.Println("✅✅✅ Foods malumotlar bazasi muvafaqqiyatli yaratildi")
}

func addFoodsHandler(foodsdb *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
		INSERT INTO foods (name,price,image_path) VALUES($1,$2,$3)
		`
		_, err := foodsdb.Exec(query, "Shorva", 20.000, "/backend/uploads/img/shorva.jpg")
		if err != nil {
			fmt.Println(err, "Foods tablega❌❌❌ Ma'lumotlarni yozishda xatolikboldi")
		}
		fmt.Println("Muvafaqqiyatli qoshildi ✅✅✅")
		rows, err := foodsdb.Query("SELECT name,price,image_path FROM foods")
		if err != nil {
			log.Println("Foods malumotini oqishda xatolik❌❌❌", err)
			http.Error(w, "Xatolik", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		fmt.Println("Ovaqatlarroyxati")
		for rows.Next() {
			var id, price int
			var name, image_path string
			err = rows.Scan(&id, &name, &price, &image_path)

			if err != nil {
				log.Println("Foods satrini oqishda xatolik")
				continue
			}
			w.WriteHeader(http.StatusOK)
			fmt.Printf(name, price, image_path)
		}
	}
}
func getFoodsHandler(foodsdb *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
