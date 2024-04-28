package web

import (
	"context"
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
	ID   string
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

	// TODO use new db functions for id ranges
	randMax := big.NewInt(db.MaxID)

	r.HEAD("/", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", struct {
			// TODO handle multiple corpora
			MaxID int
		}{db.MaxID})
	})

	// TODO retool this for pg

	r.GET("/line", func(c *gin.Context) {
		conn, err := db.Connect()
		if err != nil {
			log.Println(err.Error())
			c.String(http.StatusInternalServerError, "oh no.")
			return
		}
		defer conn.Close(context.Background())

		id, err := rand.Int(rand.Reader, randMax)
		if err != nil {
			log.Println(err.Error())
			c.String(http.StatusInternalServerError, "oh no.")
			return
		}

		row := conn.QueryRow(context.Background(), "select p.phrase, s.id, s.name from phrases p join sources s on p.sourceid = s.id where p.id = $1", id.Int64())
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
