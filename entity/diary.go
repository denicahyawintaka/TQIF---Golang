package entity

// Diary struct represents diaries table in database
type Diary struct {
	ID       uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Body     string `gorm:"type:text" json:"body"`
	Mood     string `gorm:"type:text" json:"mood"`
	Category string `gorm:"type:text" json:"category"`
	Date     string `gorm:"type:text" json:"date"`
	UserID   uint64 `gorm:"not null" json:"-"`
	User     User   `gorm:"foreignkey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
}
