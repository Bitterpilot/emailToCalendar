package interfaces

import (
	"fmt"

	"github.com/bitterpilot/emailToCalendar/pkg/user"
)

type DbHandler interface {
	Execute(statement string)
	Query(statement string) Row
}
type Row interface {
	Scan(dest ...interface{})
	Next() bool
}
type DbRepo struct {
	dbHandlers map[string]DbHandler
	dbHandler  DbHandler
}

type DbUserRepo DbRepo

func NewDbUserRepo(dbHandlers map[string]DbHandler) *DbUserRepo {
	dbUserRepo := new(DbUserRepo)
	dbUserRepo.dbHandlers = dbHandlers
	dbUserRepo.dbHandler = dbHandlers["DbUserRepo"]
	return dbUserRepo
}

func (repo *DbUserRepo) FindByID(id int) user.User {
	row := repo.dbHandler.Query(fmt.Sprintf(`
	SELECT is_admin, customer_id
	FROM users WHERE id = '%d' LIMIT 1`, id))
	var isAdmin string
	var customerID int
	row.Next()
	row.Scan(&isAdmin, &customerID)
	u := user.User{Id: id}
	u.IsAdmin = false
	if isAdmin == "yes" {
		u.IsAdmin = true
	}
	return u
}

func (repo *DbUserRepo) Store(customer user.User) {
	if repo == nil {
		panic("DbUserRepo cannot be nil")
	}
	repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO customers (name) VALUES ('%s')`,
		customer.Name))
}
