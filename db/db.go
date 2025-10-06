package db

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/DNLSalazar/gocounter/models"
)

type DatabaseService struct {
	path      string
	separator string
	data      *[]models.Counter
}

func Init(path string) *DatabaseService {
	service := &DatabaseService{
		path:      path,
		separator: "\r\n",
	}
	service.getData()
	return service
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("Error", err)
	}
}

func counterToString(c models.Counter) string {
	return fmt.Sprintf("%v;%v;%v", c.Id, c.Name, c.Value)
}

func counterFromString(str string) models.Counter {
	values := strings.Split(str, ";")
	id, err := strconv.ParseInt(strings.TrimSpace(values[0]), 10, 64)
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

func (d *DatabaseService) createFile() {
	dir := filepath.Dir(d.path)
	if dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("Error creating directories for db", err)
			panic(err)
		}
	}

	fileName := filepath.Base(d.path)

	if fileName == "" {
		panic("Invalid file name for db")
	}

	file, err := os.Create(d.path)
	if err != nil {
		log.Fatal("Error creating file ", err, d.path)
	}
	defer file.Close()
}

func (d *DatabaseService) readFile() []models.Counter {
	content, err := os.ReadFile(d.path)
	if err != nil {
		switch err.(type) {
		case *fs.PathError:
			d.createFile()
			return []models.Counter{}
		}
		log.Fatal("Error reading db file ", err)
	}

	str := strings.Split(string(content), d.separator)
	return parseContent(str)
}

func (d *DatabaseService) SaveFile() {
	str := make([]string, len(*d.data))
	for i, v := range *d.data {
		str[i] = counterToString(v)
	}

	data := []byte(strings.Join(str, d.separator))

	tempPath := fmt.Sprintf("%v.temp", d.path)

	temFile, err := os.Create(tempPath)

	if err != nil {
		log.Printf("Trying to create temp file %v counters\r\n", len(*d.data))
		log.Fatal("Error creating temp file for writing new data", err)
	}

	_, err = temFile.Write(data)
	if err != nil {
		log.Printf("Trying to write temp file %v counters\r\n", len(*d.data))
		log.Fatal("Error writting temp file", err)
	}

	err = os.Remove(d.path)
	if err != nil {
		log.Printf("Trying to delete file %v counters\r\n", len(*d.data))
		log.Fatal("Error deleting file for writing new data", err)
	}

	file, err := os.Create(d.path)
	if err != nil {
		log.Printf("Trying to open file %v counters\r\n", len(*d.data))
		log.Fatal("Error opening file for writing file", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Printf("Trying to write %v counters\r\n", len(*d.data))
		log.Fatal("Error writting file", err)
	}

	temFile.Close()
	err = os.Remove(tempPath)
	if err != nil {
		log.Printf("Trying to delete temp file %v counters\r\n", len(*d.data))
		log.Fatal("Error deleting temp file for writing new data", err)
	}
}

func (d *DatabaseService) getData() {
	content := d.readFile()
	d.data = &content
}

func (d *DatabaseService) Get() []models.Counter {
	return *d.data
}

func (d *DatabaseService) Insert(name string, value int) error {
	id := time.Now().UnixNano()
	newCounter := models.CreateCounter(id, value, name)
	(*(d.data)) = append((*(d.data)), newCounter)
	time.Sleep(time.Nanosecond)
	return nil
}

func (d *DatabaseService) Update(id int64, name string, value int) error {
	for i, v := range *d.data {
		if v.Id == id {
			(*(d.data))[i].Name = name
			(*(d.data))[i].Value = value
		}
	}
	return nil
}

func (d *DatabaseService) Delete(id int64) error {
	for i, v := range *d.data {
		if v.Id == id {
			*(d.data) = slices.Delete(*d.data, i, i+1)
		}
	}
	return nil
}

func (d *DatabaseService) DeleteDb() {
	d.data = &[]models.Counter{}
	err := os.Remove(d.path)
	if err != nil {
		log.Fatal("Error deleting db", err)
		return
	}
}
