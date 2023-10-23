package helper

import (
	"fmt"
	"net/http"
	"user/config"
	"user/models"

	"time"

	"github.com/gin-gonic/gin"
)

func StartMiddleware(c *gin.Context) {
	// before request
	fmt.Printf("%s request start path: %s time %v\n", c.Request.Method, c.Request.URL.Path, time.Now())
	c.Next()
}
func EndMiddleware(c *gin.Context) {
	c.Next()
	// after request
	fmt.Printf("%s request end path: %s time %v\n", c.Request.Method, c.Request.URL.Path, time.Now())
}
func LoggerAllInOne(c *gin.Context) {
	// before request
	t := time.Now()
	c.Next()
	// after request
	latency := time.Since(t)

	// access the status we are sending
	status := c.Writer.Status()

	fmt.Printf("Completed %s %s with status code %d in %v\n", c.Request.Method, c.Request.URL.Path, status, latency)
}
func Logger(c *gin.Context) {
	// before request
	beforeRequest(c)
	c.Next()
	// after request
	afterRequest(c)
}
func beforeRequest(c *gin.Context) {
	// before request
	t := time.Now()
	c.Set("start", t)
	c.Next()
}
func afterRequest(c *gin.Context) {
	// Get the start time from the request context
	startTime, exists := c.Get("start")
	if !exists {
		startTime = time.Now()
	}

	// Calculate the request duration
	duration := time.Since(startTime.(time.Time))

	// Log the request completion time and duration
	fmt.Printf("Completed %s %s in %v\n", c.Request.Method, c.Request.URL.Path, duration)
}

// task

func AuthMiddleware(c *gin.Context) {
	// Get the Authorization header from the request
	clientToken := c.Request.Header.Get("Authorization")
	if clientToken == "" {
		// If the Authorization header is not present, return a 403 status code
		c.JSON(403, "No Authorization header provided")
		c.Abort()
		return
	}

	CheckedInfo, err := ParseClaims(clientToken, config.JWTSecretKey)
	if err != nil {
		fmt.Println("error Parsing middleware :", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//set the claims in the context
	c.Set("user_info", CheckedInfo)
	// continue to the  next handler
	c.Next()
}

func PasswordMiddleware(c *gin.Context) {
	var user models.CreateUser
	fmt.Println("Before struct Middleware", user)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Retrieve the password from the parsed JSON body
	password := user.Password
	if !IsValidPassword(password) {
		// If the password is not valid, return a 400 status code with an error message
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		c.Abort()
		return
	}
	fmt.Println("Middleware struct holatda: ", user)
	// continue to the  next handler
	c.Next()
}
