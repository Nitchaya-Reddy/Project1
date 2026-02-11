# UF Student Marketplace - Complete YouTube Tutorial Guide
## Building a Full-Stack Application Like a Professional Software Developer

---

# 🎬 VIDEO STRUCTURE

## Part 1: Introduction & Planning (15 min)
## Part 2: Backend Setup & Database (45 min)
## Part 3: API Development & Testing with Postman (60 min)
## Part 4: Frontend Setup & Architecture (30 min)
## Part 5: Building Frontend Components (60 min)
## Part 6: Integration & Final Testing (30 min)

**Total Estimated Time: 4 hours of content**

---

# 📋 PART 1: INTRODUCTION & PLANNING

## What We're Building
"Today we're building a student marketplace - think of it like Facebook Marketplace but specifically for university students. Users can create accounts, list items for sale, browse listings, and chat with sellers."

## Technology Choices & Why

### Backend: Go with Gin Framework

**Script:**
"For our backend, we're using Go with the Gin framework. Let me explain why:

**Why Go?**
- Fast compilation and execution
- Built-in concurrency support
- Single binary deployment (no runtime needed)
- Strong typing catches errors at compile time

**Why Gin specifically?**
- Fastest Go web framework (uses httprouter under the hood)
- Minimal boilerplate
- Great middleware support
- Large community

**Alternatives you could use instead:**
1. **Node.js with Express** - More developers know JavaScript, larger npm ecosystem
2. **Python with FastAPI** - Excellent for rapid prototyping, automatic OpenAPI docs
3. **Java with Spring Boot** - Enterprise-grade, great for large teams
4. **Rust with Actix** - Even faster than Go, but steeper learning curve

**Most Optimized Choice:** For a startup/small project, **FastAPI** would be faster to develop. For high-performance production, **Go with Gin** or **Rust with Actix** are better."

### Database: SQLite with GORM

**Script:**
"For our database, we're using SQLite with GORM as our ORM.

**Why SQLite?**
- Zero configuration - it's just a file
- Perfect for development and small to medium apps
- No separate database server needed

**Why GORM?**
- Write Go code instead of raw SQL
- Auto-migrations create tables from structs
- Prevents SQL injection automatically
- Handles relationships easily

**Alternatives:**
1. **PostgreSQL** - Better for production, supports more concurrent writes
2. **MySQL** - Industry standard, great performance
3. **MongoDB** - If you need flexible schemas
4. **Raw SQL queries** - More control, better performance, but more code

**Most Optimized Choice:** For production, use **PostgreSQL**. For learning/development, **SQLite** is perfect."

### Frontend: Angular 17

**Script:**
"For the frontend, we're using Angular 17.

**Why Angular?**
- Full framework with everything included
- TypeScript by default (catches errors)
- Dependency injection built-in
- Great for large applications

**Alternatives:**
1. **React** - More flexible, larger job market
2. **Vue.js** - Easier learning curve
3. **Svelte** - Smaller bundle sizes, faster
4. **Next.js** - React with SSR built-in

**Most Optimized Choice:** For a small project, **Vue.js** or **Svelte** would be faster to build. For enterprise apps, **Angular** or **Next.js** are better."

---

# 📋 PART 2: BACKEND SETUP & DATABASE

## Step 1: Project Setup

**Script:**
"Let's start by creating our project structure. Open your terminal."

```bash
# Create project folder
mkdir uf-marketplace
cd uf-marketplace

# Create backend folder
mkdir backend
cd backend

# Initialize Go module
go mod init uf-marketplace
```

**Explain:**
"`go mod init` creates a `go.mod` file which tracks our dependencies - like `package.json` in Node.js."

## Step 2: Install Dependencies

```bash
# Gin - Web framework
go get -u github.com/gin-gonic/gin

# GORM - ORM for database
go get -u gorm.io/gorm

# SQLite driver for GORM
go get -u gorm.io/driver/sqlite

# JWT for authentication
go get -u github.com/golang-jwt/jwt/v5

# Bcrypt for password hashing
go get -u golang.org/x/crypto/bcrypt

# CORS middleware
go get -u github.com/gin-contrib/cors
```

**Explain each:**
- **Gin**: Handles HTTP requests/responses
- **GORM**: Maps Go structs to database tables
- **SQLite driver**: Connects GORM to SQLite
- **JWT**: Creates secure tokens for authentication
- **Bcrypt**: Hashes passwords (never store plain text!)
- **CORS**: Allows frontend on different port to call backend

## Step 3: Create Folder Structure

```bash
mkdir -p database handlers middleware models utils
```

**Script:**
"This is a standard Go project structure:
- `database/` - Database connection code
- `handlers/` - Request handlers (controllers in MVC)
- `middleware/` - Code that runs before handlers
- `models/` - Data structures
- `utils/` - Helper functions"

---

## Step 4: Create User Model

**File: `models/user.go`**

**Script:**
"Let's create our first model. I'll explain every line."

```go
package models
// ^ Every Go file starts with a package declaration
// All files in 'models' folder share the 'models' package

import (
    "time"
    "gorm.io/gorm"
)
// ^ Import statements - like 'require' in Node.js

// User represents a user in our system
type User struct {
    // ^ 'type' creates a new type, 'struct' is like a class in other languages
    
    gorm.Model
    // ^ This embeds GORM's Model which adds:
    // - ID uint (primary key, auto-increment)
    // - CreatedAt time.Time (set on insert)
    // - UpdatedAt time.Time (set on update)
    // - DeletedAt gorm.DeletedAt (soft delete - record hidden, not removed)
    
    Email string `gorm:"unique;not null" json:"email"`
    // ^ Field declaration: Name Type Tags
    // - 'gorm:"unique;not null"' - database constraints
    // - 'json:"email"' - JSON field name when serializing
    // Why unique? Each user needs a unique email to login
    
    Password string `gorm:"not null" json:"-"`
    // ^ json:"-" means this field is NEVER sent in JSON responses
    // SECURITY: Never expose password hashes to clients!
    
    FirstName string `gorm:"not null" json:"first_name"`
    LastName  string `gorm:"not null" json:"last_name"`
    // Note: json uses snake_case (JavaScript convention)
    // Go uses PascalCase (Go convention)
    
    ProfileImage string `json:"profile_image"`
    // ^ No 'not null' - this field is optional
    
    Phone string `json:"phone"`
    Bio   string `json:"bio"`
    
    IsAdmin bool `gorm:"default:false" json:"is_admin"`
    // ^ default:false sets the default value in database
}

// ToResponse creates a safe response without sensitive data
func (u *User) ToResponse() map[string]interface{} {
    // ^ Method on User struct
    // (u *User) is the "receiver" - like 'this' or 'self'
    // *User means pointer (we can modify the original)
    
    return map[string]interface{}{
        // ^ Generic map type - like Object in JavaScript
        "id":            u.ID,
        "email":         u.Email,
        "first_name":    u.FirstName,
        "last_name":     u.LastName,
        "profile_image": u.ProfileImage,
        "phone":         u.Phone,
        "bio":           u.Bio,
        "is_admin":      u.IsAdmin,
        "created_at":    u.CreatedAt,
    }
    // Notice: Password is not included!
}
```

**Alternative Approaches:**

1. **Without GORM (raw SQL):**
```go
type User struct {
    ID        int64
    Email     string
    // ... etc
}

// You'd write SQL manually:
db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, hashedPassword)
```
*Pros: More control, slightly faster*
*Cons: More code, risk of SQL injection if not careful*

2. **Using different ORM (sqlx):**
```go
import "github.com/jmoiron/sqlx"
// sqlx is lighter weight, closer to raw SQL but with helpers
```

---

## Step 5: Create Listing Model

**File: `models/listing.go`**

```go
package models

import "gorm.io/gorm"

// Status constants - use const instead of magic strings
const (
    StatusActive   = "active"
    StatusSold     = "sold"
    StatusInactive = "inactive"
)
// ^ Why constants? Prevents typos like "actve" and enables IDE autocomplete

// Category represents a listing category
type Category struct {
    gorm.Model
    Name        string `gorm:"unique;not null" json:"name"`
    Description string `json:"description"`
    Icon        string `json:"icon"` // Material icon name
}

// Listing represents an item for sale
type Listing struct {
    gorm.Model
    Title       string  `gorm:"not null" json:"title"`
    Description string  `json:"description"`
    Price       float64 `gorm:"not null" json:"price"`
    // ^ float64 for currency - in production, use int (cents) to avoid floating point errors
    // Example: $19.99 stored as 1999 cents
    
    CategoryID uint     `gorm:"not null" json:"category_id"`
    Category   Category `json:"category"`
    // ^ This creates a foreign key relationship
    // GORM automatically joins Category when we use Preload()
    
    SellerID uint `gorm:"not null" json:"seller_id"`
    Seller   User `json:"seller"`
    // ^ Another foreign key - who's selling this item
    
    Images []ListingImage `json:"images"`
    // ^ One-to-many relationship: one listing has many images
    
    Status    string `gorm:"default:active" json:"status"`
    Condition string `json:"condition"` // new, like_new, good, fair, poor
    Location  string `json:"location"`  // Pickup location
    Views     int    `gorm:"default:0" json:"views"`
}

// ListingImage stores images for a listing
type ListingImage struct {
    gorm.Model
    ListingID uint   `gorm:"not null" json:"listing_id"`
    ImageURL  string `gorm:"not null" json:"image_url"`
    IsPrimary bool   `gorm:"default:false" json:"is_primary"`
}
```

**Relationships Explained:**
```
User 1 ──────< Listing (one user has many listings)
Category 1 ──────< Listing (one category has many listings)
Listing 1 ──────< ListingImage (one listing has many images)
```

---

## Step 6: Create Chat Model

**File: `models/chat.go`**

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

// Chat represents a conversation between buyer and seller
type Chat struct {
    gorm.Model
    ListingID uint    `gorm:"not null" json:"listing_id"`
    Listing   Listing `json:"listing"`
    BuyerID   uint    `gorm:"not null" json:"buyer_id"`
    Buyer     User    `json:"buyer"`
    SellerID  uint    `gorm:"not null" json:"seller_id"`
    Seller    User    `json:"seller"`
}
// ^ A chat is about a specific listing, between a buyer and seller

// Message represents a single message in a chat
type Message struct {
    gorm.Model
    ChatID   uint   `gorm:"not null" json:"chat_id"`
    SenderID uint   `gorm:"not null" json:"sender_id"`
    Sender   User   `json:"sender"`
    Content  string `gorm:"not null" json:"content"`
    IsRead   bool   `gorm:"default:false" json:"is_read"`
    ReadAt   *time.Time `json:"read_at,omitempty"`
    // ^ Pointer (*time.Time) because it can be null (unread)
    // 'omitempty' - don't include in JSON if nil
}

// ChatResponse is what we send to the frontend
type ChatResponse struct {
    ID           uint      `json:"id"`
    ListingID    uint      `json:"listing_id"`
    ListingTitle string    `json:"listing_title"`
    ListingImage string    `json:"listing_image"`
    OtherUserID  uint      `json:"other_user_id"`
    OtherUserName string   `json:"other_user_name"`
    LastMessage  string    `json:"last_message"`
    LastMessageAt time.Time `json:"last_message_at"`
    UnreadCount  int       `json:"unread_count"`
}
// ^ Custom response type - formats data for frontend's needs
```

---

## Step 7: Database Connection

**File: `database/database.go`**

```go
package database

import (
    "log"
    "uf-marketplace/models"
    
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// DB is the global database instance
var DB *gorm.DB
// ^ Global variable - accessible from other packages
// In larger apps, use dependency injection instead

// InitDB initializes the database connection
func InitDB() {
    var err error
    
    // Open SQLite database
    DB, err = gorm.Open(sqlite.Open("marketplace.db"), &gorm.Config{})
    // ^ sqlite.Open() creates/opens the file
    // &gorm.Config{} is empty config - use defaults
    
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
        // ^ log.Fatal prints error and exits program
        // In production, handle this more gracefully
    }
    
    // Auto-migrate creates/updates tables based on structs
    err = DB.AutoMigrate(
        &models.User{},
        &models.Category{},
        &models.Listing{},
        &models.ListingImage{},
        &models.Chat{},
        &models.Message{},
        &models.Notification{},
    )
    // ^ AutoMigrate:
    // - Creates tables if they don't exist
    // - Adds new columns if you add fields to struct
    // - Does NOT delete columns or change types (safe)
    
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }
    
    // Seed initial categories
    seedCategories()
}

func seedCategories() {
    var count int64
    DB.Model(&models.Category{}).Count(&count)
    // ^ Count how many categories exist
    
    if count == 0 {
        // Only seed if database is empty
        categories := []models.Category{
            {Name: "Textbooks", Description: "Academic textbooks", Icon: "book"},
            {Name: "Electronics", Description: "Phones, laptops", Icon: "devices"},
            {Name: "Furniture", Description: "Dorm furniture", Icon: "chair"},
            // ... more categories
        }
        
        DB.Create(&categories)
        // ^ Batch insert all categories
        log.Println("Seeded categories")
    }
}
```

**Alternative: Using environment variables for config:**
```go
import "os"

func InitDB() {
    dbPath := os.Getenv("DATABASE_URL")
    if dbPath == "" {
        dbPath = "marketplace.db" // Default fallback
    }
    DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}
```

---

# 📋 PART 3: API DEVELOPMENT & TESTING WITH POSTMAN

## Step 8: JWT Utilities

**File: `utils/jwt.go`**

```go
package utils

import (
    "errors"
    "os"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
)

// Get secret from environment, with fallback for development
var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        // Default for development - CHANGE IN PRODUCTION!
        return "your-secret-key-change-in-production"
    }
    return secret
}

// Claims is the JWT payload structure
type Claims struct {
    UserID  uint   `json:"user_id"`
    Email   string `json:"email"`
    IsAdmin bool   `json:"is_admin"`
    jwt.RegisteredClaims
    // ^ Embeds standard claims: exp, iat, nbf, etc.
}

// GenerateToken creates a new JWT for a user
func GenerateToken(userID uint, email string, isAdmin bool) (string, error) {
    // Create claims with expiration
    claims := Claims{
        UserID:  userID,
        Email:   email,
        IsAdmin: isAdmin,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            // ^ Token expires in 24 hours
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            // ^ Token is valid starting now
        },
    }
    
    // Create token with claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    // ^ HS256 = HMAC-SHA256 signing algorithm
    // Alternative: RS256 (RSA) for asymmetric keys
    
    // Sign and return the token string
    return token.SignedString(jwtSecret)
}

// ValidateToken checks if a token is valid and returns claims
func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(
        tokenString,
        &Claims{},
        func(token *jwt.Token) (interface{}, error) {
            // Validate signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, errors.New("invalid signing method")
            }
            return jwtSecret, nil
        },
    )
    
    if err != nil {
        return nil, err
    }
    
    // Type assert to get our Claims struct
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }
    
    return claims, nil
}
```

**How JWT Works (explain in video):**
```
1. User logs in with email/password
2. Server validates credentials
3. Server creates JWT with user info + expiration
4. Server sends JWT to client
5. Client stores JWT (localStorage)
6. Client sends JWT in Authorization header for each request
7. Server validates JWT and extracts user info

JWT Structure:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.  <- Header (base64)
eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3QuLi  <- Payload (base64)
.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV...   <- Signature (HMAC)
```

**Alternatives:**
1. **Session-based auth** - Store session ID in cookie, session data on server
2. **OAuth 2.0** - Use Google/Facebook login
3. **Passport.js** (if using Node.js) - Auth middleware with strategies

---

## Step 9: Authentication Middleware

**File: `middleware/auth.go`**

```go
package middleware

import (
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
    "uf-marketplace/utils"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
    // ^ Returns a function - this is the "middleware" pattern
    // Middleware runs BEFORE the actual handler
    
    return func(c *gin.Context) {
        // c = Context - contains request/response data
        
        // Step 1: Get Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header required",
            })
            c.Abort() // Stop processing - don't call next handlers
            return
        }
        
        // Step 2: Extract token from "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid authorization format. Use: Bearer <token>",
            })
            c.Abort()
            return
        }
        
        tokenString := parts[1]
        
        // Step 3: Validate token
        claims, err := utils.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
            })
            c.Abort()
            return
        }
        
        // Step 4: Store user info in context for handlers to use
        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        c.Set("is_admin", claims.IsAdmin)
        // ^ Like adding properties to request object in Express.js
        
        // Step 5: Continue to next handler
        c.Next()
    }
}

// OptionalAuthMiddleware extracts user if token present, but doesn't require it
func OptionalAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader != "" {
            parts := strings.Split(authHeader, " ")
            if len(parts) == 2 && parts[0] == "Bearer" {
                claims, err := utils.ValidateToken(parts[1])
                if err == nil {
                    c.Set("user_id", claims.UserID)
                    c.Set("email", claims.Email)
                }
            }
        }
        c.Next() // Always continue, even without auth
    }
}
```

**Middleware Execution Flow:**
```
Request → CORS Middleware → Auth Middleware → Handler → Response
              ↓                    ↓              ↓
         Add headers         Validate JWT    Business logic
```

---

## Step 10: Authentication Handlers

**File: `handlers/auth.go`**

```go
package handlers

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    
    "uf-marketplace/database"
    "uf-marketplace/models"
    "uf-marketplace/utils"
)

// RegisterInput defines the expected JSON body
type RegisterInput struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=6"`
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
}
// ^ 'binding' tags are validated by Gin automatically
// required = field must be present
// email = must be valid email format
// min=6 = minimum 6 characters

// Register creates a new user account
func Register(c *gin.Context) {
    var input RegisterInput
    
    // Parse and validate JSON body
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // ^ ShouldBindJSON does:
    // 1. Parse JSON from request body
    // 2. Map to struct fields
    // 3. Validate based on binding tags
    // If any fail, err contains details
    
    // Check if email already exists
    var existingUser models.User
    result := database.DB.Where("email = ?", input.Email).First(&existingUser)
    // ^ Where() builds: SELECT * FROM users WHERE email = ?
    // First() executes and gets first match
    // The ? is a parameterized query - prevents SQL injection!
    
    if result.Error == nil {
        // No error means user was found = email exists
        c.JSON(http.StatusConflict, gin.H{
            "error": "Email already registered",
        })
        return
    }
    
    // Hash password with bcrypt
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(input.Password), 
        bcrypt.DefaultCost, // Cost factor 10 - balance of security/speed
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Error hashing password",
        })
        return
    }
    // ^ bcrypt automatically adds a random salt
    // Same password hashed twice = different hashes
    // This protects against rainbow table attacks
    
    // Create user
    user := models.User{
        Email:     input.Email,
        Password:  string(hashedPassword),
        FirstName: input.FirstName,
        LastName:  input.LastName,
    }
    
    if result := database.DB.Create(&user); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Error creating user",
        })
        return
    }
    // ^ Create() inserts into database
    // User.ID is automatically set after insert
    
    // Generate JWT token
    token, err := utils.GenerateToken(user.ID, user.Email, user.IsAdmin)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Error generating token",
        })
        return
    }
    
    // Return success response
    c.JSON(http.StatusCreated, gin.H{
        "token": token,
        "user":  user.ToResponse(), // Excludes password
    })
}

// LoginInput for login request
type LoginInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// Login authenticates a user
func Login(c *gin.Context) {
    var input LoginInput
    
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Find user by email
    var user models.User
    result := database.DB.Where("email = ?", input.Email).First(&user)
    if result.Error != nil {
        // User not found - but don't reveal this!
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Invalid credentials",
        })
        return
    }
    // ^ SECURITY: Same error for "user not found" and "wrong password"
    // Prevents attackers from knowing which emails are registered
    
    // Compare password with hash
    err := bcrypt.CompareHashAndPassword(
        []byte(user.Password), 
        []byte(input.Password),
    )
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Invalid credentials",
        })
        return
    }
    
    // Generate token
    token, _ := utils.GenerateToken(user.ID, user.Email, user.IsAdmin)
    
    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user":  user.ToResponse(),
    })
}

// GetMe returns the current authenticated user
func GetMe(c *gin.Context) {
    // Get user ID from context (set by auth middleware)
    userID := c.GetUint("user_id")
    // ^ c.GetUint retrieves value set by c.Set()
    
    var user models.User
    if result := database.DB.First(&user, userID); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    c.JSON(http.StatusOK, user.ToResponse())
}
```

---

## 🧪 TESTING WITH POSTMAN

**Script:**
"Now let's test our API with Postman. Postman is an industry-standard tool for API testing."

### Setting Up Postman

1. **Create a new Collection** called "UF Marketplace"

2. **Create environment variables:**
   - `base_url`: `http://localhost:8080/api`
   - `token`: (will be set after login)

3. **Create requests:**

### Test 1: Register User
```
POST {{base_url}}/auth/register

Body (JSON):
{
    "email": "test@ufl.edu",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
}

Expected Response (201 Created):
{
    "token": "eyJhbGc...",
    "user": {
        "id": 1,
        "email": "test@ufl.edu",
        "first_name": "John",
        "last_name": "Doe"
    }
}
```

**In Postman Tests tab, add:**
```javascript
// Save token to environment
if (pm.response.code === 201) {
    var jsonData = pm.response.json();
    pm.environment.set("token", jsonData.token);
}

// Test assertions
pm.test("Status is 201", function() {
    pm.response.to.have.status(201);
});

pm.test("Response has token", function() {
    pm.expect(pm.response.json()).to.have.property("token");
});
```

### Test 2: Register Duplicate Email
```
POST {{base_url}}/auth/register

Body (same as above)

Expected Response (409 Conflict):
{
    "error": "Email already registered"
}
```

### Test 3: Login
```
POST {{base_url}}/auth/login

Body (JSON):
{
    "email": "test@ufl.edu",
    "password": "password123"
}

Expected Response (200 OK):
{
    "token": "eyJhbGc...",
    "user": {...}
}
```

### Test 4: Get Current User
```
GET {{base_url}}/auth/me

Headers:
Authorization: Bearer {{token}}

Expected Response (200 OK):
{
    "id": 1,
    "email": "test@ufl.edu",
    ...
}
```

### Test 5: Access Protected Route Without Token
```
GET {{base_url}}/auth/me

(No Authorization header)

Expected Response (401 Unauthorized):
{
    "error": "Authorization header required"
}
```

---

## Step 11: Listing Handlers

**File: `handlers/listing.go`**

```go
package handlers

import (
    "fmt"
    "math"
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    
    "uf-marketplace/database"
    "uf-marketplace/models"
)

// CreateListingInput for creating listings
type CreateListingInput struct {
    Title       string   `json:"title" binding:"required"`
    Description string   `json:"description"`
    Price       float64  `json:"price" binding:"required,gt=0"`
    // ^ gt=0 means "greater than 0"
    CategoryID  uint     `json:"category_id" binding:"required"`
    Condition   string   `json:"condition"`
    Location    string   `json:"location"`
    Images      []string `json:"images"` // Array of image URLs
}

// CreateListing creates a new listing
func CreateListing(c *gin.Context) {
    userID := c.GetUint("user_id") // From auth middleware
    
    var input CreateListingInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Verify category exists
    var category models.Category
    if result := database.DB.First(&category, input.CategoryID); result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
        return
    }
    
    // Create listing
    listing := models.Listing{
        Title:       input.Title,
        Description: input.Description,
        Price:       input.Price,
        CategoryID:  input.CategoryID,
        SellerID:    userID,
        Condition:   input.Condition,
        Location:    input.Location,
        Status:      models.StatusActive,
    }
    
    if result := database.DB.Create(&listing); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Error creating listing",
        })
        return
    }
    
    // Create image records
    for i, imageURL := range input.Images {
        image := models.ListingImage{
            ListingID: listing.ID,
            ImageURL:  imageURL,
            IsPrimary: i == 0, // First image is primary
        }
        database.DB.Create(&image)
    }
    
    // Reload with relationships
    database.DB.Preload("Images").
        Preload("Category").
        Preload("Seller").
        First(&listing, listing.ID)
    // ^ Preload() eagerly loads related data
    // Without Preload, listing.Category would be empty
    
    c.JSON(http.StatusCreated, listing)
}

// GetListings returns paginated, filtered listings
func GetListings(c *gin.Context) {
    // Parse query parameters with defaults
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    // ^ DefaultQuery returns default if param not present
    // strconv.Atoi converts string to int
    
    search := c.Query("search")
    categoryID := c.Query("category_id")
    minPrice := c.Query("min_price")
    maxPrice := c.Query("max_price")
    condition := c.Query("condition")
    sortBy := c.DefaultQuery("sort", "created_at")
    sortOrder := c.DefaultQuery("order", "desc")
    
    // Calculate offset for pagination
    offset := (page - 1) * limit
    // Page 1: offset 0, Page 2: offset 20, etc.
    
    // Build query dynamically
    query := database.DB.Model(&models.Listing{}).
        Where("status = ?", models.StatusActive)
    // ^ Only show active listings by default
    
    // Apply filters conditionally
    if search != "" {
        // Search in title and description
        searchPattern := "%" + search + "%"
        query = query.Where(
            "title LIKE ? OR description LIKE ?",
            searchPattern, searchPattern,
        )
        // ^ LIKE with % is SQL pattern matching
        // %word% = contains "word" anywhere
    }
    
    if categoryID != "" {
        query = query.Where("category_id = ?", categoryID)
    }
    
    if minPrice != "" {
        query = query.Where("price >= ?", minPrice)
    }
    
    if maxPrice != "" {
        query = query.Where("price <= ?", maxPrice)
    }
    
    if condition != "" {
        query = query.Where("condition = ?", condition)
    }
    
    // Count total before pagination (for frontend's pagination UI)
    var total int64
    query.Count(&total)
    
    // Execute query with pagination and sorting
    var listings []models.Listing
    query.Preload("Images").
        Preload("Category").
        Preload("Seller").
        Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
        Offset(offset).
        Limit(limit).
        Find(&listings)
    // ^ Order() adds ORDER BY clause
    // Offset() skips records
    // Limit() restricts number of records
    // Find() returns all matches (vs First() = one)
    
    c.JSON(http.StatusOK, gin.H{
        "listings": listings,
        "total":    total,
        "page":     page,
        "limit":    limit,
        "pages":    int(math.Ceil(float64(total) / float64(limit))),
    })
}

// GetListing returns a single listing by ID
func GetListing(c *gin.Context) {
    id := c.Param("id")
    // ^ Gets :id from URL path /listings/:id
    
    var listing models.Listing
    result := database.DB.Preload("Images").
        Preload("Category").
        Preload("Seller").
        First(&listing, id)
    
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
        return
    }
    
    // Increment view count
    database.DB.Model(&listing).Update("views", listing.Views+1)
    listing.Views++ // Update local copy too
    
    c.JSON(http.StatusOK, listing)
}

// UpdateListing updates a listing
func UpdateListing(c *gin.Context) {
    id := c.Param("id")
    userID := c.GetUint("user_id")
    
    var listing models.Listing
    if result := database.DB.First(&listing, id); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
        return
    }
    
    // Check ownership
    if listing.SellerID != userID {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "Not authorized to update this listing",
        })
        return
    }
    
    // Parse update data
    var input map[string]interface{}
    c.ShouldBindJSON(&input)
    // ^ Using map allows partial updates (only sent fields)
    
    // Update
    database.DB.Model(&listing).Updates(input)
    // ^ Updates() only changes provided fields
    
    // Reload and return
    database.DB.Preload("Images").
        Preload("Category").
        Preload("Seller").
        First(&listing, id)
    
    c.JSON(http.StatusOK, listing)
}

// DeleteListing soft-deletes a listing
func DeleteListing(c *gin.Context) {
    id := c.Param("id")
    userID := c.GetUint("user_id")
    
    var listing models.Listing
    if result := database.DB.First(&listing, id); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
        return
    }
    
    // Check ownership
    if listing.SellerID != userID {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "Not authorized to delete this listing",
        })
        return
    }
    
    // Soft delete (sets deleted_at, doesn't remove from DB)
    database.DB.Delete(&listing)
    // ^ With gorm.Model, this is a soft delete
    // To hard delete: database.DB.Unscoped().Delete(&listing)
    
    c.JSON(http.StatusOK, gin.H{"message": "Listing deleted"})
}

// GetCategories returns all categories
func GetCategories(c *gin.Context) {
    var categories []models.Category
    database.DB.Find(&categories)
    c.JSON(http.StatusOK, categories)
}
```

---

## Step 12: Main Application Entry Point

**File: `main.go`**

```go
package main

import (
    "log"
    "os"
    "strings"
    
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    
    "uf-marketplace/database"
    "uf-marketplace/handlers"
    "uf-marketplace/middleware"
)

func main() {
    // Initialize database
    database.InitDB()
    log.Println("Database initialized")
    
    // Create uploads directory for images
    os.MkdirAll("./uploads", 0755)
    // ^ 0755 = owner can read/write/execute, others can read/execute
    
    // Create Gin router
    r := gin.Default()
    // ^ Default() includes Logger and Recovery middleware
    // Logger logs all requests
    // Recovery catches panics and returns 500
    
    // Configure CORS (Cross-Origin Resource Sharing)
    config := cors.DefaultConfig()
    
    // Get allowed origins from environment or use default
    allowedOrigins := os.Getenv("CORS_ORIGINS")
    if allowedOrigins != "" {
        config.AllowOrigins = strings.Split(allowedOrigins, ",")
    } else {
        config.AllowOrigins = []string{"http://localhost:4200"}
    }
    
    config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
    config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
    config.AllowCredentials = true
    
    r.Use(cors.New(config))
    // ^ Apply CORS middleware globally
    
    // Serve uploaded files statically
    r.Static("/uploads", "./uploads")
    // ^ GET /uploads/image.jpg serves ./uploads/image.jpg
    
    // API routes group
    api := r.Group("/api")
    {
        // Auth routes (public)
        auth := api.Group("/auth")
        {
            auth.POST("/register", handlers.Register)
            auth.POST("/login", handlers.Login)
            auth.GET("/me", middleware.AuthMiddleware(), handlers.GetMe)
            // ^ /me requires authentication
        }
        
        // Public listing routes
        api.GET("/listings", handlers.GetListings)
        api.GET("/listings/:id", handlers.GetListing)
        api.GET("/categories", handlers.GetCategories)
        
        // Protected listing routes
        listings := api.Group("/listings")
        listings.Use(middleware.AuthMiddleware())
        // ^ All routes in this group require auth
        {
            listings.POST("", handlers.CreateListing)
            listings.PUT("/:id", handlers.UpdateListing)
            listings.DELETE("/:id", handlers.DeleteListing)
        }
        
        // Chat routes (all protected)
        chats := api.Group("/chats")
        chats.Use(middleware.AuthMiddleware())
        {
            chats.GET("", handlers.GetChats)
            chats.POST("", handlers.CreateChat)
            chats.GET("/:id/messages", handlers.GetMessages)
            chats.POST("/:id/messages", handlers.SendMessage)
        }
        
        // User routes
        users := api.Group("/users")
        {
            users.GET("/:id", handlers.GetUser)
            users.GET("/:id/listings", handlers.GetUserListings)
            
            // Protected routes
            users.PUT("/me", middleware.AuthMiddleware(), handlers.UpdateProfile)
            users.PUT("/me/password", middleware.AuthMiddleware(), handlers.ChangePassword)
        }
        
        // Notification routes (protected)
        notifications := api.Group("/notifications")
        notifications.Use(middleware.AuthMiddleware())
        {
            notifications.GET("", handlers.GetNotifications)
            notifications.PUT("/:id/read", handlers.MarkNotificationRead)
            notifications.PUT("/read-all", handlers.MarkAllNotificationsRead)
        }
        
        // Upload route (protected)
        api.POST("/upload", middleware.AuthMiddleware(), handlers.UploadImage)
    }
    
    // Get port from environment or default to 8080
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    r.Run(":" + port)
    // ^ Starts HTTP server
    // Blocks until server stops
}
```

---

# 📋 PART 4: FRONTEND SETUP

## Step 13: Create Angular Project

**Script:**
"Now let's build the frontend. Make sure you have Node.js and Angular CLI installed."

```bash
# Navigate to project root
cd uf-marketplace

# Create Angular project
ng new frontend --routing --style=scss --standalone
# ^ --routing: Add router
# ^ --style=scss: Use SCSS for styling
# ^ --standalone: Use standalone components (Angular 17+)

cd frontend
```

## Step 14: Install Dependencies

```bash
npm install @angular/material @angular/cdk
# ^ Material UI components (optional but nice)
```

## Step 15: Create Services

**File: `src/app/services/auth.service.ts`**

```typescript
import { Injectable, inject, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { environment } from '../../environments/environment';

// Interface definitions
export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  profile_image?: string;
  phone?: string;
  bio?: string;
  is_admin: boolean;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
}

@Injectable({
  providedIn: 'root'
  // ^ Singleton - same instance everywhere
  // Alternative: provide in component for per-component instance
})
export class AuthService {
  private http = inject(HttpClient);
  // ^ Inject HttpClient using inject() function
  // Alternative: constructor(private http: HttpClient) {}
  
  private apiUrl = environment.apiUrl;
  private tokenKey = 'auth_token';
  
  // Signal for reactive state
  currentUser = signal<User | null>(null);
  // ^ signal() creates a reactive primitive
  // When value changes, Angular updates views automatically
  // Alternative: BehaviorSubject from RxJS
  
  // Computed property based on signal
  isLoggedIn = computed(() => this.currentUser() !== null);
  // ^ Automatically recomputes when currentUser changes
  
  constructor() {
    this.loadUserFromToken();
  }
  
  register(data: RegisterRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(
      `${this.apiUrl}/auth/register`, 
      data
    ).pipe(
      tap(response => {
        // tap() performs side effect without modifying data
        this.storeAuthData(response);
      })
    );
  }
  
  login(email: string, password: string): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(
      `${this.apiUrl}/auth/login`,
      { email, password }
    ).pipe(
      tap(response => this.storeAuthData(response))
    );
  }
  
  logout(): void {
    localStorage.removeItem(this.tokenKey);
    this.currentUser.set(null);
    // ^ .set() updates the signal value
  }
  
  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }
  
  private storeAuthData(response: AuthResponse): void {
    localStorage.setItem(this.tokenKey, response.token);
    this.currentUser.set(response.user);
  }
  
  private loadUserFromToken(): void {
    const token = this.getToken();
    if (token) {
      // Validate token by calling /auth/me
      this.http.get<User>(`${this.apiUrl}/auth/me`).subscribe({
        next: user => this.currentUser.set(user),
        error: () => {
          // Token invalid/expired
          this.logout();
        }
      });
    }
  }
}
```

**Why Signals instead of BehaviorSubject?**
```typescript
// Old way with RxJS BehaviorSubject:
private currentUserSubject = new BehaviorSubject<User | null>(null);
currentUser$ = this.currentUserSubject.asObservable();

// Subscribe in component:
this.authService.currentUser$.subscribe(user => {
  this.user = user;
});

// Signal way (Angular 17+):
currentUser = signal<User | null>(null);

// Use in component:
user = this.authService.currentUser; // Direct access
// In template: {{ authService.currentUser()?.name }}
```

*Signals are simpler, more performant, and don't require manual unsubscription.*

---

## Step 16: HTTP Interceptor

**File: `src/app/interceptors/auth.interceptor.ts`**

```typescript
import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { AuthService } from '../services/auth.service';

// Functional interceptor (Angular 17+)
export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);
  const token = authService.getToken();
  
  // Skip if no token
  if (!token) {
    return next(req);
  }
  
  // Clone request and add header
  const authReq = req.clone({
    setHeaders: {
      Authorization: `Bearer ${token}`
    }
  });
  // ^ Requests are immutable - must clone to modify
  
  return next(authReq);
};

// Old class-based way (pre-Angular 17):
/*
@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  constructor(private authService: AuthService) {}
  
  intercept(req: HttpRequest<any>, next: HttpHandler) {
    // Same logic...
  }
}
*/
```

**Register in app.config.ts:**
```typescript
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { authInterceptor } from './interceptors/auth.interceptor';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(withInterceptors([authInterceptor])),
    // ^ Register interceptor here
  ]
};
```

---

## Step 17: Environment Configuration

**File: `src/environments/environment.ts`**

```typescript
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080/api'
};
```

**File: `src/environments/environment.prod.ts`**

```typescript
export const environment = {
  production: true,
  apiUrl: 'https://your-railway-app.up.railway.app/api'
  // ^ Update this with your actual production URL
};
```

---

## Step 18: Create Components

**Script:**
"Let's generate our components using Angular CLI."

```bash
# Generate pages
ng g c pages/login --standalone
ng g c pages/register --standalone
ng g c pages/dashboard --standalone
ng g c pages/listing-detail --standalone
ng g c pages/create-listing --standalone
ng g c pages/profile --standalone
ng g c pages/settings --standalone

# Generate reusable components
ng g c components/navbar --standalone
ng g c components/listing-card --standalone
ng g c components/chat-widget --standalone
ng g c components/notification-dropdown --standalone

# Generate guard
ng g guard guards/auth --functional
```

---

# 📋 PART 5: BUILDING FRONTEND COMPONENTS

## Step 19: Login Component

**File: `src/app/pages/login/login.component.ts`**

```typescript
import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule, ActivatedRoute } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  // ^ CommonModule: *ngIf, *ngFor, pipes
  // ^ FormsModule: [(ngModel)] two-way binding
  // ^ RouterModule: routerLink directive
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss'
})
export class LoginComponent {
  // Inject dependencies using inject()
  private authService = inject(AuthService);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  
  // Form state
  email = '';
  password = '';
  error = '';
  isLoading = false;
  
  login(): void {
    // Validate
    if (!this.email || !this.password) {
      this.error = 'Please fill in all fields';
      return;
    }
    
    this.isLoading = true;
    this.error = '';
    
    this.authService.login(this.email, this.password).subscribe({
      next: () => {
        // Get return URL from query params or default to dashboard
        const returnUrl = this.route.snapshot.queryParams['returnUrl'] || '/';
        this.router.navigateByUrl(returnUrl);
      },
      error: (err) => {
        this.isLoading = false;
        this.error = err.error?.error || 'Login failed';
        // ^ Extract error message from backend response
      }
    });
  }
}
```

**File: `src/app/pages/login/login.component.html`**

```html
<div class="login-page">
  <div class="login-card">
    <div class="logo">
      <img src="assets/gator-logo.png" alt="UF Logo">
      <h1>UF Market</h1>
    </div>
    
    <h2>Sign in to your account</h2>
    
    <!-- Error message -->
    @if (error) {
      <div class="error-message">{{ error }}</div>
    }
    <!-- ^ @if is Angular 17's new control flow
         Old way: <div *ngIf="error">{{ error }}</div> -->
    
    <form (ngSubmit)="login()">
      <!-- (ngSubmit) handles form submission -->
      
      <div class="form-group">
        <label for="email">UF Email</label>
        <input 
          type="email" 
          id="email"
          [(ngModel)]="email"
          name="email"
          placeholder="gator@ufl.edu"
          required
        >
        <!-- [(ngModel)] = two-way binding
             Changes to input update 'email' variable
             Changes to 'email' variable update input -->
      </div>
      
      <div class="form-group">
        <label for="password">Password</label>
        <input 
          type="password" 
          id="password"
          [(ngModel)]="password"
          name="password"
          placeholder="••••••••"
          required
        >
      </div>
      
      <button type="submit" [disabled]="isLoading">
        <!-- [disabled] = property binding
             When isLoading is true, button is disabled -->
        @if (isLoading) {
          Signing in...
        } @else {
          Sign In
        }
      </button>
    </form>
    
    <p class="register-link">
      Don't have an account? 
      <a routerLink="/register">Create one</a>
      <!-- routerLink for client-side navigation -->
    </p>
  </div>
</div>
```

---

## Step 20: Dashboard Component

**File: `src/app/pages/dashboard/dashboard.component.ts`**

```typescript
import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { NavbarComponent } from '../../components/navbar/navbar.component';
import { ListingCardComponent } from '../../components/listing-card/listing-card.component';
import { ListingService, Listing, Category } from '../../services/listing.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule, 
    FormsModule, 
    RouterModule, 
    NavbarComponent, 
    ListingCardComponent
  ],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss'
})
export class DashboardComponent implements OnInit {
  private listingService = inject(ListingService);
  
  // State using signals
  listings = signal<Listing[]>([]);
  categories = signal<Category[]>([]);
  isLoading = signal(false);
  
  // Filter state (regular properties - not reactive)
  searchQuery = '';
  selectedCategory: number | null = null;
  minPrice: number | null = null;
  maxPrice: number | null = null;
  
  // Pagination
  currentPage = 1;
  totalPages = 1;
  
  ngOnInit(): void {
    // Called once when component initializes
    this.loadCategories();
    this.loadListings();
  }
  
  loadCategories(): void {
    this.listingService.getCategories().subscribe(cats => {
      this.categories.set(cats);
    });
  }
  
  loadListings(): void {
    this.isLoading.set(true);
    
    // Build filters object
    const filters: any = {
      page: this.currentPage,
      limit: 20
    };
    
    if (this.searchQuery) filters.search = this.searchQuery;
    if (this.selectedCategory) filters.category_id = this.selectedCategory;
    if (this.minPrice) filters.min_price = this.minPrice;
    if (this.maxPrice) filters.max_price = this.maxPrice;
    
    this.listingService.getListings(filters).subscribe({
      next: (response) => {
        this.listings.set(response.listings || []);
        this.totalPages = response.pages;
        this.isLoading.set(false);
      },
      error: (err) => {
        console.error('Error loading listings:', err);
        this.isLoading.set(false);
      }
    });
  }
  
  onSearch(): void {
    this.currentPage = 1; // Reset to first page
    this.loadListings();
  }
  
  selectCategory(categoryId: number | null): void {
    this.selectedCategory = categoryId;
    this.currentPage = 1;
    this.loadListings();
  }
  
  previousPage(): void {
    if (this.currentPage > 1) {
      this.currentPage--;
      this.loadListings();
    }
  }
  
  nextPage(): void {
    if (this.currentPage < this.totalPages) {
      this.currentPage++;
      this.loadListings();
    }
  }
}
```

**File: `src/app/pages/dashboard/dashboard.component.html`**

```html
<app-navbar></app-navbar>

<main class="dashboard">
  <!-- Hero Section -->
  <section class="hero">
    <h1>Find What You Need on Campus</h1>
    <p>Buy and sell with fellow Gators</p>
    
    <!-- Search Bar -->
    <div class="search-bar">
      <input 
        type="text" 
        [(ngModel)]="searchQuery"
        placeholder="Search for textbooks, electronics, furniture..."
        (keyup.enter)="onSearch()"
      >
      <!-- (keyup.enter) triggers on Enter key -->
      <button (click)="onSearch()">Search</button>
    </div>
  </section>
  
  <!-- Categories -->
  <section class="categories">
    <h2>Browse Categories</h2>
    <div class="category-grid">
      <button 
        class="category-btn"
        [class.active]="selectedCategory === null"
        (click)="selectCategory(null)"
      >
        <!-- [class.active] adds 'active' class when condition is true -->
        All
      </button>
      
      @for (category of categories(); track category.id) {
        <button 
          class="category-btn"
          [class.active]="selectedCategory === category.id"
          (click)="selectCategory(category.id)"
        >
          <span class="material-icons">{{ category.icon }}</span>
          {{ category.name }}
        </button>
      }
      <!-- @for is Angular 17's new control flow
           track is required for optimization (like React key)
           Old way: *ngFor="let category of categories; trackBy: trackById" -->
    </div>
  </section>
  
  <!-- Listings Grid -->
  <section class="listings">
    <h2>Latest Listings</h2>
    
    @if (isLoading()) {
      <div class="loading">Loading...</div>
    } @else if (listings().length === 0) {
      <div class="empty">No listings found</div>
    } @else {
      <div class="listings-grid">
        @for (listing of listings(); track listing.id) {
          <app-listing-card [listing]="listing"></app-listing-card>
          <!-- [listing] passes data to child component -->
        }
      </div>
      
      <!-- Pagination -->
      <div class="pagination">
        <button 
          (click)="previousPage()" 
          [disabled]="currentPage === 1"
        >
          Previous
        </button>
        <span>Page {{ currentPage }} of {{ totalPages }}</span>
        <button 
          (click)="nextPage()" 
          [disabled]="currentPage === totalPages"
        >
          Next
        </button>
      </div>
    }
  </section>
</main>
```

---

## Step 21: Listing Card Component

**File: `src/app/components/listing-card/listing-card.component.ts`**

```typescript
import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { Listing } from '../../services/listing.service';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-listing-card',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './listing-card.component.html',
  styleUrl: './listing-card.component.scss'
})
export class ListingCardComponent {
  @Input({ required: true }) listing!: Listing;
  // ^ @Input receives data from parent component
  // { required: true } makes it mandatory
  // ! asserts it will be defined (non-null assertion)
  
  // Computed property (getter)
  get primaryImage(): string {
    // Find primary image or first image
    const primary = this.listing.images?.find(img => img.is_primary);
    if (primary) return this.getImageUrl(primary.image_url);
    
    if (this.listing.images?.length > 0) {
      return this.getImageUrl(this.listing.images[0].image_url);
    }
    
    // Return placeholder if no images
    return 'assets/placeholder.svg';
  }
  
  get sellerInitials(): string {
    const first = this.listing.seller?.first_name?.charAt(0) || '';
    const last = this.listing.seller?.last_name?.charAt(0) || '';
    return (first + last).toUpperCase();
  }
  
  private getImageUrl(url: string): string {
    // If it's a relative path from our backend, prefix with API URL
    if (url.startsWith('/uploads/')) {
      return environment.apiUrl.replace('/api', '') + url;
    }
    return url;
  }
}
```

**File: `src/app/components/listing-card/listing-card.component.html`**

```html
<a [routerLink]="['/listing', listing.id]" class="listing-card">
  <!-- routerLink with array syntax for dynamic routes -->
  
  <div class="card-image">
    <img [src]="primaryImage" [alt]="listing.title" loading="lazy">
    <!-- loading="lazy" improves performance by lazy loading images -->
    
    @if (listing.status === 'sold') {
      <div class="sold-badge">SOLD</div>
    }
    
    @if (listing.condition) {
      <span class="condition-badge">{{ listing.condition | titlecase }}</span>
      <!-- titlecase pipe: "like_new" -> "Like New" -->
    }
  </div>
  
  <div class="card-content">
    <h3 class="card-title">{{ listing.title }}</h3>
    <span class="card-price">
      {{ listing.price | currency:'USD':'symbol':'1.0-0' }}
      <!-- currency pipe formats as $1,234 -->
    </span>
    
    <div class="card-meta">
      <span class="category">{{ listing.category?.name }}</span>
      <span class="views">{{ listing.views }} views</span>
    </div>
    
    <div class="card-footer">
      <div class="seller-info">
        @if (listing.seller?.profile_image) {
          <img [src]="listing.seller.profile_image" class="seller-avatar">
        } @else {
          <div class="seller-avatar-placeholder">{{ sellerInitials }}</div>
        }
        <span>{{ listing.seller?.first_name }}</span>
      </div>
      <span class="time">{{ listing.created_at | date:'shortDate' }}</span>
      <!-- date pipe formats dates -->
    </div>
  </div>
</a>
```

---

## Step 22: Routing Configuration

**File: `src/app/app.routes.ts`**

```typescript
import { Routes } from '@angular/router';
import { authGuard } from './guards/auth.guard';

export const routes: Routes = [
  // Public routes
  {
    path: '',
    loadComponent: () => import('./pages/dashboard/dashboard.component')
      .then(m => m.DashboardComponent)
    // ^ Lazy loading - component only loaded when route is accessed
    // Reduces initial bundle size
  },
  {
    path: 'login',
    loadComponent: () => import('./pages/login/login.component')
      .then(m => m.LoginComponent)
  },
  {
    path: 'register',
    loadComponent: () => import('./pages/register/register.component')
      .then(m => m.RegisterComponent)
  },
  {
    path: 'listing/:id',
    loadComponent: () => import('./pages/listing-detail/listing-detail.component')
      .then(m => m.ListingDetailComponent)
    // :id is a route parameter
  },
  
  // Protected routes
  {
    path: 'create-listing',
    loadComponent: () => import('./pages/create-listing/create-listing.component')
      .then(m => m.CreateListingComponent),
    canActivate: [authGuard]
    // ^ Route guard - redirects to login if not authenticated
  },
  {
    path: 'profile',
    loadComponent: () => import('./pages/profile/profile.component')
      .then(m => m.ProfileComponent),
    canActivate: [authGuard]
  },
  {
    path: 'settings',
    loadComponent: () => import('./pages/settings/settings.component')
      .then(m => m.SettingsComponent),
    canActivate: [authGuard]
  },
  
  // Wildcard - must be last
  {
    path: '**',
    redirectTo: ''
    // Redirect unknown paths to home
  }
];
```

---

## Step 23: Auth Guard

**File: `src/app/guards/auth.guard.ts`**

```typescript
import { inject } from '@angular/core';
import { Router, CanActivateFn } from '@angular/router';
import { AuthService } from '../services/auth.service';

export const authGuard: CanActivateFn = (route, state) => {
  // CanActivateFn is the type for functional guards
  // route: current route info
  // state: router state (includes full URL)
  
  const authService = inject(AuthService);
  const router = inject(Router);
  
  if (authService.isLoggedIn()) {
    return true; // Allow navigation
  }
  
  // Not logged in - redirect to login with return URL
  return router.createUrlTree(['/login'], {
    queryParams: { returnUrl: state.url }
  });
  // ^ After login, user will be redirected back to original URL
};
```

---

# 📋 PART 6: INTEGRATION & FINAL TESTING

## Step 24: Complete Testing Checklist

### Backend Tests (using Postman or curl)

| Test | Endpoint | Expected |
|------|----------|----------|
| Register with valid data | POST /auth/register | 201 + token |
| Register duplicate email | POST /auth/register | 409 error |
| Register invalid email | POST /auth/register | 400 error |
| Login valid | POST /auth/login | 200 + token |
| Login wrong password | POST /auth/login | 401 error |
| Get me without token | GET /auth/me | 401 error |
| Get me with token | GET /auth/me | 200 + user |
| Create listing | POST /listings | 201 + listing |
| Get all listings | GET /listings | 200 + array |
| Search listings | GET /listings?search=x | 200 + filtered |
| Get single listing | GET /listings/:id | 200 + listing |
| Update own listing | PUT /listings/:id | 200 |
| Update other's listing | PUT /listings/:id | 403 |
| Delete own listing | DELETE /listings/:id | 200 |
| Start chat | POST /chats | 201 |
| Send message | POST /chats/:id/messages | 201 |

### Frontend Tests

| Feature | Test |
|---------|------|
| Registration | Fill form, submit, verify redirect |
| Login | Enter credentials, verify dashboard |
| Logout | Click logout, verify redirect to login |
| Search | Enter query, verify filtered results |
| Category filter | Click category, verify listings |
| Create listing | Fill form, submit, verify appears |
| View listing | Click card, verify detail page |
| Start chat | Click "Message Seller", verify chat opens |
| Send message | Type message, verify appears |

---

# 🎯 OPTIMIZATION TIPS

## Performance Optimizations

### Backend
1. **Add indexes** to frequently searched columns:
```go
// In model definition
Email string `gorm:"uniqueIndex;not null"`
```

2. **Use pagination** instead of returning all records

3. **Cache categories** since they don't change:
```go
var categoriesCache []Category
var categoriesCacheTime time.Time

func GetCategories(c *gin.Context) {
    if time.Since(categoriesCacheTime) < 5*time.Minute {
        c.JSON(200, categoriesCache)
        return
    }
    // Fetch from DB and update cache
}
```

### Frontend
1. **Lazy load routes** (already implemented)

2. **Use OnPush change detection** for components:
```typescript
@Component({
  changeDetection: ChangeDetectionStrategy.OnPush
})
```

3. **Use trackBy** in loops (already implemented with @for track)

4. **Virtual scrolling** for long lists:
```html
<cdk-virtual-scroll-viewport itemSize="100">
  <app-listing-card *cdkVirtualFor="let listing of listings"></app-listing-card>
</cdk-virtual-scroll-viewport>
```

---

# 📚 ADDITIONAL RESOURCES

## What We Didn't Cover (for future videos)

1. **Image uploads** - File handling with multipart forms
2. **Real-time chat** - WebSockets with Gorilla WebSocket
3. **Email verification** - SMTP integration
4. **Password reset** - Token-based reset flow
5. **Admin panel** - Role-based access control
6. **Deployment** - Docker, Railway, Vercel
7. **Testing** - Unit tests, integration tests
8. **CI/CD** - GitHub Actions

## Recommended Learning Path

1. **Go**: Tour of Go (tour.golang.org)
2. **Gin**: Official docs (gin-gonic.com)
3. **Angular**: angular.dev tutorials
4. **REST APIs**: RESTful Web Services by Leonard Richardson
5. **Database Design**: SQL and Relational Theory by C.J. Date

---

**End of Tutorial Guide**

*Total estimated video time: 4-5 hours split into 6 parts*
