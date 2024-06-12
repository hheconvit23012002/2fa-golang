package main

import (
	"errors"
	"fmt"
	"github.com/0x19/goesl"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type UserBindData struct {
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Password    string `json:"password" form:"password"`
}

type UserLogin struct {
	PhoneNumber string     `json:"phone_number"`
	Password    string     `json:"password" `
	Pass2FA     int        `json:"pass_2_fa"`
	CreatedAt   *time.Time `json:"created_at"`
	Expired     *time.Time `json:"expired"`
	Enter2FA    int        `json:"enter_2_fa"`
	ChannelName string     `json:"channel_name"`
}

type UserCheck2FA struct {
	PhoneNumber string `json:"phone_number" form:"phone_number"`
}

type Event struct {
	Digit        string `json:"digit" form:"digit"`
	NumberCalled string `json:"caller_number" form:"caller_number"`
}

type GoConnectFreeSwitch struct {
	con *goesl.Client
}

func (conInstance *GoConnectFreeSwitch) connect() (bool, error) {
	conTmp, errConnect := goesl.NewClient("192.168.124.128", 8021, "ClueCon", 10)
	if errConnect != nil {
		return false, errConnect
	}
	conInstance.con = conTmp
	return true, nil
}

func (conInstance *GoConnectFreeSwitch) closeConnect() bool {
	if conInstance.con != nil {
		conInstance.con.Close()
		return true
	}
	return false
}

func (conInstance *GoConnectFreeSwitch) sendCmd(cmd string) (bool, error) {
	if conInstance.con != nil {
		err := conInstance.con.Send(cmd)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, errors.New("not isset connect freeSwitch")
}

//var (
//	conn, errConnect   = goesl.NewClient("192.168.124.128", 8021, "ClueCon", 10)
//	conn2, errConnect2 = goesl.NewClient("192.168.124.128", 8021, "ClueCon", 10)
//)

func main() {
	var ListUser []UserLogin

	router := gin.Default()
	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))

	v1 := router.Group("/api")
	{
		v1.POST("/login", login(&ListUser))
		v1.POST("/callToCustomer", callToCustomer())
		v1.POST("/receiverNumber", receiverNumber(&ListUser))
		v1.POST("/check2FA", check2FA(&ListUser))
	}

	router.Run(":8000")

}

func assignDigitNumber(event *Event, list *[]UserLogin) error {
	num, err := strconv.Atoi(event.Digit)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	for i := range *list {
		if (*list)[i].PhoneNumber == event.NumberCalled {
			(*list)[i].Enter2FA = num
			break
		}
	}
	return nil
}
func isNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func receiverNumber(list *[]UserLogin) gin.HandlerFunc {
	return func(context *gin.Context) {
		var event Event

		if errorConvertData := context.ShouldBind(&event); errorConvertData != nil {
			println(errorConvertData.Error())
			println(1)
			context.JSON(http.StatusBadRequest, gin.H{
				"error": errorConvertData.Error(),
			})
			return
		}
		num, err := strconv.Atoi(event.Digit)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if num < 10 || num > 99 || event.NumberCalled == "" {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "loi validate number hoac phone number",
			})
			return
		}

		errAssignDigit := assignDigitNumber(&event, list)
		if errAssignDigit != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "loi validate so nhap vao",
			})
			return
		}

		return

	}
}

func callToCustomer() gin.HandlerFunc {
	return func(context *gin.Context) {
		var con = GoConnectFreeSwitch{}
		var userCheck2FA UserCheck2FA

		if errGetBodyData := context.ShouldBind(&userCheck2FA); errGetBodyData != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": errGetBodyData.Error(),
			})
			return
		}

		status, errorConnect := con.connect()
		if status == false {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": errorConnect.Error(),
			})
			return
		}

		cmd := "api originate user/" + userCheck2FA.PhoneNumber + " &lua(/usr/local/freeswitch/scripts/ivr.lua)"
		_, errorSend := con.sendCmd(cmd)
		if errorSend != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": errorSend.Error(),
			})
			return
		}
		statusCLose := con.closeConnect()
		if statusCLose == false {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "can not close connect freeSwitch",
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"success": "call success",
		})
		return
	}
}

func login(list *[]UserLogin) gin.HandlerFunc {
	return func(context *gin.Context) {
		var userBind UserBindData
		var userLogin UserLogin
		rand.Seed(time.Now().UnixNano())
		if err := context.ShouldBind(&userBind); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		current := time.Now().UTC()
		expired := current.Add(40 * time.Second)
		if !isNumeric(userBind.PhoneNumber) {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "phone validate fail",
			})
			return
		}
		userLogin.PhoneNumber = userBind.PhoneNumber
		userLogin.Password = userBind.Password
		userLogin.CreatedAt = &current
		userLogin.Expired = &expired
		userLogin.Pass2FA = rand.Intn(90) + 10
		userLogin.Enter2FA = 0
		userLogin.ChannelName = ""

		for i, v := range *list {
			if v.PhoneNumber == userLogin.PhoneNumber {
				*list = append((*list)[:i], (*list)[i+1:]...)
			}
		}
		*list = append(*list, userLogin)

		context.JSON(http.StatusOK, gin.H{
			"success":      true,
			"phone_number": userLogin.PhoneNumber,
			"number2FA":    userLogin.Pass2FA,
		})
		return
	}
}

func check2FA(list *[]UserLogin) gin.HandlerFunc {
	return func(context *gin.Context) {
		var user2FA UserCheck2FA
		if err := context.ShouldBind(&user2FA); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		current := time.Now().UTC()
		index := -1
		for i, v := range *list {
			createdAt := *v.CreatedAt
			expiredAt := *v.Expired
			if createdAt.Before(current) && expiredAt.After(current) {
				if v.PhoneNumber == user2FA.PhoneNumber {
					println(v.Pass2FA)
					println(v.Enter2FA)
					if v.Pass2FA == v.Enter2FA {
						*list = append((*list)[:i], (*list)[i+1:]...)
						context.JSON(http.StatusOK, gin.H{
							"success": true,
						})
						return
					}
					context.JSON(http.StatusBadRequest, gin.H{
						"success": false,
					})
					return
				} else {
					continue
				}
			} else {
				index = i
				break
			}
		}
		if index != -1 {
			*list = append((*list)[:index], (*list)[index+1:]...)
		}
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		return
	}

}
