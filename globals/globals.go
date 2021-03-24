package globals

import (
	"github.com/baraa-almasri/shortsninja/db"
	"html/template"
	"os"
)

var (
	IPInfoToken = os.Getenv("IP_INFO_IO_TOKEN")
	DBManager   db.Database
	Templates   *template.Template
)
