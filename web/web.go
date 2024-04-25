package web

import (
	"crypto/rand"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vilmibm/trunkless/db"
)

// TODO detect max id on startup

type source struct {
	ID   int64
	Name string
}

type phrase struct {
	ID     int64
	Text   string
	Source source
}

func Serve() error {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})
	r.LoadHTMLFiles("templates/index.tmpl")
	r.StaticFile("/cutive.ttf", "./web/assets/cutive.ttf")
	r.StaticFile("/favicon.ico", "./web/assets/favicon.ico")
	r.StaticFile("/bg_light.gif", "./web/assets/bg_light.gif")
	r.StaticFile("/bg_dark.gif", "./web/assets/bg_dark.gif")
	r.StaticFile("/main.js", "./web/assets/main.js")
	r.StaticFile("/html2canvas.min.js", "./web/assets/html2canvas.min.js")

	randMax := big.NewInt(db.MaxID)

	r.HEAD("/", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", struct {
			MaxID int
			// TODO anything else?
		}{db.MaxID})
	})

	r.GET("/line", func(c *gin.Context) {
		db, err := db.Connect()
		if err != nil {
			log.Println(err.Error())
			c.String(http.StatusInternalServerError, "oh no.")
			return
		}
		defer db.Close()

		id, err := rand.Int(rand.Reader, randMax)
		if err != nil {
			log.Println(err.Error())
			c.String(http.StatusInternalServerError, "oh no.")
			return
		}

		stmt, err := db.Prepare("select p.phrase, p.id, s.name from phrases p join sources s on p.sourceid = s.id where p.id = ?")
		if err != nil {
			log.Println(err.Error())
			c.String(http.StatusInternalServerError, "oh no.")
			return
		}

		row := stmt.QueryRow(id.Int64())
		var p phrase
		var s source
		err = row.Scan(&p.Text, &s.ID, &s.Name)
		if err != nil {
			log.Println(err.Error())
			c.String(http.StatusInternalServerError, "oh no.")
		}
		p.Source = s
		p.ID = id.Int64()
		c.JSON(http.StatusOK, p)
	})

	return r.Run() // 8080
}
