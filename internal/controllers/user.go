package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TahirMontgomery/jobflex-be/internal/auth"
	"github.com/TahirMontgomery/jobflex-be/internal/database"
	"github.com/gin-gonic/gin"
)

var user database.User

// CheckRegistration Returns registration details
func CheckRegistration(c *gin.Context) {
	accessToken := c.Request.Header.Get("Authorization")
	query := c.Request.URL.Query()

	// fmt.Println(accessToken)
	// fmt.Println(query.Get("uid"))

	if accessToken == "" {
		c.JSON(400, gin.H{
			"error": "No access token provided",
		})
	}

	var sb strings.Builder
	sb.WriteString("Bearer ")
	sb.WriteString(accessToken)
	client := http.Client{}

	var url strings.Builder
	url.WriteString("https://jobflex.us.auth0.com/api/v2/users/")
	url.WriteString(query.Get("uid"))

	req, err := http.NewRequest("GET", url.String(), nil)
	req.Header.Add("Authorization", sb.String())
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	var data map[string]interface{}
	body, err := io.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &data)
	var registered bool
	var paidPlan bool

	for k := range data {
		if k == "app_metadata" {
			temp := data[k].(map[string]interface{})
			for k1 := range temp {
				if k1 == "registration" && temp[k1] == "complete" {
					registered = true
				}
				if k1 == "plan" && temp[k1] != nil {
					paidPlan = true
				}
			}
		}
	}

	if !registered || !paidPlan {
		c.JSON(200, gin.H{
			"error": gin.H{
				"code":       http.StatusForbidden,
				"registered": registered,
				"paidPlan":   paidPlan,
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
	return
}

// RegisterData model
type RegisterData struct {
	FirstName   string      `json:"firstName" binding:"required"`
	LastName    string      `json:"lastName" binding:"required"`
	CompanyName string      `json:"companyName" binding:"required"`
	CompanySize json.Number `json:"companySize" binding:"required"`
}

// Register register user to database
func Register(c *gin.Context) {
	var data RegisterData

	// accessToken := c.Request.Header.Get("Authorization")
	// fmt.Println(accessToken)
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(200, gin.H{
			"error": gin.H{
				"code":    http.StatusForbidden,
				"message": err,
			},
		})
	}

	user := database.User{ID: c.GetString("uid"), FirstName: data.FirstName, LastName: data.LastName}
	response := database.DB.Create(&user)
	if response.Error != nil {
		c.JSON(200, gin.H{
			"error": gin.H{
				"code":    http.StatusForbidden,
				"message": response.Error,
			},
		})
		return
	}

	company := database.Company{CompanyName: data.CompanyName, CompanySize: data.CompanySize}
	response = database.DB.Create(&company)
	if response.Error != nil {
		c.JSON(200, gin.H{
			"error": gin.H{
				"code":    http.StatusForbidden,
				"message": response.Error,
			},
		})
		return
	}

	err = database.DB.Model(&company).Association("Users").Append([]database.User{user})
	if err != nil {
		c.JSON(200, gin.H{
			"error": gin.H{
				"code":    http.StatusForbidden,
				"message": err,
			},
		})
		return
	}

	database.DB.Where("UID=?", user.ID).First(&user)

	payload := map[string]interface{}{
		"app_metadata": map[string]interface{}{
			"registration": "complete",
			"companyId":    user.CompanyID,
		},
	}

	resp, err1 := auth.A.Update(payload, user.ID)
	if err1 != nil {
		c.JSON(200, gin.H{
			"error": err1,
		})
	}

	fmt.Println(resp)

	c.JSON(200, gin.H{
		"success": true,
		"user":    user,
	})
	return
}

// PlanData struct
type PlanData struct {
	Plan string `json:"plan" binding:"required"`
}

func RegisterPlan(c *gin.Context) {
	var data PlanData
	uid := c.GetString("uid")

	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(200, gin.H{
			"error": gin.H{
				"code":    http.StatusForbidden,
				"message": err,
			},
		})
	}

	payload := map[string]interface{}{
		"app_metadata": map[string]interface{}{
			"plan": data.Plan,
		},
	}

	resp, err1 := auth.A.Update(payload, uid)
	if err1 != nil {
		c.JSON(200, gin.H{
			"error": err1,
		})
	}

	c.JSON(200, gin.H{
		"success": true,
		"user":    resp,
	})
}
