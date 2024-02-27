package service

import (
	"fmt"
	"meepShopTest/internal/apierr"
	"meepShopTest/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SaveMoneyRequestStruct struct {
	NewDeposit float64 `json:"newDeposit" binding:"required"`
}

type WithdrawMoneyRequestStruct struct {
	NewDeposit float64 `json:"newDeposit" binding:"required"`
}

type TransferMoneyStruct struct {
	NewDeposit           float64 `json:"newDeposit" binding:"required"`
	TransferOutUserIdStr string  `json:"transferOutUserIdStr" binding:"required"`
	TransferInUserIdStr  string  `json:"transferInUserIdStr" binding:"required"`
	TransferOutUserId    uuid.UUID
	TransferInUserId     uuid.UUID
}

func (svc *Service) SaveMoney(userId uuid.UUID, newDeposit float64) error {
	user, err := svc.db.GetUserByID(svc.db.DB, userId)
	if err != nil {
		return err
	}
	updatedUserDeposit, err := CheckTotalDepositIsGreaterThanZero(user.Deposit, newDeposit)
	if err != nil {
		return err
	}

	updateUserDepositStruct := &database.UpdateUserDepositStruct{
		Tx:                 svc.db.DB,
		User:               user,
		UpdatedUserDeposit: updatedUserDeposit,
	}

	err = svc.db.UpdateUserDeposit(updateUserDepositStruct)
	if err != nil {
		return err
	}
	return nil
}

func (svc *Service) WithdrawMoney(userId uuid.UUID, newDeposit float64) error {
	user, err := svc.db.GetUserByID(svc.db.DB, userId)
	if err != nil {
		return err
	}
	// 領錢用負數的概念去看存錢這件事情
	updatedUserDeposit, err := CheckTotalDepositIsGreaterThanZero(user.Deposit, -newDeposit)
	if err != nil {
		return err
	}

	updateUserDepositStruct := &database.UpdateUserDepositStruct{
		Tx:                 svc.db.DB,
		User:               user,
		UpdatedUserDeposit: updatedUserDeposit,
	}

	err = svc.db.UpdateUserDeposit(updateUserDepositStruct)
	if err != nil {
		return err
	}
	return nil
}

// 個人的習慣是傳入參數超過三個就用struct包起來 看起來比較不會很混亂
func (svc *Service) TransferMoney(transferMoneyStruct TransferMoneyStruct) error {
	newDeposit := transferMoneyStruct.NewDeposit
	transferOutUserId := transferMoneyStruct.TransferOutUserId
	transferInUserId := transferMoneyStruct.TransferInUserId

	tx := svc.db.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err := tx.Transaction(func(tx *gorm.DB) error {
		// 先看轉帳人的錢轉完後會不會變負數
		transferOutUser, err := svc.db.GetUserByID(tx, transferOutUserId)
		if err != nil {
			return err
		}

		// 領錢用負數的概念去看存錢這件事情
		transferOutUpdatedUserDeposit, err := CheckTotalDepositIsGreaterThanZero(transferOutUser.Deposit, -newDeposit)
		if err != nil {
			return err
		}

		transferOutUpdateUserDepositStruct := &database.UpdateUserDepositStruct{
			Tx:                 tx,
			User:               transferOutUser,
			UpdatedUserDeposit: transferOutUpdatedUserDeposit,
		}

		err = svc.db.UpdateUserDeposit(transferOutUpdateUserDepositStruct)
		if err != nil {
			return err
		}

		// 再看收到錢的人收到後會不會變負數
		transferInUser, err := svc.db.GetUserByID(tx, transferInUserId)

		if err != nil {
			return err
		}

		transferInUpdatedUserDeposit, err := CheckTotalDepositIsGreaterThanZero(transferInUser.Deposit, newDeposit)
		if err != nil {
			return err
		}

		transferInUpdateUserDepositStruct := &database.UpdateUserDepositStruct{
			Tx:                 tx,
			User:               transferInUser,
			UpdatedUserDeposit: transferInUpdatedUserDeposit,
		}

		err = svc.db.UpdateUserDeposit(transferInUpdateUserDepositStruct)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 檢查user計算前和計算後的存款會不會是負數 (基本上這邊就是不管領錢,存錢還是轉帳都會進行檢查,以防有人去db亂改資料或是用其他未知的api達到存款是負數的情形發生,這邊也不針對極端案例做額外的處理,一視同仁)
// 在這邊把相加完的結果 算完後順便回傳 讓call他的人不用再算一次
func CheckTotalDepositIsGreaterThanZero(userDeposit, newDeposit float64) (float64, error) {
	totalDeposit := userDeposit + newDeposit
	if userDeposit < 0 || totalDeposit < 0 {
		return 0, fmt.Errorf(apierr.ErrTotalDepositIsNegativeNumber.Message)
	}
	return totalDeposit, nil
}
