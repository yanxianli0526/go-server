package service_test

import (
	"meepShopTest/internal/service"
	"testing"
)

func Test_unit_CheckTotalDepositIsGreaterThanZero(t *testing.T) {
	testCases := []struct {
		userDeposit float64
		newDeposit  float64
		expected    float64
		wantErr     bool
	}{
		// WantErr is false
		{userDeposit: 0.0, newDeposit: 150.0, expected: 150.0, wantErr: false},
		{userDeposit: 0.0, newDeposit: -150.0, expected: 0, wantErr: false},
		{userDeposit: 100.0, newDeposit: 50.0, expected: 150.0, wantErr: false},
		{userDeposit: 100.0, newDeposit: -50.0, expected: 50.0, wantErr: false},
		{userDeposit: 100.0, newDeposit: -100.0, expected: 0, wantErr: false},
		{userDeposit: -100.0, newDeposit: 50.0, expected: 0, wantErr: false},
		{userDeposit: -100.0, newDeposit: -50.0, expected: 0, wantErr: false},
		// WantErr is true
		{userDeposit: 0.0, newDeposit: -50.0, expected: 0, wantErr: true},
		{userDeposit: 100.0, newDeposit: -50.0, expected: 50, wantErr: true},
		{userDeposit: -100.0, newDeposit: -50.0, expected: 0, wantErr: true},
	}

	for _, tc := range testCases {
		result, err := service.CheckTotalDepositIsGreaterThanZero(tc.userDeposit, tc.newDeposit)
		if err != nil {
			if tc.wantErr != false && result != tc.expected {
				t.Errorf("CheckTotalDepositIsGreaterThanZero(%f, %f) = %f; err %f", tc.userDeposit, tc.newDeposit, result, tc.expected)
			}
		}
		if result != tc.expected {
			t.Errorf("CheckTotalDepositIsGreaterThanZero(%f, %f) = %f; expected %f", tc.userDeposit, tc.newDeposit, result, tc.expected)
		}

	}
}
