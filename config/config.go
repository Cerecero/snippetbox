package config

import (
	"log"

	"github.com/Cerecero/snippetbox/internal/models"
)

type Application struct{
	ErrorLog *log.Logger
	InfoLog *log.Logger
	Snippets *models.SnippetModel
}
