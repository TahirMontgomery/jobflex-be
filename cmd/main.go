package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/TahirMontgomery/jobflex-be/internal/auth"
	"github.com/TahirMontgomery/jobflex-be/internal/controllers"
	"github.com/TahirMontgomery/jobflex-be/internal/database"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var clientID = "CLIENT_ID"
var secretKey = "SECRET_KEY"
var domain = "https://jobflex.us.auth0.com/"

func initDb() {
	var err error
	database.DB, err = gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})
	if err != nil {
		panic("Error")
	}
	database.DB.SetupJoinTable(&database.Job{}, "Recruiters", &database.UserJob{})
	database.DB.SetupJoinTable(&database.User{}, "AssignedJobs", &database.UserJob{})

	database.DB.AutoMigrate(&database.User{}, &database.Company{}, &database.Job{}, &database.Application{},
		&database.ApplicantFile{}, &database.Benefits{}, &database.CustomMilestone{}, &database.EducationHistory{},
		&database.EmployerHistory{}, &database.Milestone{})
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://jobflex.us.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

var jwtChecker = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			panic(err.Error())
		}

		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	},
	SigningMethod: jwt.SigningMethodRS256,
	Extractor: func(r *http.Request) (string, error) {
		return r.Header.Get("Authorization"), nil
	},
})

func checkJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtMid := *jwtChecker
		if err := jwtMid.CheckJWT(c.Writer, c.Request); err != nil {
			c.AbortWithStatus(401)
		}

		accessToken := c.Request.Header.Get("Authorization")
		claims := jwt.MapClaims{}
		jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		})
		c.Set("uid", claims["sub"])
	}
}

func main() {
	fmt.Println("Hello guys")
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"*"}
	r.Use(cors.New(config))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "start",
		})
	})

	r.GET("/user/registration", controllers.CheckRegistration)
	r.POST("/user/register", checkJWT(), controllers.Register)
	r.POST("/user/plan", checkJWT(), controllers.RegisterPlan)

	initDb()

	data := auth.GetToken()

	auth.A = auth.Auth{
		ClientID:    clientID,
		SecretKey:   secretKey,
		URI:         domain,
		AccessToken: data["access_token"],
	}

	r.Run()
}
