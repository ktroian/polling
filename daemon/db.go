package main

import (
	"log"
    "strings"
    "context"

    "github.com/go-pg/pg/v9"
    "github.com/go-pg/pg/v9/orm"
    "github.com/ktroian/polling/daemon/action"
)

const (
    SelectOne = iota
	InsertOne = iota
	DeleteOne = iota
)

func connectToDatabase(ctx context.Context, uname, password string) {
    go func() {
        db := pg.Connect(&pg.Options{
            User: uname,
            Password: password,
        })
        defer db.Close()
        createSchema(db)

        for {
            select {
            case <- ctx.Done():
                return

            default:
                var err error
            	act := action.Pop()
                company := (*Company)(act.Target)
            	
            	switch act.Type {
                case SelectOne:
                    err = db.Select(company)

        	    case InsertOne:
        	    	err = db.Insert(company)
                    
                    if err != nil {
                        if strings.Contains(err.Error(), "#23505") {
                            // #23505 is for duplicating data in a DB
                            // so duplicates will be skipped
                            continue
                        }
                    }

            	case DeleteOne:
                    err = db.Delete(company)
            	}

        	    if err != nil {
        	        log.Println(err)
        	    }
            }

        }
    }()
}

func createSchema(db *pg.DB) error {
    for _, model := range []interface{}{(*Company)(nil)} {
        err := db.CreateTable(model, &orm.CreateTableOptions{
            Temp: false,
        })
        if err != nil {
            return err
        }
    }
    return nil
}
