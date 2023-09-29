package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	_ "ariga.io/atlas-go-sdk/recordriver"

	"github.com/Smylet/symlet-backend/utilities/db"
)

func main() {
    stmts, err := gormschema.New("postgres").Load(
		db.GetMigrateModels()...
	)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
        os.Exit(1)
    }
    io.WriteString(os.Stdout, stmts)
}