package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"os"
)

// Day 5: Maps, Struct Tags, and JSON
//
// Read EXERCISE.md before starting.
// Implement all 4 parts in this file.

func WordCount(sentence string) map[string]int {
	m := make(map[string]int)
	words := strings.Fields(strings.ToLower(sentence))
	for _, w := range words {
		w = strings.Trim(w, ",.!?;:")
		m[w]++
	}
	return m
}

func TopN(counts map[string]int, n int) []string {
	type kv struct {
		key string
		value int
	}
	var ss []kv
	for k,v := range counts {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func (i,j int) bool {
		return ss[i].value > ss[j].value
	})

	var ans []string
	for idx, item := range ss {
		if idx == (n) {
			break
		}
		ans = append(ans, item.key)
	}
	return ans
}

func MergeCounts(a, b map[string]int) map[string]int {
	for k,v := range b {
		prev, ok := a[k]
		if !ok {
			prev = 0
		}
		a[k] = prev + v
	}
	return a
}

type Person struct {
    ID        int   	`json:"id"`
    FirstName string	`json:"first_name"`
    LastName  string	`json:"last_name"`
    Email     string	`json:"email"`
    Age       int		`json:"age,omitzero"`
    Active    bool		`json:"active"`
    Tags      []string	`json:"tags,omitempty"`
}

type Config struct {
    AppName  string
    Port     int
    Debug    bool
    Database struct {
        Host     string
        Port     int
        Name     string
        Password string
    }
    AllowedOrigins []string
}

func SaveConfig(cfg Config, path string) error {
	// serialized, err := json.MarshalIndent(cfg, "", "  ")
	// if err != nil {
	// 	return err
	// }
	// f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE, 0644)
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()
	// res, err := f.WriteString(string(serialized))
	// if err != nil || res < 0 {
	// 	return err
	// }
	// return nil
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		return err
	}
	return nil
}

func LoadConfig(path string) (Config, error) {
	// var cfg Config
	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	return cfg, err
	// }
	// dErr := json.Unmarshal(data, &cfg)
	// if dErr != nil {
	// 	return cfg, dErr
	// }
	// return cfg, nil
	var lCfg Config
	f, err := os.Open(path)
	if err != nil {
		return lCfg, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&lCfg); err != nil {
		return lCfg, err
	}
	return lCfg, nil
}

func main() {

	cfg := Config{
		AppName: "fafo",
		Port:    8080,
		Debug:   true,
		AllowedOrigins: []string{"http://localhost:3000", "https://example.com"},
	}
	cfg.Database.Host = "localhost"
	cfg.Database.Port = 5432
	cfg.Database.Name = "fafo_db"
	cfg.Database.Password = "secret"

	if err := SaveConfig(cfg, "config.json"); err != nil {
		fmt.Println("save error:", err)
		return
	}
	fmt.Println("Config saved.")

	loaded, err := LoadConfig("config.json")
	if err != nil {
		fmt.Println("load error:", err)
		return
	}
	fmt.Printf("AppName: %s\n", loaded.AppName)
	fmt.Printf("Port: %d\n", loaded.Port)
	fmt.Printf("Debug: %v\n", loaded.Debug)
	fmt.Printf("DB Host: %s\n", loaded.Database.Host)
	fmt.Printf("DB Port: %d\n", loaded.Database.Port)
	fmt.Printf("DB Name: %s\n", loaded.Database.Name)
	fmt.Printf("Allowed Origins: %v\n", loaded.AllowedOrigins)

	// counts := WordCount("the cat sat on the mat the cat")

	// top3 := TopN(counts, 3)

	// for k,v := range counts {
	// 	fmt.Printf("word:%s - count:%d\n", k, v)
	// }

	// fmt.Print("The top 3 words are:")
	// for _, i := range top3 {
	// 	fmt.Println(i)
	// }

	// a := map[string]int{"the": 3, "cat": 2, "sat": 1}
	// b := map[string]int{"cat": 5, "dog": 1, "the": 2}
	// merged := MergeCounts(a, b)
	// // Expected:
	// // "the" → 3 + 2 = 5  (shared key, counts should sum)
	// // "cat" → 2 + 5 = 7  (shared key, counts should sum)
	// // "sat" → 1          (only in a, should survive unchanged)
	// // "dog" → 1          (only in b, should be added)
	// fmt.Println(merged["the"])  // 5
	// fmt.Println(merged["cat"])  // 7
	// fmt.Println(merged["sat"])  // 1
	// fmt.Println(merged["dog"])  // 1


	// p := Person{
	// 	ID: 1, FirstName: "Ada", LastName: "Lovelace",
	// 	Email: "ada@example.com", Age: 36, Active: true,
	// 	Tags: []string{"engineer", "mathematician"},
	// }

	// // p1 := Person{
	// // 	ID: 1, FirstName: "Ada", LastName: "Lovelace",
	// // 	Email: "ada@example.com", Age: 0, Active: true,
	// // 	Tags: []string{},
	// // }

	// data, err := json.MarshalIndent(p, "", "  ")
	// if err != nil {
	// 	fmt.Print("There was a error while searalizing the payload")
	// }
	
	// var p1 Person

	// err1 := json.Unmarshal(data, &p1)
	// if err1 != nil {
	// 	fmt.Print("Error deserializing payload")
	// }
	// fmt.Print(p1)

	

}
