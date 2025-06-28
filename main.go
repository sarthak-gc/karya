package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/fatih/color"
)

type Status string

const (
	PENDING Status = "PENDING"
	DONE    Status = "DONE"
)

const todosFile = "FILE_PATH_HERE(ABSOLUTE)"

type Todo struct {
	Id     int    `json:"id"`
	Task   string `json:"task"`
	Note   string `json:"note"`
	Status Status `json:"status"`
}

func addTodo(task string, note string) error {
	file, err := os.OpenFile(todosFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	newData := []string{strconv.Itoa(rand.Intn(90000) + 10000), task, note, string(PENDING)}
	fmt.Println(newData)
	if err := writer.Write(newData); err != nil {
		fmt.Println("Error writing record:", err)
		return err
	}
	fmt.Println(task, note)
	return nil
}

func listTodos() ([]Todo, error) {
	file, err := os.OpenFile(todosFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	fmt.Println(records)
	if err != nil {
		return nil, err
	}

	var todos []Todo
	for _, rec := range records {
		id, _ := strconv.Atoi(rec[0])
		todos = append(todos, Todo{
			Id:     id,
			Task:   rec[1],
			Note:   rec[2],
			Status: Status(rec[3]),
		})
	}

	return todos, nil
}
func markTodo(toActOn int64, action string) error {

	file, err := os.OpenFile(todosFile, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	var updated bool
	for _, rec := range records {
		id, _ := strconv.Atoi(rec[0])
		if id == int(toActOn) {
			writer := csv.NewWriter(file)
			defer writer.Flush()
			if action == "pending" {
				rec[3] = string(PENDING)
				updated = true
				break
			}
			if action == "done" {
				rec[3] = string(DONE)
				updated = true
				break
			}
		}
	}

	if !updated {
		return fmt.Errorf("invalid id")
	}

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("cannot truncate")
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("cannot seek file")
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("failed to write updated records: %v", err)
	}

	return nil

}
func deleteTodo(toActOn int64) error {
	file, err := os.OpenFile(todosFile, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	var updatedRecords [][]string
	found := false

	for _, rec := range records {
		id, _ := strconv.Atoi(rec[0])
		if id == int(toActOn) {
			found = true
			continue
		}
		updatedRecords = append(updatedRecords, rec)
	}
	if !found {
		return fmt.Errorf("no todo found")
	}
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(updatedRecords); err != nil {
		return err
	}

	return nil
}

func getNote(toActOn int64) (string, error) {
	file, err := os.OpenFile(todosFile, os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	for _, rec := range records {
		id, _ := strconv.Atoi(rec[0])
		if id == int(toActOn) {
			return rec[2], nil
		}
	}
	return "", err
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("No arguments")
		return
	}
	action := args[1]
	switch action {
	case "add":
		if len(args) < 2 {
			log.Fatal("No arguments")
			return
		}
		task := args[2]
		note := args[3]
		err := addTodo(task, note)
		if err != nil {
			log.Fatal("Error while adding the task to the file")
		}
		color.Green("Added")

	case "list":
		todos, err := listTodos()
		if err != nil {
			log.Fatal("Error listing the tasks")
		}
		if len(args) <= 2 {
			log.Fatal("No arguments")
			return
		}
		toActOn, err := strconv.ParseInt(args[2], 0, 0)
		if err == nil {
			for index, todo := range todos {
				if todo.Id == int(toActOn) {
					fmt.Println(todo)
					return
				}
				if index == len(todos)+1 {
					log.Fatal("Invalid id")
					return
				}
			}
			return
		}

		command := args[2]
		switch command {
		case "all":
			for _, todo := range todos {
				line := fmt.Sprintf("ID: %d | Task: %s | Note: %s | Status: %s",
					todo.Id, todo.Task, todo.Note, todo.Status)
				if todo.Status == DONE {
					color.New(color.FgBlack, color.BgGreen).Println(line)
				} else {
					color.New(color.FgBlack, color.BgYellow).Println(line)
				}
			}
		case "pending":
			for _, todo := range todos {
				if todo.Status == PENDING {
					line := fmt.Sprintf("ID: %d | Task: %s | Note: %s | Status: %s",
						todo.Id, todo.Task, todo.Note, todo.Status)
					color.New(color.FgBlack, color.BgYellow).Println(line)
				}
			}
		case "complete":
			for _, todo := range todos {
				if todo.Status == DONE {
					line := fmt.Sprintf("ID: %d | Task: %s | Note: %s | Status: %s",
						todo.Id, todo.Task, todo.Note, todo.Status)
					color.New(color.FgBlack, color.BgGreen).Println(line)
				}
			}
		}

	case "mark":
		if len(args) <= 3 {
			log.Fatal("No arguments")
			return

		}

		toActOn, err := strconv.ParseInt(args[2], 0, 0)
		if err != nil {
			log.Fatal("Invalid id")
		}
		status := args[3]

		err = markTodo(toActOn, status)
		if err == nil {
			color.New(color.FgBlack, color.BgGreen).Println("Marked successfully")
		}

	case "delete":
		toActOn, err := strconv.ParseInt(args[2], 0, 0)
		if err != nil {
			log.Fatal("Invalid arguments")
		}
		err = deleteTodo(toActOn)
		if err == nil {
			color.New(color.FgBlack, color.BgRed).Println("Deleted successfully")
		}

	case "get-note":
		toActOn, err := strconv.ParseInt(args[2], 0, 0)
		if err != nil {
			log.Fatal("Invalid arguments")
		}
		note, err := getNote(toActOn)
		if err != nil {
			log.Fatal("Error while reading file")
		}
		color.New(color.FgBlack, color.BgRed).Println(note)
	}
}
