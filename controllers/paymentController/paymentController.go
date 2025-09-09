package paymentcontroller

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/models"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func BuyClass(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	if role != "student" {
		c.JSON(403, models.Response{
			Message: "only student can buy classes"})
		return
	}

	var student models.User
	if err := config.DB.Preload("Classes").First(&student, userID).Error; err != nil {
		c.JSON(404, models.Response{
			Message: "user not found"})
		return
	}

	classIDParam := c.Param("id")
	classID, err := strconv.Atoi(classIDParam)
	if err != nil {
		c.JSON(400, models.Response{
			Message: "invalid class id"})
		return
	}

	var class models.Class
	if err := config.DB.First(&class, classID).Error; err != nil {
		c.JSON(404, models.Response{
			Message: "class not found"})
		return
	}

	for _, cl := range student.Classes {
		if cl.Id == class.Id {
			c.JSON(400, models.Response{
				Message: "already enrolled"})
			return
		}
	}

	orderID := fmt.Sprintf("ORDER-%d-%d", student.Id, time.Now().Unix())

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(class.Price),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: student.Name,
			Email: student.Email,
		},
	}

	snapResp, errSnap := config.SnapClient.CreateTransaction(req)
	if errSnap != nil {
		c.JSON(500, models.Response{
			Message: fmt.Sprintf("midtrans error: %+v", err),
		})
		return
	}

	if snapResp == nil {
		c.JSON(500, models.Response{
			Message: "empty response from midtrans"})
		return
	}

	transaction := models.Transaction{
		OrderID: orderID,
		UserID:  student.Id,
		ClassID: class.Id,
		Amount:  int64(class.Price),
		Status:  "pending",
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		c.JSON(500, models.Response{
			Message: "failed to save transaction"})
		return
	}

	c.JSON(200, models.Response{
		Message: "Successfully created payment",
		Data:    snapResp,
	})
}

func PaymentNotification(c *gin.Context) {
	var notif map[string]interface{}

	if err := c.BindJSON(&notif); err != nil {
		c.JSON(400, models.Response{
			Message: "invalid payload"})
		return
	}

	orderID, _ := notif["order_id"].(string)
	statusCode, _ := notif["status_code"].(string)
	grossAmount, _ := notif["gross_amount"].(string)
	transactionStatus, _ := notif["transaction_status"].(string)
	fraudStatus, _ := notif["fraud_status"].(string)
	signatureKey, _ := notif["signature_key"].(string)

	if !isValidSignature(orderID, statusCode, grossAmount, signatureKey) {
		c.JSON(401, models.Response{
			Message: "invalid signature"})
		return
	}

	var trx models.Transaction
	if err := config.DB.Where("order_id = ?", orderID).First(&trx).Error; err != nil {
		c.JSON(404, models.Response{
			Message: "transaction not found"})
		return
	}

	switch transactionStatus {
	case "capture":
		if fraudStatus == "challenge" {
			trx.Status = "challenge"
		} else if fraudStatus == "accept" {
			trx.Status = "paid"
			assignUserToClass(trx)
		}
	case "settlement":
		trx.Status = "paid"
		assignUserToClass(trx)
	case "deny", "cancel", "expire":
		trx.Status = "failed"
	case "pending":
		trx.Status = "pending"
	default:
		trx.Status = "unknown"
	}

	config.DB.Save(&trx)

	c.JSON(200, models.Response{
		Message: "notification processed"})
}

func assignUserToClass(trx models.Transaction) {

	var user models.User
	if err := config.DB.Preload("Classes").First(&user, trx.UserID).Error; err != nil {
		return
	}

	var class models.Class
	if err := config.DB.First(&class, trx.ClassID).Error; err != nil {
		return
	}

	config.DB.Model(&user).Association("Classes").Append(&class)
}

func isValidSignature(orderID, statusCode, grossAmount, signatureKey string) bool {
	serverKey := config.MidtransSKey
	raw := orderID + statusCode + grossAmount + serverKey

	hash := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(hash[:])

	return expected == signatureKey
}
