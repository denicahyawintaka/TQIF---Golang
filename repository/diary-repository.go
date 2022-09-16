package repository

import (
	"log"
	"tqif-golang/entity"

	"gorm.io/gorm"
)

// DiaryRepository is a ....
type DiaryRepository interface {
	InsertDiary(b entity.Diary) entity.Diary
	UpdateDiary(b entity.Diary) entity.Diary
	DeleteDiary(b entity.Diary)
	AllDiary(userID string) []entity.Diary
	FindDiaryByID(diaryID uint64) entity.Diary
}

type diaryConnection struct {
	connection *gorm.DB
}

// NewDiaryRepository creates an instance DiaryRepository
func NewDiaryRepository(dbConn *gorm.DB) DiaryRepository {
	return &diaryConnection{
		connection: dbConn,
	}
}

func (db *diaryConnection) InsertDiary(b entity.Diary) entity.Diary {
	db.connection.Save(&b)
	db.connection.Preload("User").Find(&b)
	return b
}

func (db *diaryConnection) UpdateDiary(b entity.Diary) entity.Diary {
	db.connection.Save(&b)
	db.connection.Preload("User").Find(&b)
	return b
}

func (db *diaryConnection) DeleteDiary(b entity.Diary) {
	db.connection.Delete(&b)
}

func (db *diaryConnection) FindDiaryByID(diaryID uint64) entity.Diary {
	var diary entity.Diary
	db.connection.Preload("User").Find(&diary, diaryID)
	return diary
}

func (db *diaryConnection) AllDiary(userID string) []entity.Diary {
	var diaries []entity.Diary
	log.Println("datanya ", userID)

	db.connection.Preload("User").Where("user_id = ?", userID).Find(&diaries)
	return diaries
}
