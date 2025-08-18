package storage

import (
	"golang-gin/models"
	"log"
)

func GetAllDinosaurus() *[]models.Dinosaurus {
	var dinosaurus *[]models.Dinosaurus

	if result := db.Find(&dinosaurus); result.Error != nil {
		log.Printf("Failed to query dinosaurus: %v", result.Error)
		return nil
	}
	return dinosaurus
}

func GetDinosaurByID(id int) *models.Dinosaurus {
	var dinosaur models.Dinosaurus
	if result := db.First(&dinosaur, id); result.Error != nil {
		log.Printf("Failed to query item ID %d: %v", id, result.Error)
		return nil
	}
	return &dinosaur
}

// Функция для добавление новой записи о динозавре
func AddDinosaur(newDino models.Dinosaurus) models.Dinosaurus {
	var dinosaurus *[]models.Dinosaurus
	// Использование GORM для выполнения SQL-запроса INSERT и сохранения нового динозавра в базе данных
	if err := db.Create(&newDino); err != nil {
		log.Printf("Failed to update item ID %d: %v", newDino.ID, err)
	}
	db.Find(&dinosaurus)
	return newDino // Возвращение созданного динозавра
}

// Функция для обновления существующего динозавра по ID
func UpdateDinosaurByID(id int, updateDino models.Dinosaurus) *models.Dinosaurus {
	var dinosaur models.Dinosaurus
	// Использование GORM для выполнения SQL-запроса SELECT с условием WHERE id = id динозавра
	if result := db.First(&dinosaur, id); result.Error != nil {
		log.Printf("Failed to find item ID %d: %v", updateDino.ID, result.Error)
		return nil
	}

	dinosaur.Species = updateDino.Species
	dinosaur.Types = updateDino.Types
	dinosaur.Height = updateDino.Height
	dinosaur.Length = updateDino.Length
	dinosaur.Weight = updateDino.Weight
	dinosaur.Aquatic = updateDino.Aquatic
	dinosaur.Flying = updateDino.Flying
	// Использование GORM для выполнения SQL-запроса UPDATE и сохранения обновленной записи о динозавре в базе данных
	if err := db.Save(dinosaur); err != nil {
		log.Printf("Failed to update item ID %d: %v", id, err)
	}
	return &dinosaur // Возвращение обновленной записи
}

// Функция для удаления динозавра по ID
func DeleteDinosaurByID(id int) bool {
	// Использование GORM для выполнения SQL-запроса DELETE с условием WHERE id = id заметки
	if result := db.Delete(&models.Dinosaurus{}, id); result.Error != nil {
		log.Printf("Failed to delete item ID %d: %v", id, result.Error)
		return false // Возвращение false, если удаление не удалось
	}
	return true // Возвращение true при успешном удалении записи
}
