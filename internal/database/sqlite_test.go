package database

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/sharovik/devbot/internal/config"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"runtime"
	"testing"
)

const testSQLiteDatabasePath = "./test/testdata/database/devbot.sqlite"

var (
	cfg        config.Config
	dictionary SQLiteDictionary
)

func init() {
	//We switch pointer to the root directory for control the path from which we need to generate test-data file-paths
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	_ = os.Chdir(dir)
}

func TestSQLiteDictionary_InitDatabaseConnection(t *testing.T) {
	cfg.DatabaseHost = "./wrong_path"
	dictionary.cfg = cfg

	db, err := dictionary.InitDatabaseConnection()

	assert.Error(t, err)
	assert.Empty(t, db)

	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.cfg = cfg

	db, err = dictionary.InitDatabaseConnection()

	assert.NoError(t, err)
	assert.NotEmpty(t, db)

	defer db.Close()

	sqlStmt := `
	drop table if exists foo;
	create table foo (id integer not null primary key, name text);
	`
	_, err = db.Exec(sqlStmt)
	assert.NoError(t, err)

	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	assert.NoError(t, err)

	rows, err := db.Query("select id, name from foo where id = 1")
	assert.NoError(t, err)
	defer rows.Close()

	var id int
	var name string

	for rows.Next() {
		err = rows.Scan(&id, &name)
		assert.NoError(t, err)
	}

	assert.Equal(t, 1, id)
	assert.Equal(t, "foo", name)

	err = rows.Err()
	assert.NoError(t, err)
}
