package db

import (
	"os"
	"testing"
)

func TestDbServiceInit(t *testing.T) {
	path := "./db_test.txt"
	service := Init(path)
	data := service.Get()
	t.Run("Empty Creation", func(t *testing.T) {
		if len(data) != 0 {
			t.Error("Error creating db Service")
		}
	})

	t.Run("Diff ids", func(t *testing.T) {
		service.Insert("1", 34)
		service.Insert("2", 54)
		data = service.Get()
		if data[0].Id == data[1].Id {
			t.Error("Error with Ids: id1 {}, id2 {}", data[0].Id, data[1].Id)
		}
	})

	service.DeleteDb()
	_, err := os.ReadFile(path)
	if err == nil {
		t.Error("Error, file still exitst", path)
	}
}

func TestCrudDB(t *testing.T) {
	path := "./db_test_crud.txt"
	service := Init(path)

	t.Run("Creation", func(t *testing.T) {
		service.Insert("1", 34)
		service.Insert("2", 54)
		data := service.Get()
		if len(data) != 2 {
			t.Error("Error with Ids: id1 {}, id2 {}", data[0].Id, data[1].Id)
		}
	})

	t.Run("Update", func(t *testing.T) {
		newData := service.Get()
		oldValue := newData[0]
		service.Update(newData[0].Id, "name", 60)
		if oldValue.Name == newData[0].Name ||
			oldValue.Value == newData[0].Value {
			t.Error("Error updating {}, {}", newData[0], newData[0])
		}
	})

	t.Run("Delete", func(t *testing.T) {
		newData := service.Get()
		service.Delete(newData[0].Id)
		newData2 := service.Get()
		if len(newData) == len(newData2) {
			t.Error("Error deleting {}, {}", len(newData), len(newData2))
		}
	})

	service.DeleteDb()
	_, err := os.ReadFile(path)
	if err == nil {
		t.Error("Error, file still exitst", path)
	}
}

func TestSaveFile(t *testing.T) {
	path := "./db_test_save.txt"
	t.Run("Saving File", func(t *testing.T) {
		service := Init(path)

		service.Insert("Name1", 49)
		service.Insert("Name2", 70)
		service.SaveFile()

		file, err := os.Open(path)
		if err != nil {
			t.Error("Error saving file ", err)
		}
		defer file.Close()
	})

	t.Run("Reopen and save", func(t *testing.T) {
		service := Init(path)
		data := service.Get()[0]
		service.Update(data.Id, "NewName1", data.Value)
		service.SaveFile()

		service = Init(path)
		allData := service.Get()
		newData := allData[0]

		if data.Name == newData.Name {
			t.Error("Error on files ", data, newData)
		}

		if len(*service.data) > 2 {
			t.Error("Error updating file there are {} rows", len(*service.data))
		}
		service.DeleteDb()
	})
}
