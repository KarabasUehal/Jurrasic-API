package storage

import (
	"Dinosaurus/models"
	"log"
)

func GetAllDinosaurus() *[]models.Dinosaurus {
	var dinosaurus *[]models.Dinosaurus

	if result := DB.Find(&dinosaurus); result.Error != nil {
		log.Printf("Failed to query dinosaurus: %v", result.Error)
		return nil
	}
	return dinosaurus
}

func GetDinosaurByID(id int) *models.Dinosaurus {
	var dinosaur models.Dinosaurus
	if result := DB.First(&dinosaur, id); result.Error != nil {
		log.Printf("Failed to query item ID %d: %v", id, result.Error)
		return nil
	}
	return &dinosaur
}

func AddDinosaur(newDino *models.Dinosaurus) *models.Dinosaurus {
	var dinosaurus []models.Dinosaurus

	if err := DB.Create(&newDino); err != nil {
		log.Printf("Failed to create dinosaur: %v", err)
	}
	DB.Find(&dinosaurus)
	return newDino
}

func UpdateDinosaurByID(id int, updateDino models.Dinosaurus) *models.Dinosaurus {
	var dinosaur models.Dinosaurus
	if result := DB.First(&dinosaur, id); result.Error != nil {
		log.Printf("Failed to find dinosaur ID %d: %v", updateDino.ID, result.Error)
		return nil
	}

	dinosaur.Species = updateDino.Species
	dinosaur.Types = updateDino.Types
	dinosaur.Height = updateDino.Height
	dinosaur.Length = updateDino.Length
	dinosaur.Weight = updateDino.Weight
	dinosaur.Aquatic = updateDino.Aquatic
	dinosaur.Flying = updateDino.Flying

	if err := DB.Save(dinosaur); err != nil {
		log.Printf("Failed to update dinosaur ID %d: %v", id, err)
	}
	return &dinosaur
}

func DeleteDinosaurByID(id int) bool {

	if result := DB.Delete(&models.Dinosaurus{}, id); result.Error != nil {
		log.Printf("Failed to delete dinosaur ID %d: %v", id, result.Error)
		return false
	}
	return true
}
