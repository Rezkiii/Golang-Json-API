package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "sync"
)

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

var users []User
var mu sync.Mutex

// Load users from the JSON file
func loadUsers() {
    file, err := ioutil.ReadFile("users.json")
    if err != nil {
        log.Fatalf("Failed to read file: %v", err)
    }
    json.Unmarshal(file, &users)
}

// Save users to the JSON file
func saveUsers() {
    data, err := json.MarshalIndent(users, "", "  ")
    if err != nil {
        log.Fatalf("Failed to marshal users: %v", err)
    }
    ioutil.WriteFile("users.json", data, 0644)
}

// CORS handling
func handleCORS(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "http://192.168.43.154:3000") // Allow requests from React app
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // Handle preflight requests
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusNoContent)
        return
    }
}

// Get all users
func getUsers(w http.ResponseWriter, r *http.Request) {
    handleCORS(w, r) // Handle CORS
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// Add a new user
func postUser(w http.ResponseWriter, r *http.Request) {
    handleCORS(w, r) // Handle CORS
    var newUser User
    err := json.NewDecoder(r.Body).Decode(&newUser)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    mu.Lock()
    users = append(users, newUser)
    saveUsers()
    mu.Unlock()

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newUser)
}

// Login handler
func loginUser(w http.ResponseWriter, r *http.Request) {
    handleCORS(w, r) // Handle CORS
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Cek apakah username dan password ada di dalam daftar pengguna
    for _, u := range users {
        if u.Username == user.Username && u.Password == user.Password {
            w.WriteHeader(http.StatusOK) // Login berhasil
            return
        }
    }

    // Jika username atau password salah
    http.Error(w, "Invalid credentials", http.StatusUnauthorized) // Kembalikan 401
}

func main() {
    loadUsers()

    http.HandleFunc("/users", getUsers)
    http.HandleFunc("/users/add", postUser)      // Endpoint untuk menambah pengguna
    http.HandleFunc("/users/login", loginUser)    // Endpoint untuk login

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
