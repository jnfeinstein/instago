package models

import (
  "fmt"
  "time"
)

type User struct {
  Id      int64     `db:"id"`
  Created time.Time `db:"created_at"`
  Updated time.Time `db:"updated_at"`
}

func (u User) String() string {
  return fmt.Sprintf("User %d: created=%s / updated=%s", u.Id, u.Created.String(), u.Updated.String())
}
