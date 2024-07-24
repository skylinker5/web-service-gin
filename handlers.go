package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

/*
 * store user uploaded file at server local
 */
func uploadSingleFile(c *gin.Context) {
	// Single file
	file, err := c.FormFile("file")
	log.Println(file.Filename)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, "./uploaded_files/"+file.Filename)

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

/*
 * generate jwt token and send back to frontend
 */
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func generateJWTHandler(w http.ResponseWriter, r *http.Request) {
	// Try to get the JWT from the cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If no cookie is found, generate a new UUID
			userID := uuid.New().String()
			expirationTime := time.Now().Add(2 * time.Hour)

			claims := &Claims{
				UserID: userID,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(expirationTime),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Set the JWT in the cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
				Path:    "/",
			})

			fmt.Fprintf(w, "Generated and set new token for User ID: %s", userID)
			return
		}

		// For any other error, return a bad request status
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// If the cookie is found, validate the JWT
	tokenString := cookie.Value
	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !parsedToken.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "User ID from token: %s", claims.UserID)

	// Validate and extract claims
	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		fmt.Println("Claims:", claims)
	} else {
		fmt.Println("Invalid token")
	}
}
