package migrate

import (
	"fmt"

	"github.com/Dhar01/incident_resp/config"
	"github.com/Dhar01/incident_resp/internal/database"
	"github.com/Dhar01/incident_resp/internal/model"
)

type auth model.Auth
type user model.User

func StartMigration(configure config.Configuration) error {
	db := database.GetDB()

	configureDB := configure.Database.RDBMS
	driver := configureDB.Env.Driver

	if driver == "postgres" {
		if err := db.AutoMigrate(
			&auth{},
			&user{},
		); err != nil {
			return err
		}
	}

	fmt.Println("new database migrated successfully!")

	return nil
}
