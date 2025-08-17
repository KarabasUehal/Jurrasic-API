package storage

import (
	"golang-gin/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Функция инициализации базы данных
func InitDatabase() error {
	dsn := "host=localhost user=postgres password=1998ki31 dbname=api port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	err = db.AutoMigrate(&models.Dinosaurus{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	var count int64
	db.Model(&models.Dinosaurus{}).Count(&count)
	if count == 0 {
		dinosaurus := []models.Dinosaurus{
			{Species: "Allosaurus", Types: "Carnivorous", Height: 4.00, Length: 11.00, Weight: 4.00, Aquatic: false, Flying: false},
			{Species: "Argentinosaurus", Types: "Herbivorous", Height: 7.00, Length: 33.00, Weight: 100.00, Aquatic: false, Flying: false},
			{Species: "Brachiosaurus", Types: "Herbivorous", Height: 10.00, Length: 20.00, Weight: 50.00, Aquatic: false, Flying: false},
			{Species: "Carnotaurus", Types: "Carnivorous", Height: 3.00, Length: 8.00, Weight: 3.00, Aquatic: false, Flying: false},
			{Species: "Diplodocus", Types: "Herbivorous", Height: 6.00, Length: 35.00, Weight: 100.00, Aquatic: false, Flying: false},
			{Species: "Giganotosaurus", Types: "Carnivorous", Height: 5.50, Length: 14.00, Weight: 14.00, Aquatic: false, Flying: false},
			{Species: "Hatzegopteryx", Types: "Carnivorous", Height: 3.00, Length: 12.00, Weight: 0.25, Aquatic: false, Flying: true},
			{Species: "Ichthyotitan severnensis", Types: "Carnivorous", Height: 6.00, Length: 25.00, Weight: 34.00, Aquatic: true, Flying: false},
			{Species: "Mosasaurus", Types: "Carnivorous", Height: 4.00, Length: 17.00, Weight: 15.00, Aquatic: true, Flying: false},
			{Species: "Parasaurolophus", Types: "Herbivorous", Height: 5.00, Length: 10.00, Weight: 2.70, Aquatic: false, Flying: false},
			{Species: "Quetzalcoatlus", Types: "Carnivorous", Height: 3.00, Length: 11.00, Weight: 0.40, Aquatic: false, Flying: true},
			{Species: "Shonisaurus sikanniensis", Types: "Carnivorous", Height: 5.00, Length: 21.00, Weight: 29.00, Aquatic: true, Flying: false},
			{Species: "Spinosaurus", Types: "Carnivorous", Height: 4.00, Length: 15.00, Weight: 6.00, Aquatic: true, Flying: false},
			{Species: "Triceratops", Types: "Herbivorous", Height: 3.00, Length: 12.00, Weight: 9.00, Aquatic: false, Flying: false},
			{Species: "Tyrannosaurus", Types: "Carnivorous", Height: 5.00, Length: 12.80, Weight: 8.50, Aquatic: false, Flying: false},
		}
		if err := db.Create(&dinosaurus).Error; err != nil {
			log.Fatalf("Failed to insert initial data: %v", err)
		}
	}
	return nil
}

// Функция для получения экземпляра базы данных
func GetDB() *gorm.DB {
	// Возвращение глобальной переменной db, содержащей подключение к базе данных
	return db
}
