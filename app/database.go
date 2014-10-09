package app

import (
  "database/sql"
  //"fmt"
  "github.com/coopernurse/gorp"
  _ "github.com/mattn/go-sqlite3"
  . "instago/app/models"
  "time"
)

var dayDuration time.Duration = 24 * time.Hour

type Database struct {
  *gorp.DbMap
}

func MakeDatabase(dbFile string) (*Database, error) {
  db, err := sql.Open("sqlite3", dbFile)
  if err != nil {
    return nil, err
  }
  return &Database{
    &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}},
  }, nil
}

func (db *Database) Orders(query string) ([]Order, error) {
  orders := []Order{}
  _, err := db.Select(&orders, query)
  return orders, err
}

func (db *Database) Cohorts(d time.Duration, length time.Duration) ([]*Cohort, error) {
  cohorts := []*Cohort{}
  orders, err := db.Orders("SELECT * FROM orders ORDER BY created_at DESC")
  if err != nil {
    return cohorts, err
  }

  lastDay := orders[0].Created.Truncate(dayDuration)
  startDay := lastDay.Add(-length)

  for currentDay := startDay; currentDay.Before(lastDay); currentDay = currentDay.Add(d) {
    newCohort := &Cohort{currentDay, currentDay.Add(d), []User{}}
    cohorts = append(cohorts, newCohort)
  }

  users := []User{}
  _, err = db.Select(&users, "SELECT * FROM users WHERE created_at >= ? ORDER BY created_at ASC", startDay)
  if err != nil {
    return cohorts, err
  }

  cohortIdx := 0
  cohort := cohorts[cohortIdx]
  for _, user := range users {
    if user.Created.After(cohort.End) {
      cohortIdx++
      cohort = cohorts[cohortIdx]
    }
    cohort.Users = append(cohort.Users, user)
  }
  return cohorts, nil
}

func (db *Database) CacheFirstOrders(userIds string) {
  db.Exec("DROP TABLE IF EXISTS cached_first_orders")
  db.Exec("CREATE TEMP TABLE cached_first_orders AS SELECT user_id, min(created_at) AS created_at FROM orders o WHERE o.user_id IN (" + userIds + ") GROUP BY o.user_id")
}
