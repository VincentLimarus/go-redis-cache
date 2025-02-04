package configs

import (
	models "VincentLimarus/go-redis/models/database"
	"fmt"

	"github.com/jaswdr/faker"
)

func Faker() {
	db := GetDB()

	fake := faker.New()

	for i := 0; i < 10000; i++ {
		student := models.Students{
			Name:    fake.Person().Name(),
			Email:   fake.Internet().Email(),
			Address: fake.Address().Address(),
		}

		result := db.Create(&student)
		if result.Error != nil{
			fmt.Println("Error inserting student:", result.Error)
		} else{
			fmt.Println("Inserted student:", student.Name)
		}
	}
}
