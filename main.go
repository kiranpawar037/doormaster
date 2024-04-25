package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DoorStatus struct {
	DoorID   string    `json:"door_id"`
	DoorName string    `json:"door_name"`
	Time     time.Time `json:"time"`
	Date     string    `json:"date"`
	Value    int       `json:"value"`
}

var doorStatus DoorStatus

func main() {
	r := gin.Default()

	baseURL := "http://localhost:8080"

	// Define endpoints
	r.POST("/door", GetDoorStatus)
	r.POST("/open-door/:doorID", OpenDoor)

	// Handlers for base URL and "/door" endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server is listening", "baseURL": baseURL})
	})
	r.GET("/door", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Door is also listening", "doorEndpoint": baseURL + "/door"})
	})

	r.Run(":8080")
}

func GetDoorStatus(c *gin.Context) {
	var newStatus DoorStatus
	if err := c.BindJSON(&newStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newStatus.Value != 0 && newStatus.Value != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value. Value must be 0 or 1."})
		return
	}

	doorStatus = newStatus

	var message string
	if newStatus.Value == 1 {
		message = "Door is open"
		// Start a Goroutine to close the door after 10 seconds
		go closeDoor(newStatus.DoorID, newStatus.DoorName)
	} else {
		message = "Door is closed"
	}

	response := gin.H{"message": message, "doorStatus": doorStatus}
	c.JSON(http.StatusOK, response)
}

func OpenDoor(c *gin.Context) {
	doorID := c.Param("doorID")
	// Here you can write code to open the door with the given doorID.
	// For now, let's assume the door is opened successfully.
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Door with ID %s is opened", doorID)})
}

func closeDoor(doorID, doorName string) {
	// Sleep for 10 seconds
	time.Sleep(10 * time.Second)

	// Update doorStatus to indicate the door is closed
	doorStatus.Value = 0
	doorStatus.Time = time.Now()
	doorStatus.Date = time.Now().Format("2006-01-02")
	message := "Door is closed"

	// Send the response indicating that the door is closed
	response := gin.H{"message": message, "doorStatus": doorStatus}
	fmt.Println("Response:", response)
}
