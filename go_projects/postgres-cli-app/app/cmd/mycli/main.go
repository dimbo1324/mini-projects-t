package main

import (
	"flag"
	"log"
	"mycli/internal/app"
	"mycli/internal/config"
	"mycli/internal/db"
)

var flags struct {
	ConfigPath    string
	DirectoryPath string
}

func init() {
	flag.StringVar(&flags.ConfigPath, "c", "", "path to yaml schema config files")
	flag.StringVar(&flags.DirectoryPath, "d", "", "path to directory with files")
}

const postgresqlConnString = "postgres://user:password@localhost:5432/postgres"

func main() {
	flag.Parse()

	if flags.ConfigPath == "" || flags.DirectoryPath == "" {
		log.Fatal("Both -c and -d flags are required")
	}

	database, err := db.NewClient(postgresqlConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	schemas, err := config.LoadConfig(flags.ConfigPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	processor := app.NewProcessor(database, schemas)

	log.Println("Starting processing...")
	if err := processor.ProcessDirectory(flags.DirectoryPath); err != nil {
		log.Fatalf("Processing failed: %v", err)
	}
	log.Println("Done.")
}
