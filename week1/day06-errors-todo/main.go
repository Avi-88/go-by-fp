package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"time"
)

// Day 6: Errors, Custom Types, and Wrapping
//
// Read EXERCISE.md before starting.
// Parts 1-3 can live in this file.
// Part 4 (Todo CLI) will grow — feel free to split into multiple files.

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("Cannot divide by zero")
	}
	return a / b, nil
}

func sqrt(n float64) (float64, error) {
	if n < 0 {
		return 0, errors.New("Cannot find sqaure root of a negative number")
	}
	return math.Sqrt(n), nil
}

func parseAge(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("The passed value does not seem to be a number-%w", err)
	}
	if num < 0 || num > 150 {
		return 0, errors.New("The age value passed is invalid")
	}
	return num, nil
}

type ValidationError struct {
	Field   string
	Message string
}

type NotFoundError struct {
	Resource string
	ID       int
}

func (ve *ValidationError) Error() string {
	return "Error validating field -" + ve.Field + "error -" + ve.Message
}

func (nf *NotFoundError) Error() string {
	return "Resource -" + nf.Resource + "with ID -" + strconv.Itoa(nf.ID) + "not found"
}

func validateUser(name string, age int) error {
	if len(name) == 0 {
		return &ValidationError{Field: "name", Message: "Cant pass empty string as value"}
	}
	if age < 0 || age > 150 {
		return &ValidationError{Field: "age", Message: "Invalid age passed"}
	}
	return nil
}

func findUser(id int) (string, error) {
	users := map[int]string{1: "Ada", 2: "Eve", 3: "Dan", 4: "John"}
	name, ok := users[id]
	if !ok {
		return "", &NotFoundError{ID: id, Resource: "User"}
	}
	return name, nil
}

type User struct {
	name string
	age  int
}

// tries to open a file — if it fails, wrap: fmt.Errorf("readConfig: %w", err)
func readConfig(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("readConfig: %w", err)
	}
	defer f.Close()
	var cfg string
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return "", fmt.Errorf("readConfig: %w", err)
	}
	return cfg, nil
}

// calls readConfig — if it fails, wrap: fmt.Errorf("loadApp: %w", err)
func loadApp(configPath string) error {
	cfg, err := readConfig(configPath)
	if err != nil {
		return fmt.Errorf("loadApp: %w", err)
	}
	fmt.Printf("Loaded config - %s", cfg)
	return nil
}

// calls loadApp — if it fails, wrap: fmt.Errorf("startServer: %w", err)
func startServer(configPath string) error {
	if err := loadApp(configPath); err != nil {
		return fmt.Errorf("startServer: %w", err)
	}
	return nil
}

type Todo struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
}

func loadTodos() []Todo {
	var todos []Todo = []Todo{}
	f, err := os.Open("todos.json")
	if err != nil {
		fmt.Println("File not found - starting with empty list")
		return todos
	}
	if err := json.NewDecoder(f).Decode(&todos); err != nil {
		fmt.Println("Error loading todos - starting with empty list")
		return todos
	}
	return todos
}

func updateTodos(newTodos []Todo) error {
	f, err := os.OpenFile("todos.json", os.O_WRONLY| os.O_CREATE | os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file to update todos")
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(newTodos); err != nil {
		fmt.Println("Error persisting updated todos")
		return err
	}
	return nil
}

func add(todo string) {
	currTodos := loadTodos()
	var newId int
	if len(currTodos) == 0 {
		newId = 1
	} else {
		newId = currTodos[len(currTodos)-1].ID + 1
	}

	newTodo := Todo{
		ID:        newId,
		Text:      todo,
		Done:      false,
		CreatedAt: time.Now().String(),
	}
	currTodos = append(currTodos, newTodo)
	if err := updateTodos(currTodos); err != nil {
		fmt.Printf("Error adding the new todo - %s\n", err)
		return
	}
}

func list() {
	todoList := loadTodos()
	fmt.Println("------TODO LIST------")
	for _, i := range todoList {
		var status string = "[ ]"
		if i.Done {
			status = "[X]"
		}
		fmt.Printf("%d) %s - %s\n", i.ID, i.Text, status)
	}
}

func done(id int) {
	currTodos := loadTodos()
	idx := slices.IndexFunc(currTodos, func(todo Todo) bool {
		return todo.ID == id
	})
	if idx == -1 {
		fmt.Println("Invalid todo id - todo with the provided ID not found")
		return
	}
	currTodos[idx].Done = true
	if err := updateTodos(currTodos); err != nil {
		fmt.Printf("Error updating the todo status - %s\n", err)
		return
	}
}

func delete(id int) {
	currTodos := loadTodos()
	idx := slices.IndexFunc(currTodos, func(todo Todo) bool {
		return todo.ID == id
	})
	if idx == -1 {
		fmt.Println("Invalid todo id - todo with the provided ID not found")
		return
	}
	currTodos = slices.Delete(currTodos, idx, idx+1)
	if err := updateTodos(currTodos); err != nil {
		fmt.Printf("Error updating the todo status - %s\n", err)
		return
	}
}

func main() {
	// var user User = User{name: "", age: -1}
	// err := validateUser(user.name, user.age)
	// var ve *ValidationError
	// if errors.As(err, &ve){
	// 	fmt.Printf("Validation failed - field: %s, message: %s\n", ve.Field, ve.Message)
	// }else{
	// 	fmt.Printf("Validation successfull for user - %s, age - %d\n", user.name, user.age)
	// }

	// usr, err := findUser(8)
	// var nf *NotFoundError
	// if errors.As(err, &nf){
	// 	fmt.Printf("Not found - Resource:%s, ID:%d\n", nf.Resource, nf.ID)
	// }else{
	// 	fmt.Printf("Found user - %s\n", usr)
	// }

	// err := startServer("missing.json")
	// if errors.Is(err, os.ErrNotExist){
	// 	fmt.Println("The error is indeed file not exist")
	// }
	// fmt.Println(err)                           // startServer: loadApp: readConfig: ...
	// fmt.Println(errors.Unwrap(err))            // loadApp: readConfig: ...
	// fmt.Println(errors.Unwrap(errors.Unwrap(err)))  // readConfig: ...
	// fmt.Println(errors.Unwrap(errors.Unwrap(errors.Unwrap(err)))) // the raw os error
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: todo <add|list|done|delete> [args]")
		return
	}

	switch args[0] {
	case "add":
		if len(args) < 2 {
			fmt.Println("Invalid format - missing argument [task]")
			return
		}
		add(args[1])
		return

	case "list":
		list()
		return
	
	case "done":
		if len(args) < 2 {
			fmt.Println("Invalid format - missing argument [task]")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid id format passed - should be a number")
			return
		}
		done(id)
		return

	case "delete":
		if len(args) < 2 {
			fmt.Println("Invalid format - missing argument [task]")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid id format passed - should be a number")
			return
		}
		delete(id)
		return
	
	default:
		fmt.Printf("Unknown comamnd: %s\n", args[0])
		fmt.Println("Usage: todo <add|list|done|delete> [args]")
	}
}
