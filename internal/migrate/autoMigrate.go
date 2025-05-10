package migrate

// type auth model.Auth
// type user model.User
// type incident model.Incident
// func StartMigration(configure config.Configuration) error {
// 	db := database.GetDB()
// 	configureDB := configure.Database.RDBMS
// 	driver := configureDB.Env.Driver
// 	if driver == "postgres" {
// 		if err := db.AutoMigrate(
// 			&auth{},
// 			&user{},
// 			&incident{},
// 		); err != nil {
// 			return err
// 		}
// 	}
// 	fmt.Println("new database migrated successfully!")
// 	return nil
// }
