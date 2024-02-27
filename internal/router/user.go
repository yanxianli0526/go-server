package router

import (
	"time"

	"go.uber.org/zap"

	"meepShopTest/internal/apierr"
	"meepShopTest/internal/database"
	"meepShopTest/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterUser(
	db *database.GormDatabase,
	logger *zap.Logger,
	routerGroup *gin.RouterGroup,
) {
	userHandler := NewUserHandler(db, logger)
	VitalLinkUserRouter := routerGroup.Group("/user")

	{
		// 如果要新增user的話 可以用這個 不過因為是作業 就只寫一個殼 meepShopTest.sql已經有兩個user可以使用了
		// 有需要再新增user可以直接下sql指令
		VitalLinkUserRouter.POST("/", userHandler.CreateUser)
		VitalLinkUserRouter.PATCH("/:userId/saveMoney", userHandler.SaveMoney)
		VitalLinkUserRouter.PATCH("/:userId/withdrawMoney", userHandler.WithdrawMoney)
		VitalLinkUserRouter.PATCH("/transferMoney", userHandler.TransferMoney)
	}
}

type UserAPI struct {
	DB     *database.GormDatabase
	logger *zap.Logger
}

func NewUserHandler(db *database.GormDatabase, logger *zap.Logger) *UserAPI {
	return &UserAPI{
		DB:     db,
		logger: logger,
	}
}

func (u *UserAPI) CreateUser(ctx *gin.Context) {
	// do createUser
}

func (u *UserAPI) SaveMoney(ctx *gin.Context) {
	userIdStr := ctx.Param("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		apierr.SuccessOrAbort(ctx, http.StatusBadRequest, err)
		return
	}
	param := service.SaveMoneyRequestStruct{}
	err = ctx.Bind(&param)
	// 這邊要判斷傳入的值是不是負數(存0元是一個挺不合理的行為)
	// 也可以拆開兩個不同的 error給前端 但因為跟金錢有關的功能,以防被一些奇怪的攻擊,所以不這麼做
	if err != nil || param.NewDeposit <= 0 {
		apierr.SuccessOrAbort(ctx, http.StatusBadRequest, err)
		return
	}
	svc := service.New(ctx, u.DB)
	err = svc.SaveMoney(userId, param.NewDeposit)
	if err != nil || param.NewDeposit < 0 {
		apierr.SuccessOrAbort(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (u *UserAPI) WithdrawMoney(ctx *gin.Context) {
	userIdStr := ctx.Param("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		apierr.SuccessOrAbort(ctx, http.StatusBadRequest, err)
		return
	}
	param := service.WithdrawMoneyRequestStruct{}
	err = ctx.Bind(&param)
	// 這邊要判斷傳入的值是不是負數(領0元是一個挺不合理的行為)
	// 也可以拆開兩個不同的 error給前端 但因為跟金錢有關的功能,以防被一些奇怪的攻擊,所以不這麼做
	if err != nil || param.NewDeposit <= 0 {
		apierr.SuccessOrAbort(ctx, http.StatusBadRequest, err)
		return
	}
	svc := service.New(ctx, u.DB)
	err = svc.WithdrawMoney(userId, param.NewDeposit)
	if err != nil || param.NewDeposit < 0 {
		apierr.SuccessOrAbort(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (u *UserAPI) TransferMoney(ctx *gin.Context) {
	startTimeUnix := time.Now().Unix()
	param := service.TransferMoneyStruct{}
	err := ctx.Bind(&param)
	// 用來辦別開始跟結束
	transferUUIDStr := uuid.New().String()
	u.logger.Info("transferMoney start",
		zap.String("transferUUIDStr", transferUUIDStr),
		zap.Int64("time", startTimeUnix),
		zap.Float64("newDeposit", param.NewDeposit),
		zap.String("backoff", param.TransferInUserIdStr),
	)
	// 這邊要判斷傳入的值是不是負數
	// 也可以拆開兩個不同的 error 給前端 但因為跟金錢有關的功能,以防被一些奇怪的攻擊,所以不這麼做
	if err != nil || param.NewDeposit <= 0 {
		apierr.SuccessOrAbort(ctx, http.StatusBadRequest, err)
		return
	}

	param.TransferInUserId, err = uuid.Parse(param.TransferInUserIdStr)
	if success := apierr.SuccessOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}

	param.TransferOutUserId, err = uuid.Parse(param.TransferOutUserIdStr)
	if success := apierr.SuccessOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}

	svc := service.New(ctx, u.DB)
	err = svc.TransferMoney(param)
	if err != nil || param.NewDeposit < 0 {
		apierr.SuccessOrAbort(ctx, http.StatusBadRequest, err)
		return
	}

	// 結果加上一個執行總時間 未來可以針對duration > x 秒的api進行優化
	u.logger.Info("transferMoney end",
		zap.String("transferUUIDStr", transferUUIDStr),
		zap.Int64("time", time.Now().Unix()),
		zap.Int64("duration", time.Now().Unix()-startTimeUnix),
		zap.Float64("newDeposit", param.NewDeposit),
		zap.String("backoff", param.TransferInUserIdStr),
	)
	ctx.JSON(http.StatusNoContent, nil)
}
