package database

import (
	"fmt"
	"meepShopTest/internal/apierr"
	"meepShopTest/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateUserDepositStruct struct {
	Tx                 *gorm.DB
	User               *model.User
	UpdatedUserDeposit float64
}

func (d *GormDatabase) GetUserByID(tx *gorm.DB, userId uuid.UUID) (*model.User, error) {
	user := model.User{}
	err := tx.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf(apierr.ErrNoData.Message)
	}
	return &user, nil
}

func (d *GormDatabase) UpdateUserDeposit(updateUserDepositStruct *UpdateUserDepositStruct) error {
	tx := updateUserDepositStruct.Tx
	user := updateUserDepositStruct.User
	updatedUserDeposit := updateUserDepositStruct.UpdatedUserDeposit
	result := tx.Model(user).Update("deposit", updatedUserDeposit)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf(apierr.ErrNoData.Message)
	}

	return nil
	// 其實也可以把計算都丟給 sql來決定有沒有負數 有點類似下面的作法 但我比較喜歡用程式本身來做判斷,也比較好寫test來做測試
	// 	result := tx.Exec(`
	// 	UPDATE users
	// 	SET deposit = deposit + ?
	// 	WHERE id = ?
	// 	AND deposit + ? >= 0;
	// `, newDeposit, user.ID, updatedUserDeposit)

}

func (d *GormDatabase) CreateUser() error {
	userId, err := uuid.Parse("6073a4b0711d160028829eb1")
	if err != nil {
		return err
	}
	user := &model.User{
		ID:            userId,
		FirstName:     "yanxian",
		LastName:      "li",
		DisplayName:   "yanxianli",
		AccountNumber: "0987654321",
		Deposit:       200,
	}
	return d.DB.Create(user).Error
}
