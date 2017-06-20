package gsql

import (
	"database/sql"
	"sort"
	"strings"
)

var Conn *sql.DB

// this allows us to the the go ` and replaces the MySql ` with ""
// so `"ID" int` becomes `ID` int
var quotes = strings.NewReplacer(`"`, "`")

type Migration struct {
	Up   string
	Down string
}

var Migrations = map[string]Migration{
	"000000000000_create_migrations": Migration{
		Up: quotes.Replace(`
			CREATE TABLE IF NOT EXISTS "migrations"(
			  "id" int unsigned not null auto_increment,
			  PRIMARY KEY("id"),
			  "Name" varchar(255)
			);`),
		Down: "DROP TABLE `migrations`;",
	},
}

func Migrate() {
	hasRun := getMigrationsRan()
	var run []string
	for name, _ := range Migrations {
		if !hasRun[name] {
			run = append(run, name)
		}
	}
	sort.Slice(run, func(i, j int) bool {
		return run[i] < run[j]
	})

	for _, name := range run {
		runMigration(name)
	}
}

func Describe() string {
	hasRan := getMigrationsRan()
	var all []string
	for name, _ := range Migrations {
		all = append(all, name)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i] < all[j]
	})
	for i, s := range all {
		if hasRan[s] {
			all[i] = "* " + s
		} else {
			all[i] = "  " + s
		}
	}
	return strings.Join(all, "\n")
}

func AddMigration(name, up, down string) Migration {
	m := Migration{
		Up:   quotes.Replace(up),
		Down: quotes.Replace(down),
	}
	Migrations[name] = m
	return m
}
func getMigrationsRan() map[string]bool {
	hasRan := make(map[string]bool)

	rows, err := db.Conn.Query("SHOW TABLES LIKE 'migrations';")
	if err != nil {
		panic(err)
	}
	if !rows.Next() {
		// not even migration migration has run yet
		return hasRan
	}

	rows, err = db.Conn.Query("SELECT `name` FROM `migrations`")
	if err != nil {
		panic(err)
	}
	var name string

	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			panic(err)
		}
		hasRan[name] = true
	}
	return hasRan
}

func runMigration(migration string) {
	_, err := db.Conn.Exec(Migrations[migration].Up)
	if err != nil {
		panic(err)
	}
	db.Conn.Exec("INSERT INTO migrations (name) VALUES (?)", migration)
}

func Rollback() string {
	rows, _ := db.Conn.Query("SELECT `id`,`Name` FROM `migrations` ORDER BY `Name` DESC LIMIT 1;")
	if !rows.Next() {
		// not even migration migration has run yet
		return ""
	}
	var id int
	var name string
	rows.Scan(&id, &name)
	db.Conn.Exec(Migrations[name].Down)
	db.Conn.Exec("DELETE FROM `migrations` WHERE `id`=?;", id)
	return name
}
