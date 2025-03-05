package main

import (
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
)
//struct for task
type Task struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

var (
	tasks Task
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	//CRUD
	r.POST("/addTask", addTask)
	r.GET("/getTasks", getTasks)
	r.PUT("/updateTask", updateTask)
	r.DELETE("/deleteTask/:id", deleteTask)

	r.Run(":8080")
}
//add new task
func addTask(c *gin.Context) {
	var newTask Task
	err := json.NewDecoder(c.Request.Body).Decode(&newTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input!"})
		return
	}

	mu.Lock()
	newTask.ID = nextID
	nextID++
	tasks = append(tasks, newTask)
	mu.Unlock()

	c.JSON(http.StatusCreated, newTask)
}
//get the task
func getTasks(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, tasks)
}
//update task
func updateTask(c *gin.Context) {
	var updatedTask Task
	err := json.NewDecoder(c.Request.Body).Decode(&updatedTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input!"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, task := range tasks {
		if task.ID == updatedTask.ID {
			tasks[i] = updatedTask
			c.JSON(http.StatusOK, updatedTask)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"Error": "Task not found!"})
}
//del task
func deleteTask(c *gin.Context) {
	id := c.Param("id")
	mu.Lock()
	defer mu.Unlock()

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"Notice": "Task deleted"})
			return
		}
	}
	//error msg
	c.JSON(http.StatusNotFound, gin.H{"Error": "Task not found!"})
}
