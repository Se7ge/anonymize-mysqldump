package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syreclabs.com/go/faker"
	"text/template"

	"github.com/DekodeInteraktiv/anonymize-mysqldump/internal/embed"
)

type Config struct {
	Name        string
	ProcessName string
	Version     string
	Commit      string
	Date        string
	WD          string
	Patterns    []ConfigPattern `json:"patterns"`
}

type ConfigPattern struct {
	TableName      string         `json:"tableName"`
	TableNameRegex string         `json:"tableNameRegex"`
	Purge          bool           `json:"purge"`
	Fields         []PatternField `json:"fields"`
}

type PatternField struct {
	Field       string                   `json:"field"`
	Position    int                      `json:"position"`
	Type        string                   `json:"type"`
	Template    Template                 `json:"template"`
	Constraints []PatternFieldConstraint `json:"constraints"`
}

type PatternFieldConstraint struct {
	Field    string `json:"field"`
	Position int    `json:"position"`
	Value    string `json:"value"`
}

type Template struct {
	Tpl *template.Template
}

var tplFuncMap = template.FuncMap{
	"fakerInternet": faker.Internet,
	"fakerUser":     faker.Internet().UserName,
	"fakerAddress":  faker.Address().String,
}

func (t *Template) UnmarshalJSON(data []byte) error {
	var v string

	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Printf("failed unmarshaling response | %s", err.Error())
		return err
	}
	if data != nil {
		t.Tpl = template.Must(template.New("template").Funcs(tplFuncMap).Parse(v))
	}
	return nil
}

// New creates a new Config from flags and environment variables
func New(version, commit, date string) *Config {
	c := &Config{
		Name:        "Anomymize MySQLDump",
		ProcessName: filepath.Base(os.Args[0]),
		Version:     version,
		Commit:      commit,
		Date:        date,
	}

	// Get Working Dir
	wd, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	c.WD = wd

	return c
}

// ParseConfig parses a default or user provided config file.
func (c *Config) ParseConfig(filepath string) {
	var jsonConfig []byte
	var err error

	jsonConfig = []byte(embed.DefaultConfig)

	if filepath != "" {
		jsonConfig, err = ioutil.ReadFile(filepath)
		if err != nil {
			log.Fatalf("Failed reading config file: %s", err)
		}
	}

	err = json.Unmarshal(jsonConfig, &c)

	// Make sure the JSON read is valid.
	if err != nil {
		log.Fatalf("JSON file not valid!")
	}
}
