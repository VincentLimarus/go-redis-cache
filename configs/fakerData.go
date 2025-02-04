package configs

import (
	models "VincentLimarus/go-redis/models/database"
	"fmt"

	"github.com/jaswdr/faker"
)

func Faker() {
    db := GetDB()
    fake := faker.New()

    var students []models.Students
    for i := 0; i < 100000; i++ {
        student := models.Students{
            Name:    fake.Person().Name(),
            Email:   fake.Internet().Email(),
            Address: fake.Address().Address(),
        }
        students = append(students, student)

        if len(students) >= 1000 {
            if err := db.Create(&students).Error; err != nil {
                fmt.Println("Error inserting batch:", err)
            } else {
                fmt.Println("Inserted batch of students")
            }
            students = nil 
        }
    }

    if len(students) > 0 { 
        db.Create(&students)
    }
}

