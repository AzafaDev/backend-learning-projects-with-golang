package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: task-cli <command> [argumetns]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli add <description>")
			os.Exit(1)
		}
		description := os.Args[2]
		if err := AddTask(description); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "list":
		filter := ""
		if len(os.Args) >= 3 {
			filter = os.Args[2]
		}
		if err := ListTask(filter); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli delete <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		if err := DeleteTask(id); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: task-cli update <description>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		description := os.Args[3]
		if err := UpdateTask(id, description); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli mark-in-progress <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		if err := MarkTask(id, "in-progress"); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli mark-done <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		if err := MarkTask(id, "done"); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}

func LoadTask() ([]Task, error) {
	data, err := os.ReadFile("tasks.json")
	if os.IsNotExist(err) {
		data = []byte("[]")
		if err := os.WriteFile("tasks.json", data, 0644); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func SaveTask(inputTasks []Task) error {
	data, err := json.MarshalIndent(inputTasks, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile("tasks.json", data, 0644)
}

func AddTask(description string) error {
	tasks, err := LoadTask()
	if err != nil {
		return err
	}

	maxID := 0
	now := time.Now()
	for _, task := range tasks {
		if task.Id > maxID {
			maxID = task.Id
		}
	}

	newTask := Task{
		Id:          maxID + 1,
		Description: description,
		Status:      "todo",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, newTask)

	if err := SaveTask(tasks); err != nil {
		return err
	}

	fmt.Printf("Task added successfully (ID: %d) \n", newTask.Id)

	return nil
}

func ListTask(filter string) error {
	data, err := LoadTask()
	if err != nil {
		return err
	}

	var tasks []Task
	for _, val := range data {
		if filter != "" && val.Status != filter {
			continue
		}
		tasks = append(tasks, val)
	}

	for _, task := range tasks {
		fmt.Printf("ID: %d, description: %s, status: %s\n", task.Id, task.Description, task.Status)
	}

	return nil
}

func DeleteTask(id int) error {
	data, err := LoadTask()
	if err != nil {
		return err
	}

	var newTasks []Task
	var found bool
	for _, v := range data {
		if v.Id == id {
			found = true
			continue
		}
		newTasks = append(newTasks, v)
	}

	if !found {
		return fmt.Errorf("task with id %d not found\n", id)
	}

	if err := SaveTask(newTasks); err != nil {
		return err
	}

	fmt.Printf("task with id (%d) is deleted successfully!", id)

	return nil
}

func UpdateTask(id int, description string) error {
	data, err := LoadTask()
	if err != nil {
		return err
	}

	var found bool

	for idx := range data {
		if data[idx].Id == id {
			data[idx].Description = description
			data[idx].UpdatedAt = time.Now()
			found = true
		}
	}

	if !found {
		return fmt.Errorf("id with (%d) not found", id)
	}

	if err := SaveTask(data); err != nil {
		return err
	}

	fmt.Printf("task with id %d updated successfully\n", id)

	return nil
}

func MarkTask(id int, status string) error {
	data, err := LoadTask()
	if err != nil {
		return err
	}

	var found bool

	for idx := range data {
		if data[idx].Id == id {
			data[idx].Status = status
			data[idx].UpdatedAt = time.Now()
			found = true
		}
	}

	if !found {
		return fmt.Errorf("id with (%d) not found", id)
	}

	if err := SaveTask(data); err != nil {
		return err
	}

	fmt.Printf("task with id %d marked %s successfully\n", id, status)

	return nil
}
