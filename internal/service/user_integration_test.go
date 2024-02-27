package service_test

import (
	"context"
	"fmt"
	"math/rand"
	"meepShopTest/config"
	"meepShopTest/internal/apierr"
	"meepShopTest/internal/database"
	"meepShopTest/internal/service"
	"strconv"

	"meepShopTest/internal/model"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func Test_integration_SaveMoney(t *testing.T) {
	// 建立測試資料庫
	testDB, err := setupTestDatabase()
	if err != nil {
		t.Errorf("setupTestDatabase error: %s", err.Error())
		return
	}
	tests := []struct {
		name                      string
		userId                    uuid.UUID
		accountNumber             string
		deposit                   float64
		newDeposit                float64
		unnecessaryCreateUserToDB bool
		expected                  float64
		err                       error
	}{
		// {
		// 這是一個蠻特別的case 是領0元或是負數的case 不過這個會在router層被擋下來 如果想測試的話 可能要寫e2e的測試 或是把那段檢查搬到service層
		// 後續有其他newDeposit為0的case就不列舉了
		// {
		// 	name:                      "empty user",
		// 	userId:                    uuid.MustParse("ad073e32-ff3e-4a0a-9ba8-ca2cc8c59aeb"),
		// 	deposit:                   100,
		// 	newDeposit:                0,
		// 	expected:                  100,
		// 	err:                       fmt.Errorf(apierr.ErrNoData.Message),
		// },
		{
			name:       "origin deposit is zero",
			userId:     uuid.MustParse("f273668e-a960-40f1-8c05-58ecccb90db6"),
			deposit:    0,
			newDeposit: 100,
			expected:   100,
			err:        nil,
		}, {
			name:       "origin deposit is positive number",
			userId:     uuid.MustParse("985a05f1-49f3-4549-9368-8b4c94918081"),
			deposit:    100,
			newDeposit: 100,
			expected:   200,
			err:        nil,
		}, {
			name:       "origin deposit is negative",
			userId:     uuid.MustParse("4d23d1fb-6690-4b65-9724-f4adeed65ec1"),
			deposit:    -100,
			newDeposit: -100,
			expected:   -100,
			err:        fmt.Errorf(apierr.ErrNoData.Message),
		}, {
			name:       "origin deposit is negative",
			userId:     uuid.MustParse("7eab7b04-bc90-4bc4-a564-9004806e6c0b"),
			deposit:    -100,
			newDeposit: 200,
			expected:   -100,
			err:        fmt.Errorf(apierr.ErrNoData.Message),
		},
	}

	// 測試邏輯 1.建立User 2.比對User的存款是否正確 3.進行領錢 4.比對存完錢的存款是否正確
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var accountNumber string
			for i := 0; i < 15; i++ {
				num := rand.Intn(10)
				accountNumber += strconv.Itoa(num)
			}

			testUser := &model.User{
				ID:            tc.userId,
				FirstName:     "user",
				LastName:      "test",
				DisplayName:   "test user",
				AccountNumber: accountNumber,
				Deposit:       tc.deposit,
			}

			// 初始化db(1.建立User)
			err := initializeTestUserData(testDB, testUser)
			if err != nil {
				t.Errorf("initializeTestUserData failed: %v", err)
				return
			}

			defer func() {
				if err := teardownTestDatabase(testDB, tc.userId); err != nil {
					t.Errorf("teardownTestDatabase failed: %v", err)
				}
			}()

			// 建立 Service
			ctx := context.Background()
			svc := service.New(ctx, testDB)

			// 2.比對User的存款是否正確
			user, err := testDB.GetUserByID(testDB.DB, tc.userId)
			if err != nil {
				t.Errorf("first GetUserByID failed: %v", err)
				return
			}
			if user.Deposit != tc.deposit {
				t.Errorf("createdUser's Deposit is incorrect, expected %f, got %f", tc.deposit, user.Deposit)
				return
			}

			// call SaveMoney (3.進行存錢 也是最主要的測試部分)
			err = svc.SaveMoney(tc.userId, tc.newDeposit)
			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("SaveMoney failed: %v", err)
				return
			}
			// 4.比對存完錢的存款是否正確
			user, err = testDB.GetUserByID(testDB.DB, tc.userId)
			if err != nil {
				t.Errorf("second GetUserByID failed: %v", err)
				return
			}
			if user.Deposit != tc.expected {
				t.Errorf("final Deposit is incorrect, expected %f, got %f", tc.expected, user.Deposit)
				return
			}
		})
	}
}

func Test_integration_SaveMoneyToEmptyUser(t *testing.T) {
	// 建立測試資料庫
	testDB, err := setupTestDatabase()
	if err != nil {
		t.Errorf("setupTestDatabase error: %s", err.Error())
		return
	}
	tests := []struct {
		name          string
		userId        uuid.UUID
		accountNumber string
		deposit       float64
		newDeposit    float64
		expected      float64
		err           error
	}{
		{
			name:       "empty user",
			userId:     uuid.MustParse("0162f9ca-6c7f-4bf9-9d79-ec1af9dd5dde"),
			deposit:    0,
			newDeposit: 100,
			expected:   100,
			err:        fmt.Errorf(apierr.ErrNoData.Message),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var accountNumber string
			for i := 0; i < 15; i++ {
				num := rand.Intn(10)
				accountNumber += strconv.Itoa(num)
			}

			testUser := &model.User{
				ID:            tc.userId,
				FirstName:     "user",
				LastName:      "test",
				DisplayName:   "test user",
				AccountNumber: accountNumber,
				Deposit:       tc.deposit,
			}
			// 初始化db
			err := initializeTestUserData(testDB, testUser)
			if err != nil {
				t.Errorf("initializeTestUserData failed: %v", err)
				return
			}

			defer func() {
				if err := teardownTestDatabase(testDB, tc.userId); err != nil {
					t.Errorf("teardownTestDatabase failed: %v", err)
				}
			}()

			// 建立 Service
			ctx := context.Background()
			svc := service.New(ctx, testDB)
			// call SaveMoney (原則上就是判斷有沒有非預期的error而已)
			err = svc.SaveMoney(tc.userId, tc.newDeposit)
			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("SaveMoney failed: %v", err)
				return
			}
		})
	}
}

func Test_integration_WithdrawMoney(t *testing.T) {
	// 建立測試資料庫
	testDB, err := setupTestDatabase()
	if err != nil {
		t.Errorf("setupTestDatabase error: %s", err.Error())
		return
	}
	// 這裡沒有做 newDeposit是負數的case 這個會在router層被擋下來 如果想測試的話 可能要寫e2e的測試 或是把那段檢查搬到service層
	tests := []struct {
		name          string
		userId        uuid.UUID
		accountNumber string
		deposit       float64
		newDeposit    float64
		expected      float64
		err           error
	}{
		// {
		// 這是一個蠻特別的case 是領0元或是負數的case 不過這個會在router層被擋下來 如果想測試的話 可能要寫e2e的測試 或是把那段檢查搬到service層
		// 後續有其他newDeposit為0的case就不列舉了
		// {
		// 	name:                      "empty user",
		// 	userId:                    uuid.MustParse("e03845a4-8749-468e-ae74-dac44d85959f"),
		// 	deposit:                   100,
		// 	newDeposit:                0,
		// 	expected:                  100,
		// 	err:                       fmt.Errorf(apierr.ErrNoData.Message),
		// },
		{
			name:       "empty user",
			userId:     uuid.MustParse("0e2e5dfd-0e26-4732-aedf-83f81ed53590"),
			deposit:    0,
			expected:   0,
			newDeposit: 100,
			err:        fmt.Errorf(apierr.ErrNoData.Message)},
		{
			name:       "origin deposit is zero",
			userId:     uuid.MustParse("35a6c515-ac88-4890-8de7-7015c380c9cd"),
			deposit:    0,
			expected:   0,
			newDeposit: 100,
			err:        fmt.Errorf(apierr.ErrNoData.Message),
		}, {
			name:       "origin deposit is positive number",
			userId:     uuid.MustParse("690211a6-ed99-41b2-8b70-d74b2025bfa7"),
			deposit:    100,
			expected:   0,
			newDeposit: 100,
			err:        nil,
		}, {
			name:       "origin deposit is negative",
			userId:     uuid.MustParse("f51d7a33-9f6b-4f8d-90d5-08ce13b633dc"),
			deposit:    -100,
			newDeposit: 100,
			expected:   -100,
			err:        fmt.Errorf(apierr.ErrNoData.Message),
		},
	}
	// 測試邏輯 1.建立User 2.比對User的存款是否正確 3.進行領錢 4.比對領完錢的存款是否正確
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var accountNumber string
			for i := 0; i < 15; i++ {
				num := rand.Intn(10)
				accountNumber += strconv.Itoa(num)
			}

			testUser := &model.User{
				ID:            tc.userId,
				FirstName:     "user",
				LastName:      "test",
				DisplayName:   "test user",
				AccountNumber: accountNumber,
				Deposit:       tc.deposit,
			}

			// 初始化db(1.建立User)
			err := initializeTestUserData(testDB, testUser)
			if err != nil {
				t.Errorf("initializeTestUserData failed: %v", err)
				return
			}

			defer func() {
				if err := teardownTestDatabase(testDB, tc.userId); err != nil {
					t.Errorf("teardownTestDatabase failed: %v", err)
				}
			}()

			// 建立 Service
			ctx := context.Background()
			svc := service.New(ctx, testDB)

			// 2.比對User的存款是否正確
			user, err := testDB.GetUserByID(testDB.DB, tc.userId)
			if err != nil {
				t.Errorf("first GetUserByID failed: %v", err)
				return
			}
			if user.Deposit != tc.deposit {
				t.Errorf("createdUser's Deposit is incorrect, expected %f, got %f", tc.deposit, user.Deposit)
				return
			}

			// call WithdrawMoney (3.進行領錢 也是最主要的測試部分)
			err = svc.WithdrawMoney(tc.userId, tc.newDeposit)
			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("WithdrawMoney failed: %v", err)
				return
			}
			// 4.比對存完錢的存款是否正確
			user, err = testDB.GetUserByID(testDB.DB, tc.userId)
			if err != nil {
				t.Errorf("second GetUserByID failed: %v", err)
				return
			}
			if user.Deposit != tc.expected {
				t.Errorf("final Deposit is incorrect, expected %f, got %f", tc.expected, user.Deposit)
				return
			}
		})
	}
}

func Test_integration_WithdrawMoneyToEmptyUser(t *testing.T) {
	// 建立測試資料庫
	testDB, err := setupTestDatabase()
	if err != nil {
		t.Errorf("setupTestDatabase error: %s", err.Error())
		return
	}
	// 這裡沒有做 newDeposit是負數的case 這個會在router層被擋下來 如果想測試的話 可能要寫e2e的測試 或是把那段檢查搬到service層
	tests := []struct {
		name                      string
		userId                    uuid.UUID
		accountNumber             string
		deposit                   float64
		newDeposit                float64
		unnecessaryCreateUserToDB bool
		expected                  float64
		err                       error
	}{
		{
			name:                      "empty user",
			userId:                    uuid.MustParse("0e2e5dfd-0e26-4732-aedf-83f81ed53590"),
			deposit:                   0,
			expected:                  0,
			newDeposit:                100,
			unnecessaryCreateUserToDB: false,
			err:                       fmt.Errorf(apierr.ErrNoData.Message),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var accountNumber string
			for i := 0; i < 15; i++ {
				num := rand.Intn(10)
				accountNumber += strconv.Itoa(num)
			}

			testUser := &model.User{
				ID:            tc.userId,
				FirstName:     "user",
				LastName:      "test",
				DisplayName:   "test user",
				AccountNumber: accountNumber,
				Deposit:       tc.deposit,
			}

			// 初始化db
			err := initializeTestUserData(testDB, testUser)
			if err != nil {
				t.Errorf("initializeTestUserData failed: %v", err)
				return
			}

			defer func() {
				if err := teardownTestDatabase(testDB, tc.userId); err != nil {
					t.Errorf("teardownTestDatabase failed: %v", err)
				}
			}()
			// 建立 Service
			ctx := context.Background()
			svc := service.New(ctx, testDB)
			// call SaveMoney (原則上就是判斷有沒有非預期的error而已)
			err = svc.WithdrawMoney(tc.userId, tc.newDeposit)
			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("WithdrawMoney failed: %v", err)
				return
			}
		})
	}
}

// 其實這裡還可以針對 logger去做驗證 可以把logger繼續往下記錄到service在這邊驗證 (詳細的作法就不特別做了 只是想表達這也是可以測試的部分)
func Test_integration_TransferMoney(t *testing.T) {
	// 建立測試資料庫
	testDB, err := setupTestDatabase()
	if err != nil {
		t.Errorf("setupTestDatabase error: %s", err.Error())
		return
	}
	tests := []struct {
		name                     string
		transferOutUserId        uuid.UUID
		transferInUserId         uuid.UUID
		transferOutAccountNumber string
		transferInAccountNumber  string
		transferOutDeposit       float64
		transferInDeposit        float64
		newDeposit               float64
		transferOutExpected      float64
		transferInExpected       float64
		err                      error
	}{
		// {
		// 這是一個蠻特別的case 是轉0元或是負數的case 不過這個會在router層被擋下來 如果想測試的話 可能要寫e2e的測試 或是把那段檢查搬到service層
		// 後續有其他newDeposit為0的case就不列舉了
		// name:               "transferOutUserId and transferInUserId deposit is zero and newDeposit is zero",
		// transferOutUserId:  uuid.MustParse("f7ca3aa3-82de-4ae7-9681-df453f3ccb6a"),
		// transferInUserId:   uuid.MustParse("e2b2163d-9547-466f-ae1d-6c9fb19a56fc"),
		// transferOutDeposit:  100,
		// transferInDeposit:   0,
		// newDeposit:          0,
		// transferOutExpected: 100,
		// transferInExpected:  0,
		// err:                nil,
		//},
		{
			name:                "empty TransferInUser",
			transferOutUserId:   uuid.MustParse("815b9ad6-9c1b-4766-9bcb-35095601e317"),
			transferInUserId:    uuid.MustParse("ddfd8a15-3302-4bf2-8d1e-8b8e1f75f48c"),
			transferOutDeposit:  100,
			transferInDeposit:   0,
			newDeposit:          50,
			transferOutExpected: 50,
			transferInExpected:  50,
			err:                 fmt.Errorf(apierr.ErrNoData.Message),
		},
		{
			name:                "empty TransferOutUser",
			transferOutUserId:   uuid.MustParse("e6decf0b-0dfa-4dee-9204-d381de6a5f28"),
			transferInUserId:    uuid.MustParse("238f29d3-40c7-4b12-a7fa-3ed0120c1e7e"),
			transferOutDeposit:  100,
			transferInDeposit:   0,
			newDeposit:          50,
			transferOutExpected: 50,
			transferInExpected:  50,
			err:                 fmt.Errorf(apierr.ErrNoData.Message),
		},
		{
			name:                "empty TransferOutUser and TransferInUser",
			transferOutUserId:   uuid.MustParse("bc67d214-f4fc-4d47-89af-9219f93d8a46"),
			transferInUserId:    uuid.MustParse("4b37606b-b695-435a-bbc5-270a102cf5bc"),
			transferOutDeposit:  100,
			transferInDeposit:   0,
			newDeposit:          50,
			transferOutExpected: 50,
			transferInExpected:  50,
			err:                 fmt.Errorf(apierr.ErrNoData.Message),
		},
		{
			name:                "origin deposit is zero",
			transferOutUserId:   uuid.MustParse("c558ffcd-b21c-4cf3-885b-f32b6a25cdd8"),
			transferInUserId:    uuid.MustParse("d4306904-0724-49ef-a176-9083baef026f"),
			transferOutDeposit:  0,
			transferInDeposit:   0,
			newDeposit:          10,
			transferOutExpected: 0,
			transferInExpected:  0,
			err:                 fmt.Errorf(apierr.ErrNoData.Message),
		},
		{
			name:                "origin deposit is negative number and newDeposit > transferOutDeposit",
			transferOutUserId:   uuid.MustParse("271f13ee-fa9e-4940-9ef5-20d90e84d7fa"),
			transferInUserId:    uuid.MustParse("8b58d996-5763-49d2-840a-168d39fde846"),
			transferOutDeposit:  -10,
			transferInDeposit:   0,
			newDeposit:          10,
			transferOutExpected: -10,
			transferInExpected:  0,
			err:                 fmt.Errorf(apierr.ErrNoData.Message),
		},
		{
			name:                "origin deposit is positive number and newDeposit < transferOutDeposit",
			transferOutUserId:   uuid.MustParse("995728f2-cf11-479d-a86b-1ca612c1c144"),
			transferInUserId:    uuid.MustParse("922defc3-16d1-4e0a-90c0-b89c98a8ea10"),
			transferOutDeposit:  100,
			transferInDeposit:   0,
			newDeposit:          50,
			transferOutExpected: 50,
			transferInExpected:  50,
			err:                 nil,
		},
		{
			name:                "origin deposit is positive number and newDeposit = transferOutDeposit",
			transferOutUserId:   uuid.MustParse("36f26736-fd18-42dd-9109-4a62b3547ee6"),
			transferInUserId:    uuid.MustParse("bb0ae4f3-f6e4-4e38-85e5-726c36b73d9a"),
			transferOutDeposit:  100,
			transferInDeposit:   0,
			newDeposit:          100,
			transferOutExpected: 0,
			transferInExpected:  100,
			err:                 nil,
		},
		{
			name:                "origin deposit is positive number and newDeposit > transferOutDeposit",
			transferOutUserId:   uuid.MustParse("4c170466-822e-497c-baee-17c7657a09ee"),
			transferInUserId:    uuid.MustParse("0b16747a-5e66-49b8-b1bc-68e2baf80d7f"),
			transferOutDeposit:  100,
			transferInDeposit:   0,
			newDeposit:          110,
			transferOutExpected: 100,
			transferInExpected:  0,
			err:                 fmt.Errorf(apierr.ErrNoData.Message),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var transferOutUserAccountNumber string
			var transferInUserAccountNumber string

			for i := 0; i < 15; i++ {
				transferOutNum := rand.Intn(10)
				transferInNum := rand.Intn(10)
				transferOutUserAccountNumber += strconv.Itoa(transferOutNum)
				transferInUserAccountNumber += strconv.Itoa(transferInNum)
			}

			var testUsers []*model.User
			TransferOutUser := &model.User{
				ID:            tc.transferOutUserId,
				FirstName:     "user",
				LastName:      "test",
				DisplayName:   "test user",
				AccountNumber: transferOutUserAccountNumber,
				Deposit:       tc.transferOutDeposit,
			}
			testUsers = append(testUsers, TransferOutUser)

			TransferInUser := &model.User{
				ID:            tc.transferInUserId,
				FirstName:     "user",
				LastName:      "test",
				DisplayName:   "test user",
				AccountNumber: transferInUserAccountNumber,
				Deposit:       tc.transferInDeposit,
			}
			testUsers = append(testUsers, TransferInUser)

			// 避免transferOut跟transferIn都沒有的情況
			if len(testUsers) > 0 {
				// 初始化db(1.建立User)
				err := initializeTestUsersData(testDB, testUsers)
				if err != nil {
					t.Errorf("initializeTestUsersData failed: %v", err)
					return
				}
			}
			defer func() {
				if err := teardownTestDatabase(testDB, tc.transferOutUserId); err != nil {
					t.Errorf("teardownTestDatabase failed: %v", err)
				}
			}()

			defer func() {
				if err := teardownTestDatabase(testDB, tc.transferInUserId); err != nil {
					t.Errorf("teardownTestDatabase failed: %v", err)
				}
			}()

			// 建立 Service
			ctx := context.Background()
			svc := service.New(ctx, testDB)

			// 2.比對User的存款是否正確
			transferInUser, err := testDB.GetUserByID(testDB.DB, tc.transferInUserId)
			if err != nil {
				t.Errorf("first GetUserByID transferInUser failed: %v", err)
				return
			}

			transferOutUser, err := testDB.GetUserByID(testDB.DB, tc.transferOutUserId)
			if err != nil {
				t.Errorf("first GetUserByID transferOutUser failed: %v", err)
				return
			}

			if transferInUser.Deposit != tc.transferInDeposit {
				t.Errorf("createdUser's transferInUser.Deposit is incorrect, expected %f, got %f", tc.transferInDeposit, transferInUser.Deposit)
				return
			}

			if transferOutUser.Deposit != tc.transferOutDeposit {
				t.Errorf("createdUser's transferOutUser.Deposit is incorrect, expected %f, got %f", tc.transferOutDeposit, transferOutUser.Deposit)
				return
			}

			// call TransferMoney (3.進行轉帳 也是最主要的測試部分)
			transferMoneyStruct := &service.TransferMoneyStruct{
				NewDeposit:        tc.newDeposit,
				TransferOutUserId: tc.transferOutUserId,
				TransferInUserId:  tc.transferInUserId,
			}

			err = svc.TransferMoney(*transferMoneyStruct)
			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("WithdrawMoney failed: %v", err)
				return
			}

			// 4.比對存完錢的存款是否正確
			transferInUser, err = testDB.GetUserByID(testDB.DB, tc.transferInUserId)
			if err != nil {
				t.Errorf("second GetUserByID transferInUser failed: %v", err)
				return
			}

			transferOutUser, err = testDB.GetUserByID(testDB.DB, tc.transferOutUserId)
			if err != nil {
				t.Errorf("second GetUserByID transferOutUser failed: %v", err)
				return
			}

			if transferInUser.Deposit != tc.transferInExpected {
				t.Errorf("final transferInUser.Deposit is incorrect, expected %f, got %f", tc.transferInExpected, transferInUser.Deposit)
				return
			}

			if transferOutUser.Deposit != tc.transferOutExpected {
				t.Errorf("final transferOutUser.Deposit is incorrect, expected %f, got %f", tc.transferOutExpected, transferOutUser.Deposit)
				return
			}
		})
	}
}

func Test_integration_TransferMoneyToEmptyUser(t *testing.T) {
	// 建立測試資料庫
	testDB, err := setupTestDatabase()
	if err != nil {
		t.Errorf("setupTestDatabase error: %s", err.Error())
		return
	}
	tests := []struct {
		name                                 string
		transferOutUserId                    uuid.UUID
		transferInUserId                     uuid.UUID
		transferOutAccountNumber             string
		transferInAccountNumber              string
		transferOutDeposit                   float64
		transferInDeposit                    float64
		newDeposit                           float64
		unnecessaryCreateTransferOutUserToDB bool
		unnecessaryCreateTransferInUserToDB  bool
		transferOutExpected                  float64
		transferInExpected                   float64
		err                                  error
	}{
		{
			name:                                 "empty TransferInUser",
			transferOutUserId:                    uuid.MustParse("815b9ad6-9c1b-4766-9bcb-35095601e317"),
			transferInUserId:                     uuid.MustParse("ddfd8a15-3302-4bf2-8d1e-8b8e1f75f48c"),
			transferOutDeposit:                   100,
			transferInDeposit:                    0,
			newDeposit:                           50,
			unnecessaryCreateTransferOutUserToDB: false,
			unnecessaryCreateTransferInUserToDB:  true,
			transferOutExpected:                  50,
			transferInExpected:                   50,
			err:                                  fmt.Errorf(apierr.ErrNoData.Message),
		}, {
			name:                                 "empty TransferOutUser",
			transferOutUserId:                    uuid.MustParse("e6decf0b-0dfa-4dee-9204-d381de6a5f28"),
			transferInUserId:                     uuid.MustParse("238f29d3-40c7-4b12-a7fa-3ed0120c1e7e"),
			transferOutDeposit:                   100,
			transferInDeposit:                    0,
			newDeposit:                           50,
			unnecessaryCreateTransferOutUserToDB: true,
			unnecessaryCreateTransferInUserToDB:  false,
			transferOutExpected:                  50,
			transferInExpected:                   50,
			err:                                  fmt.Errorf(apierr.ErrNoData.Message),
		}, {
			name:                                 "empty TransferOutUser and TransferInUser",
			transferOutUserId:                    uuid.MustParse("bc67d214-f4fc-4d47-89af-9219f93d8a46"),
			transferInUserId:                     uuid.MustParse("4b37606b-b695-435a-bbc5-270a102cf5bc"),
			transferOutDeposit:                   100,
			transferInDeposit:                    0,
			newDeposit:                           50,
			unnecessaryCreateTransferOutUserToDB: true,
			unnecessaryCreateTransferInUserToDB:  true,
			transferOutExpected:                  50,
			transferInExpected:                   50,
			err:                                  fmt.Errorf(apierr.ErrNoData.Message),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var transferOutUserAccountNumber string
			var transferInUserAccountNumber string

			for i := 0; i < 15; i++ {
				transferOutNum := rand.Intn(10)
				transferInNum := rand.Intn(10)
				transferOutUserAccountNumber += strconv.Itoa(transferOutNum)
				transferInUserAccountNumber += strconv.Itoa(transferInNum)
			}

			var testUsers []*model.User
			if !tc.unnecessaryCreateTransferOutUserToDB {
				TransferOutUser := &model.User{
					ID:            tc.transferOutUserId,
					FirstName:     "user",
					LastName:      "test",
					DisplayName:   "test user",
					AccountNumber: transferOutUserAccountNumber,
					Deposit:       tc.transferOutDeposit}
				testUsers = append(testUsers, TransferOutUser)
			}

			if !tc.unnecessaryCreateTransferInUserToDB {
				TransferInUser := &model.User{
					ID:            tc.transferInUserId,
					FirstName:     "user",
					LastName:      "test",
					DisplayName:   "test user",
					AccountNumber: transferInUserAccountNumber,
					Deposit:       tc.transferInDeposit}
				testUsers = append(testUsers, TransferInUser)
			}

			// 避免transferOut跟transferIn都沒有的情況
			if len(testUsers) > 0 {
				// 初始化db
				err := initializeTestUsersData(testDB, testUsers)
				if err != nil {
					t.Errorf("initializeTestUsersData failed: %v", err)
					return
				}
			}
			if !tc.unnecessaryCreateTransferOutUserToDB {
				defer func() {
					if err := teardownTestDatabase(testDB, tc.transferOutUserId); err != nil {
						t.Errorf("teardownTestDatabase failed: %v", err)
					}
				}()
			}
			if !tc.unnecessaryCreateTransferInUserToDB {
				defer func() {
					if err := teardownTestDatabase(testDB, tc.transferInUserId); err != nil {
						t.Errorf("teardownTestDatabase failed: %v", err)
					}
				}()
			}
			// 建立 Service
			ctx := context.Background()
			svc := service.New(ctx, testDB)
			// call SaveMoney (原則上就是判斷有沒有非預期的error而已)
			transferMoneyStruct := &service.TransferMoneyStruct{
				NewDeposit:        tc.newDeposit,
				TransferOutUserId: tc.transferOutUserId,
				TransferInUserId:  tc.transferInUserId,
			}

			err = svc.TransferMoney(*transferMoneyStruct)
			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("WithdrawMoney failed: %v", err)
				return
			}
		})
	}
}

func setupTestDatabase() (*database.GormDatabase, error) {
	// 設置相對於config文件的路徑
	viper.SetConfigFile("../../app.yaml")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	cfg := new(config.Config)
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	postgresClient := database.GetPostgresCli(cfg)
	return postgresClient, nil
}

func teardownTestDatabase(db *database.GormDatabase, testUserId uuid.UUID) error {
	return db.DB.Exec(`DELETE FROM users WHERE id = ?`, testUserId).Error
}

func initializeTestUserData(db *database.GormDatabase, testUser *model.User) error {
	return db.DB.Create(&testUser).Error
}

func initializeTestUsersData(db *database.GormDatabase, testUsers []*model.User) error {
	return db.DB.Create(testUsers).Error
}
