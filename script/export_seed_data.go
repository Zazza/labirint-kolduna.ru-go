package script

import (
	"flag"
	"log"
	"os"

	"gamebook-backend/database"
	"gamebook-backend/database/seeders/exporter"
)

func main() {
	var (
		host     string
		port     string
		user     string
		password string
		dbname   string
		sslmode  string
		output   string
	)

	flag.StringVar(&host, "host", "localhost", "Database host")
	flag.StringVar(&port, "port", "5432", "Database port")
	flag.StringVar(&user, "user", "postgres", "Database user")
	flag.StringVar(&password, "password", "", "Database password")
	flag.StringVar(&dbname, "dbname", "gamebook_develop", "Database name")
	flag.StringVar(&sslmode, "sslmode", "disable", "SSL mode")
	flag.StringVar(&output, "output", "./database/seeders/json", "Output directory")

	flag.Parse()

	if password == "" {
		if envPassword := os.Getenv("DB_PASSWORD"); envPassword != "" {
			password = envPassword
		} else {
			log.Fatal("Database password is required. Use -password flag or DB_PASSWORD environment variable")
		}
	}

	db, err := database.Connect(host, port, user, password, dbname, sslmode)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("Connected to database: %s", dbname)

	exporter := exporter.NewExporter(db, output)

	log.Println("Starting data export...")

	if err := exporter.ExportAll(); err != nil {
		log.Fatalf("Export failed: %v", err)
	}

	log.Println("Export completed successfully!")
}
