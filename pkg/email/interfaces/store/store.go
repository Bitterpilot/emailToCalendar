package store

import (
	"fmt"

	"github.com/bitterpilot/emailToCalendar/pkg/email"
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

type DbEmailRepo DbRepo

func NewDbEmailRepo(dbHandlers map[string]DbHandler) *DbEmailRepo {
	dbEmailRepo := new(DbEmailRepo)
	dbEmailRepo.dbHandlers = dbHandlers
	dbEmailRepo.dbHandler = dbHandlers["DbEmailRepo"]
	return dbEmailRepo
}

func (repo *DbEmailRepo) Store(msg *email.Msg) {
	if repo == nil {
		panic("DbEmailRepo cannot be nil")
	}
	repo.dbHandler.Execute(fmt.Sprintf(`
		INSERT INTO emails (ExternalID, ThreadID, ReceivedTime, Body)
		VALUES ('%s','%s','%d','%s')`,
		msg.ExternalID, msg.ThreadID, msg.ReceivedTime, msg.Body))
}

func (repo *DbEmailRepo) FindByID(id int) email.Msg {
	if repo == nil {
		panic("DbEmailRepo cannot be nil")
	}
	row := repo.dbHandler.Query(fmt.Sprintf(`
	SELECT
		ExternalID,
		ThreadID,
		ReceivedTime,
		Body
	FROM
		emails
	WHERE
		id = '%d'
	LIMIT 1`, id))

	var ExternalID,
		ThreadID,
		Body string

	var ReceivedTime int64

	row.Next()
	row.Scan(&ExternalID, &ThreadID, &ReceivedTime, &Body)
	m := email.Msg{
		ID:           id,
		ExternalID:   ExternalID,
		ThreadID:     ThreadID,
		ReceivedTime: ReceivedTime,
		Body:         Body,
	}

	return m
}

func (repo *DbEmailRepo) FindByExternalID(ExternalID string) email.Msg {
	if repo == nil {
		panic("DbEmailRepo cannot be nil")
	}
	row := repo.dbHandler.Query(fmt.Sprintf(`
	SELECT
		ID,
		ThreadID,
		ReceivedTime,
		Body
	FROM
		emails
	WHERE
		"ExternalID" = '%s'
	LIMIT 1;`, ExternalID))

	var id int
	var ThreadID,
		Body string

	var ReceivedTime int64

	row.Next()
	row.Scan(&id, &ThreadID, &ReceivedTime, &Body)
	m := email.Msg{
		ID:           id,
		ExternalID:   ExternalID,
		ThreadID:     ThreadID,
		ReceivedTime: ReceivedTime,
		Body:         Body,
	}

	return m
}

func (repo *DbEmailRepo) FindByExternalThreadID(ThreadID string) []email.Msg {
	if repo == nil {
		panic("DbEmailRepo cannot be nil")
	}
	row := repo.dbHandler.Query(fmt.Sprintf(`
	SELECT
		id,
		ExternalID,
		ReceivedTime,
		Body
	FROM
		emails
	WHERE
		ThreadID = '%s'`, ThreadID))

	var id int
	var ExternalID,
		Body string
	var ReceivedTime int64
	var m []email.Msg

	for row.Next() {
		row.Scan(&id, &ExternalID, &ReceivedTime, &Body)
		msg := email.Msg{
			ID:           id,
			ExternalID:   ExternalID,
			ThreadID:     ThreadID,
			ReceivedTime: ReceivedTime,
			Body:         Body,
		}
		m = append(m, msg)
	}

	return m
}
