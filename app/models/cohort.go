package models

import (
	"fmt"
	"strconv"
	"time"
)

type Cohort struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Users []User    `json:"users"`
}

func (c *Cohort) String() string {
	return fmt.Sprintf("Start=%s / End=%s / Len=%d", c.Start.String(), c.End.String(), len(c.Users))
}

func (c *Cohort) UserIds() []string {
	ids := []string{}
	for _, user := range c.Users {
		ids = append(ids, strconv.FormatInt(user.Id, 10))
	}
	return ids
}
