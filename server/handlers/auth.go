package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "time"
"log"
    "github.com/aws/aws-lambda-go/events"
    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_very_secret_key") // replace this securely!

type Claims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*") // adjust for stricter CORS if needed
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

// LoginHttpHandler for local Go dev server
func LoginHttpHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("LoginHttpHandler called")

    if r.Method == http.MethodOptions {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method != http.MethodPost {
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
        return
    }

    var creds struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
        return
    }

    if creds.Username == "admin" && creds.Password == "pass123" {
        tokenStr, _ := generateToken(creds.Username) // Assuming this exists
        writeJSON(w, http.StatusOK, map[string]string{"token": tokenStr})
        return
    }

    writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
}

// LoginHandler issues a token on valid login
func LoginHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var creds struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    json.Unmarshal([]byte(req.Body), &creds)

    if creds.Username == "admin" && creds.Password == "pass123" {
        tokenStr, _ := generateToken(creds.Username)
        return events.APIGatewayProxyResponse{
            StatusCode: http.StatusOK,
            Headers:    getCORSHeaders(),
            Body:       fmt.Sprintf(`{"token": "%s"}`, tokenStr),
        }, nil
    }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusUnauthorized,
        Headers:    getCORSHeaders(),
        Body:       `{"error": "invalid credentials"}`,
    }, nil
}

func generateToken(username string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)

    claims := &Claims{
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func VerifyJWT(authHeader string) error {
    if !strings.HasPrefix(authHeader, "Bearer ") {
        return fmt.Errorf("invalid token format")
    }

    tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil || !token.Valid {
        return fmt.Errorf("unauthorized")
    }

    return nil
}
