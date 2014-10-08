package models

import (
	"fmt"
	"time"
)

type Order struct {
	Id      int64     `db:"id"`
	Num     int64     `db:"order_num"`
	User_id int64     `db:"user_id"`
	Created time.Time `db:"created_at"`
	Updated time.Time `db:"updated_at"`
}

func (o Order) String() string {
	return fmt.Sprintf("Order %d: num=%d / user_id=%d / created=%s / updated=%s", o.Id, o.Num, o.User_id, o.Created.String(), o.Updated.String())
}
