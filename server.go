package main

import (
  "fmt"
  "github.com/codegangsta/martini-contrib/binding"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "html/template"
  "instago/app"
  . "instago/app/models"
  "instago/config"
  "os"
  "strings"
  "time"
  "strconv"
  //"runtime"
)

var dayDuration time.Duration = time.Hour * 24

type CohortDetails struct {
  Start    float64 `json:"start"`
  End      float64 `json:"end"`
  Orderers int64   `json:"orderers"`
  Firsts   int64   `json:"firsts"`
}

func (c CohortDetails) String() string {
  return fmt.Sprintf("start=%f / end=%f / orderers=%f / firsts=%f", c.Start, c.End, c.Orderers, c.Firsts)
}

type CohortBlob struct {
  Start   time.Time       `json:"start"`
  End     time.Time       `json:"end"`
  Users   int             `json:"users"`
  Details []CohortDetails `json:"details"`
}

func analyzeCohort(db *app.Database, c *Cohort, end time.Time, d time.Duration) ([]CohortDetails, error) {
  details := []CohortDetails{}
  userIds := strings.Join(c.UserIds(), ",")
  db.CacheFirstOrders(userIds)
  for t,i := c.Start, 0; t.Before(end); t,i = t.Add(d), i +1 {
    startNumDays := int64(time.Duration(i) * d / dayDuration)
    endNumDays := int64(time.Duration(i+1) * d / dayDuration)
    queryOrderers := `SELECT COUNT(DISTINCT(u.id)) FROM users u
                      JOIN ORDERS o ON o.user_id = u.id
                      WHERE u.id IN (` + userIds + `)
                      AND o.created_at >= DATETIME(u.created_at, '+` + strconv.FormatInt(startNumDays,10) +` days')
                      AND o.created_at < DATETIME(u.created_at, '+` + strconv.FormatInt(endNumDays,10) +` days')`
    countOrderers, err := db.SelectInt(queryOrderers)
    if err != nil {
      return details, err
    }
    queryFirsts := `SELECT COUNT(DISTINCT(u.id)) FROM users u
                    JOIN cached_first_orders z ON z.user_id = u.id
                    WHERE u.id IN (` + userIds + `)
                    AND z.created_at >= DATETIME(u.created_at, '+` + strconv.FormatInt(startNumDays,10) +` days')
                    AND z.created_at < DATETIME(u.created_at, '+` + strconv.FormatInt(endNumDays,10) +` days')`
    countFirsts, err := db.SelectInt(queryFirsts, t, t.Add(d))
    if err != nil {
      return details, err
    }
    newDetails := CohortDetails{
      t.Sub(c.Start).Hours() / 24,
      t.Add(d).Sub(c.Start).Hours() / 24,
      countOrderers,
      countFirsts,
    }
    details = append(details, newDetails)
  }

  return details, nil
}

type DataRequest struct {
  NumCohorts     int `form:"numCohorts" json:"numCohorts"`
  CohortDuration int `form:"cohortDuration" json:"cohortDuration"`
}

func main() {
  //runtime.GOMAXPROCS(runtime.NumCPU())

  db, err := app.MakeDatabase("db/database.bin")
  if err != nil {
    fmt.Printf("Error opening db: %s\n", err.Error())
    os.Exit(1)
  }

  m := martini.Classic()

  config.Initialize(m)

  m.Use(render.Renderer(render.Options{
    Funcs: []template.FuncMap{
      {
        "heroku": config.IsHeroku,
      },
    },
    Layout: "app",
  }))

  m.Get("/", func(r render.Render) {
    r.HTML(200, "index", nil)
  })

  m.Get("/data", binding.Bind(DataRequest{}), func(r render.Render, req DataRequest) {
    cohortDuration := time.Duration(req.CohortDuration) * dayDuration
    cohorts, _ := db.Cohorts(cohortDuration, cohortDuration*time.Duration(req.NumCohorts))
    endTime := cohorts[len(cohorts)-1].End

    blobs := []CohortBlob{}

    for _, c := range cohorts {
      details, err := analyzeCohort(db, c, endTime, cohortDuration)
      if err != nil {
        r.JSON(500, nil)
      }
      blobs = append(blobs, CohortBlob{c.Start, c.End, len(c.Users), details})
    }
    r.JSON(200, blobs)
  })

  m.Run()
}
