package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rate-api/config"
	"rate-api/internal/logger"
	"rate-api/internal/mailerror"
	"rate-api/internal/model"

	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/gin-gonic/gin"
)

type EmailServiceInterface interface {
	AddEmail(email model.Email) error
	DeleteEmail(email model.Email) error
	GetAllEmails() []string
}

type EmailHandler struct {
	emailServ EmailServiceInterface
	log       logger.LoggerInterface
}

func NewEmailHandler(emailServ EmailServiceInterface, log logger.LoggerInterface) EmailHandler {
	return EmailHandler{emailServ, log}
}

func (h *EmailHandler) EmailPlusCustomerSubscribe(context *gin.Context) {
	h.log.Debug("EmailPlusCustomerSubscribe started")
	var newEmail model.Email

	if err := context.Bind(&newEmail); err != nil {
		h.log.Error(fmt.Sprintf("Error occurred during adding subscription: %s", err.Error()))
		http.Error(context.Writer, err.Error(), http.StatusConflict)
		return
	}
	h.log.Debug(newEmail)

	dtmAddr := config.Cfg.DTMAddr
	globalTransactionID := dtmcli.MustGenGid(dtmAddr)

	emailsServURL := config.Cfg.EmailsServerURL
	customersServURL := config.Cfg.CustomersServerURL

	registerSubscriberURL := emailsServURL + "/register-subscriber"
	registerSubscriberCompensateURL := emailsServURL + "/register-subscriber-compensate"

	registerCustomerURL := customersServURL + "/register-customer"
	registerCustomerCompensateURL := customersServURL + "/register-customer-compensate"

	reqBody, err := structToMap(newEmail)
	if err != nil {
		context.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
		return
	}

	saga := dtmcli.
		NewSaga(dtmAddr, globalTransactionID).
		Add(registerSubscriberURL, registerSubscriberCompensateURL, reqBody).
		Add(registerCustomerURL, registerCustomerCompensateURL, reqBody)

	saga.WaitResult = true
	saga.WithRetryLimit(1)

	if err = saga.Submit(); err != nil {
		code, data := dtmcli.Result2HttpJSON(err)
		dataMap, ok := data.(map[string]string)
		if !ok {
			context.JSON(code, data)
			return
		}
		context.Data(code, "application/json", []byte(dataMap["error"]))
		return
	}
	context.Status(http.StatusOK)
}

func (h *EmailHandler) Subscribe(context *gin.Context) interface{} {
	h.log.Debug("SubscribeEmail started")
	var newEmail model.Email
	err := context.BindJSON(&newEmail)
	if err != nil {
		return err
	}

	newEmail.TransactionID = context.Query("gid")
	if err := h.emailServ.AddEmail(newEmail); err != nil {
		if errors.Is(err, mailerror.ErrEmailSubscribed) {
			h.log.Debug("Email already subscribed")
			http.Error(context.Writer, err.Error(), http.StatusConflict)
			return dtmcli.StatusSucceed
		}
		if errors.Is(err, mailerror.ErrEmailNotValid) {
			h.log.Error("Not valid email provided")
			http.Error(context.Writer, err.Error(), http.StatusBadRequest)
			return dtmcli.ErrFailure
		}

		h.log.Error(fmt.Sprintf("Error occurred during adding subscription: %s", err.Error()))
		http.Error(context.Writer, err.Error(), http.StatusInternalServerError)
		return dtmcli.ErrFailure
	}

	return nil
}

func (h *EmailHandler) SubscribeEmailCompensate(context *gin.Context) interface{} {
	h.log.Debug("SubscribeEmailCompensate started")
	var newEmail model.Email
	err := context.BindJSON(&newEmail)
	if err != nil {
		return err
	}
	newEmail.TransactionID = context.Query("gid")

	err = h.emailServ.DeleteEmail(newEmail)
	if err != nil {
		return err
	}

	return nil
}

func structToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &newMap)
	return
}

func extractCode(err error) int {
	statusCode := http.StatusOK
	if err != nil {
		statusCode = http.StatusInternalServerError
	}
	return statusCode
}
