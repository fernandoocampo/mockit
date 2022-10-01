package databases

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

var gormDB *gorm.DB
var errUnableToGetPet = errors.New("unable to get pet")

type Friend struct {
	ID   string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name string `gorm:"type:varchar(64);uniqueIndex;not null" json:"name" binding:"required,min=1,max=64"`
}

type Pet struct {
	ID       string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name     string `gorm:"type:varchar(64);uniqueIndex;not null" json:"name" binding:"required,min=1,max=64"`
	Breed    string `gorm:"type:varchar(64);uniqueIndex;not null" json:"breed" binding:"required,min=1,max=64"`
	FriendID string `gorm:"type:uuid;not null" json:"-"`
	Friend   Friend `gorm:"foreignKey:FriendID;joinForeignKey:ID" json:"-"`
}

type PetFriend struct {
	Name   string `json:"name"`
	Breed  string `json:"breed"`
	Friend string `json:"friend"`
}

func GetPet(id string) (*PetFriend, error) {
	var pet PetFriend

	result := gormDB.
		Model(&Pet{}).
		Select("name", "breed").
		Where("id = ?", id).
		Find(&pet)
	if result.Error != nil {
		log.Println(
			"unable to get pet",
			"id", id,
			"error", result.Error)

		return nil, errUnableToGetPet
	}

	return &pet, nil
}

func GetPetAndFriend(id string) (*PetFriend, error) {
	var pet PetFriend

	result := gormDB.
		Model(&Pet{}).
		Select("pets.name as name", "pets.breed as breed", "friends.name as friend").
		Joins("join friends on pets.friend_id = friends.id").
		Where("pets.id = ?", id).
		Find(&pet)
	if result.Error != nil {
		log.Println(
			"unable to get pet",
			"id", id,
			"error", result.Error)

		return nil, errUnableToGetPet
	}

	return &pet, nil
}

func SetORM(orm *gorm.DB) {
	gormDB = orm
}
