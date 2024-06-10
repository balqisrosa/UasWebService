package auth

import (
	"database/sql"
	"encoding/json"
	"onlineshop/database"
	"log"
	"net/http"
	"time"
	"strings"

	//"onlineshop/database"
	"onlineshop/model/user"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte ("your_secret_key")

type Claims struct{
	Username string `json:"username"`
	jwt.StandardClaims
}
func Registration(w http.ResponseWriter, r *http.Request) {
	var creds user.User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var existingUser user.User
	err = database.DB.QueryRow("SELECT username FROM users WHERE username = (?)", creds.Username).Scan(&existingUser.Username)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal server error select", http.StatusInternalServerError)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error password", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO users (username, password, email, phone_number, address) VALUES (?, ?, ?, ?, ?)",
		creds.Username, hashedPassword, creds.Email, creds.PhoneNumber, creds.Address)
	if err != nil {
		http.Error(w, "Internal server error insert", http.StatusInternalServerError)
		return
	}

	// Berikan Respon Sukses
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "User registered successfully",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func Login(w http.ResponseWriter, r *http.Request){
	var creds user.User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user user.User
	err = database.DB.QueryRow("SELECT user_id, username, password, email, phone_number, address FROM users WHERE username= (?)", creds.Username).Scan(&user.UserId, &user.Username, &user.Password,  &user.Email, &user.PhoneNumber, &user.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w,"User not found", http.StatusUnauthorized)
			return
		}
		http.Error(w,"Internal server error",
		http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid Password", http.StatusUnauthorized)
		return
	}
	expirationTime := time.Now().Add(60 * time.Minute)
    claims := &Claims{
        Username: creds.Username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
w.Header().Set("Content-Type", "application/json")
    response := map[string]interface{}{
        "message": "Login successful",
        "token": tokenString,
    }
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
    }
}
func ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return false, err
	}
	return token.Valid, nil
}
func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		sttArr := strings.Split(bearerToken, " ")
		if len(sttArr) == 2 {
			isValid, _ := ValidateToken(sttArr[1])
			if isValid {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}