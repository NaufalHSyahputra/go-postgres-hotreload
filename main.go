package main

import (
	"fmt"
	"go-postgres-docker/lib"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var schema = `
CREATE SCHEMA IF NOT EXISTS public;
DROP TABLE IF EXISTS public.product;
CREATE TABLE IF NOT EXISTS product (
	id serial NOT NULL,
	name character varying(255) NULL
  );
  TRUNCATE TABLE public.product;
  INSERT INTO public.product(name) VALUES ('Keyboard');
  INSERT INTO public.product(name) VALUES ('Mouse');
  INSERT INTO public.product(name) VALUES ('Mousepad');
  INSERT INTO public.product(name) VALUES ('Keycaps');
  INSERT INTO public.product(name) VALUES ('Laptop');
  INSERT INTO public.product(name) VALUES ('Memory RAM');
  INSERT INTO public.product(name) VALUES ('Processor');
  INSERT INTO public.product(name) VALUES ('Harddisk');
  INSERT INTO public.product(name) VALUES ('SSD');
  INSERT INTO public.product(name) VALUES ('Flash Disk');
  INSERT INTO public.product(name) VALUES ('Headphone');
  INSERT INTO public.product(name) VALUES ('Headset');
  INSERT INTO public.product(name) VALUES ('Speaker');
  INSERT INTO public.product(name) VALUES ('Mic');
  INSERT INTO public.product(name) VALUES ('Fan Cooler');
  INSERT INTO public.product(name) VALUES ('Printer');
  INSERT INTO public.product(name) VALUES ('Charger');
  INSERT INTO public.product(name) VALUES ('Other Accessories');
  INSERT INTO public.product(name) VALUES ('Other components');
`

type Product struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sqlx.Connect("postgres", dsn)
	log.Println(dsn)
	log.Println("test")
	db.MustExec(schema)

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		lib.ResponseJSON(w, 200, map[string]any{
			"message": "Hello WOrld",
		})
	})
	r.Get("/get-all", func(w http.ResponseWriter, r *http.Request) {
		product := []Product{}
		err = db.Select(&product, "SELECT * FROM public.product")
		lib.ResponseJSON(w, 200, map[string]any{
			"message": "Get Dataa Success",
			"data":    product,
			"error":   err,
		})
	})
	http.ListenAndServe(":8080", r)
}
