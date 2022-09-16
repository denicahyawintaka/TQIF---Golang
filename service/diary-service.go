package service

import (
	"fmt"
	"log"

	"tqif-golang/dto"
	"tqif-golang/entity"
	"tqif-golang/repository"

	"github.com/mashingan/smapping"
)

// DiaryService is a ....
type DiaryService interface {
	Insert(b dto.DiaryCreateDTO) entity.Diary
	Update(b dto.DiaryUpdateDTO) entity.Diary
	Delete(b entity.Diary)
	All(userID string) []entity.Diary
	FindByID(diaryID uint64) entity.Diary
	IsAllowedToEdit(userID string, diaryID uint64) bool
}

type diaryService struct {
	diaryRepository repository.DiaryRepository
}

// NewDiaryService .....
func NewDiaryService(diaryRepo repository.DiaryRepository) DiaryService {
	return &diaryService{
		diaryRepository: diaryRepo,
	}
}

func (service *diaryService) Insert(b dto.DiaryCreateDTO) entity.Diary {
	diary := entity.Diary{}
	err := smapping.FillStruct(&diary, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.diaryRepository.InsertDiary(diary)
	return res
}

func (service *diaryService) Update(b dto.DiaryUpdateDTO) entity.Diary {
	diary := entity.Diary{}
	err := smapping.FillStruct(&diary, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.diaryRepository.UpdateDiary(diary)
	return res
}

func (service *diaryService) Delete(b entity.Diary) {
	service.diaryRepository.DeleteDiary(b)
}

func (service *diaryService) All(userID string) []entity.Diary {
	return service.diaryRepository.AllDiary(userID)
}

func (service *diaryService) FindByID(diaryID uint64) entity.Diary {
	return service.diaryRepository.FindDiaryByID(diaryID)
}

func (service *diaryService) IsAllowedToEdit(userID string, diaryID uint64) bool {
	b := service.diaryRepository.FindDiaryByID(diaryID)
	id := fmt.Sprintf("%v", b.UserID)
	return userID == id
}
