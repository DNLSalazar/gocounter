package db

import (
	"fmt"
	"goCounter/models"
	"log"
	"os"
	"strconv"
	"strings"
)

type DatabaseService struct {
	path      string
	separator string
}

func Init(path string) *DatabaseService {
	return &DatabaseService{
		path:      path,
		separator: "---",
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("Error", err)
	}
}

func counterToString(c models.Counter) string {
  return fmt.Sprintf("%v;%v;%v\r\n", c.Id, c.Name, c.Value)
}

func counterFromString(str string) models.Counter {
	values := strings.Split(str, ";")
	id, err := strconv.Atoi(strings.TrimSpace(values[0]))
	checkErr(err)
	name := strings.TrimSpace(values[1])
	value, err := strconv.Atoi(strings.TrimSuffix(values[2], "\r\n"))
	checkErr(err)
	return models.CreateCounter(id, value, name)
}

func parseContent(str []string) []models.Counter {
	counters := []models.Counter{}
	for _, v := range str {
		if v != "" {
			counters = append(counters, counterFromString(v))
		}
	}
  return counters
}

func (d *DatabaseService) readFile() []models.Counter {
	content, err := os.ReadFile(d.path)
	if err != nil {
		if err == os.ErrNotExist {
			_, err := os.Create(d.path)
			if err != nil {
				log.Fatal("Error creating file", err, d.path)
			}
			return []models.Counter{}
		} else {
			log.Fatal("Error reading file", err, d.path)
		}
	}

	str := strings.Split(string(content), d.separator)
	return parseContent(str)
}

func (d *DatabaseService) saveFile() {

}

func (d *DatabaseService) ReadDb() []models.Counter {
  return d.readFile()
}

func (d *DatabaseService) Insert(c models.Counter) error {

	return nil
}

func (d *DatabaseService) Update(c models.Counter) error {
	return nil
}

func (d *DatabaseService) Create(c models.Counter) error {
	return nil
}
