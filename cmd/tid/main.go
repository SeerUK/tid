package main

import (
	"log"
	"os"

	"github.com/SeerUK/tid/pkg/state/bolt"
	"github.com/SeerUK/tid/pkg/state/migrate"
	_ "github.com/SeerUK/tid/pkg/state/migrate/versions"
	"github.com/SeerUK/tid/pkg/tid"
	"github.com/SeerUK/tid/pkg/tid/cli"
	"github.com/SeerUK/tid/pkg/util"
	boltdb "github.com/boltdb/bolt"
)

func main() {
	db := getBoltDB()
	defer db.Close()

	backend := bolt.NewBoltBackend(db)

	// Initialise the backend, preparing it for use, ensuring it's up-to-date.
	err := migrate.Backend(backend)
	fatal(err)

	factory := util.NewStandardFactory(backend)
	kernel := cli.NewTidKernel(backend, factory)

	os.Exit(cli.CreateApplication(kernel).Run(os.Args[1:], os.Environ()))
}

// getBoltDB gets a Bolt DB instance.
func getBoltDB() *boltdb.DB {
	dir, err := tid.GetLocalDirectory()
	fatal(err)

	// Open Bolt database.
	db, err := bolt.Open(dir)
	fatal(err)

	return db
}

// fatal kills the application upon error.
func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
