package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/code"
	"gobeacon/model"
	"gobeacon/service"
)

// CurrentGroupAccounts godoc
// @Summary Получение списка пользователей текущей подписки
// @Description Получение списка пользователей текущей подписки
// @Accept json
// @Produce json
// @Router /subscription/my/group [get]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Subscription
func CurrentGroupAccounts(c *gin.Context) {
	userId := getUserId(c)
	result, err := service.GetAllAccountWithMySubscription(userId)
	sendObjResponse(result, err, c)
}

// BuySubscription godoc
// @Summary Совершение покупки подписки
// @Description Покупка подписки для авторизованного пользователя и его связанных аккаунтов, в сумме не больше 5
// @Accept json
// @Produce json
// @Param request body model.BuySubscriptionRequest true "Запрос на добавление подписки"
// @Router /subscription/my/buy [post]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Subscription
func BuySubscription(c *gin.Context) {
	req := model.BuySubscriptionRequest{UserId:getUserId(c)}
	if e := c.Bind(&req); e != nil {
		sendResponse([]int{code.ParseRequest}, c)
	}
	err := service.BuySubscription(&req)
	sendResponse(err, c)
}

// CurrentSubscription godoc
// @Summary Текущая подписка
// @Description Получить действующую подписку для пользователя
// @Accept json
// @Produce json
// @Router /subscription/my/current [get]
// @Success 200 "ok"
// @Failure 500 "err"
// @Tags Subscription
func CurrentSubscription(c *gin.Context) {
	userId := getUserId(c)
	result, err := service.CurrentSubscription(userId)
	sendObjResponse(result, err, c)
}

// AllActiveSubscription godoc
// @Summary Подписки пользователя
// @Description Подписки пользователя активная или те, которые будут активны в будущем
// @Accept json
// @Produce json
// @Router /subscription/my/all [get]
// @Success 200 "ok"
// @Failure 500 "err"
// @Tags Subscription
func AllActiveSubscription(c *gin.Context) {
	userId := getUserId(c)
	result, err := service.AllActiveSubscription(userId)
	sendObjResponse(result, err, c)
}

// Subscriptions godoc
// @Summary Список подписок
// @Description Список подписок доступных для покупки
// @Accept json
// @Produce json
// @Router /subscription/available-buy [get]
// @Success 200 "ok"
// @Failure 500 "err"
// @Tags Subscription
func Subscriptions(c *gin.Context) {
	userId := getUserId(c)
	result, err := service.Subscriptions(userId)
	sendObjResponse(result, err, c)
}
