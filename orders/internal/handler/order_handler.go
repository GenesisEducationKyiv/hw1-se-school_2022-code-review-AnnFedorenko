package handler

import (
	"encoding/json"
	"net/http"
	"orders/internal/model"
	"os"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var dtmCoordinatorAddress = os.Getenv("DTM_COORDINATOR")
var ordersServerURL = os.Getenv("ORDERS_SERVICE_URL")
var customersServerURL = os.Getenv("CUSTOMERS_SERVICE_URL")

type OrderHandler struct {
	db *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	createOrderRequest := struct {
		IdCustomer uint   `json:"idCustomer"`
		Currency   string `json:"currency"`
		Amount     uint   `json:"amount"`
	}{}

	err := c.BindJSON(&createOrderRequest)
	if err != nil {
		c.JSON(extractCode(err), err)
		return
	}

	globalTransactionId := dtmcli.MustGenGid(dtmCoordinatorAddress)
	req, _ := StructToMap(createOrderRequest)
	err = dtmcli.
		NewSaga(dtmCoordinatorAddress, globalTransactionId).
		Add(ordersServerURL+"/register-order", ordersServerURL+"/register-order-compensate", req).
		Add(customersServerURL+"/withdraw-money", customersServerURL+"/withdraw-money-compensate", req).
		Submit()

	createOrderResponse := struct {
		Gid string `json:"gid"`
	}{Gid: globalTransactionId}

	c.JSON(extractCode(err), createOrderResponse)

}

func (h *OrderHandler) RegisterOrder(c *gin.Context) interface{} {
	registerOrderRequest := struct {
		IdCustomer uint   `json:"idCustomer"`
		Currency   string `json:"currency"`
		Amount     uint   `json:"amount"`
	}{}
	transactionId := c.Query("gid")

	err := c.BindJSON(&registerOrderRequest)
	if err != nil {
		return err
	}

	return h.db.
		Create(&model.Order{
			IDTransaction: transactionId,
			IDCustomer:    registerOrderRequest.IdCustomer,
			Currency:      registerOrderRequest.Currency,
			Amount:        registerOrderRequest.Amount,
			Status:        "created",
		}).
		Error
}

func (h *OrderHandler) RegisterOrderCompensate(c *gin.Context) interface{} {
	transactionId := c.Query("gid")

	return h.db.
		Model(&model.Order{}).
		Where("id_transaction = ?", transactionId).
		Update("status", "canceled").
		Limit(1).
		Error
}

func extractCode(err error) int {
	statusCode := http.StatusOK
	if err != nil {
		statusCode = http.StatusInternalServerError
	}
	return statusCode
}

func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &newMap)
	return
}
