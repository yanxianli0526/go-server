package apierr

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ApiErr struct { // can be modified to wrap error and log
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	ErrInvalidParam = ApiErr{Code: 1, Message: "invalid params"}
	ErrInternal     = ApiErr{Code: 2, Message: "internal error"}
	// 定義了幾個 會出現的錯誤 Message 下面都用some error 混淆,有需要的話 可以把Code放到log觀測 (實際上要看前端業務邏輯決定有沒有需要做不同的Message)
	ErrNoData                       = ApiErr{Code: 3, Message: "some error"}
	ErrTotalDepositIsNegativeNumber = ApiErr{Code: 4, Message: "some error"}
)

func (e ApiErr) Error() string {
	return e.Message
}

// SuccessOrAbort 會判斷是否讓程式繼續向後執行，若 err 不存在則繼續執行
func SuccessOrAbort(
	ctx *gin.Context,
	statusCode int,
	err error,
) bool {
	// 表示在 ctx 中已經就有錯誤存在
	if len(ctx.Errors) > 0 {
		return false
	}

	// 當錯誤是不合法的 JSON syntax 時
	var jsonSyntaxErr *json.SyntaxError
	if errors.As(err, &jsonSyntaxErr) {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid JSON syntax"))
		return false
	}

	// 當錯誤是 request 中不正確的 params，把 ErrorType 設成 ErrorTypeBind
	var validationErrs validator.ValidationErrors
	if ok := errors.As(err, &validationErrs); ok {
		_ = ctx.AbortWithError(http.StatusBadRequest, validationErrs).SetType(gin.ErrorTypeBind)
		return false
	}

	// 如果是其他類型的錯誤
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return false
	}

	return true
}
