package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User
type User struct {
	ID            uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt     time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     *time.Time
	DeletedAt     gorm.DeletedAt
	FirstName     string
	LastName      string
	DisplayName   string
	AccountNumber string `gorm:"uniqueIndex"` // 原則上應該還會有一個帳戶,帳戶應該也是唯一的值
	// 其實這裡也可以用這個方式阻擋有人刻意改db的值為負數 `gorm:"check:deposit > 0"` 但我覺得有點隱晦就沒這麼做
	Deposit float64 // 這邊用浮點數,是為了保留有角或是分的可能性
}
