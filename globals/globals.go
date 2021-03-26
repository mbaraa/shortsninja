package globals

import (
	"github.com/baraa-almasri/shortsninja/config"
	"github.com/baraa-almasri/shortsninja/db"
	"html/template"
)

var (
	Config    = config.LoadConfig()
	DBManager db.Database
	Templates *template.Template
)
