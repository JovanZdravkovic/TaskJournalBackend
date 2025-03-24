package handlers

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

var (
	Users    = regexp.MustCompile(`^/user/*$`)
	UserIcon = regexp.MustCompile(`^/user/icon/*$`)
)

const iconUploadDirectory = "uploads/profile_icons"
const maxIconSize = 500000

type UserHandler struct {
	DBService *db.DatabaseService
}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	preflight := EnableCORS(w, r)
	if preflight {
		return
	}

	tokenString := r.Header.Get("X-Auth-Token")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := uuid.Parse(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch {
	case r.Method == http.MethodGet && Users.MatchString(r.URL.Path):
		u.GetUser(w, r, token)
		return
	case r.Method == http.MethodPut && Users.MatchString(r.URL.Path):
		u.UpdateUser(w, r, token)
		return
	case r.Method == http.MethodGet && UserIcon.MatchString(r.URL.Path):
		u.GetIcon(w, r, token)
		return
	case r.Method == http.MethodPost && UserIcon.MatchString(r.URL.Path):
		u.UploadIcon(w, r, token)
		return
	default:
		return
	}
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	user, err := u.DBService.GetUserInfo(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var user db.UserPut
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}
	err = u.DBService.UpdateUser(user, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	responseJson, err := json.Marshal(db.Success{Success: true})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func (u *UserHandler) UploadIcon(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	err := r.ParseMultipartForm(maxIconSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Icon can't be larger than 500kb"))
		return
	}

	file, handler, err := r.FormFile("icon")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid file"))
		return
	}
	defer file.Close()

	contentType := handler.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/jpg" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("File is not jpeg format"))
		return
	}

	iconSize := handler.Size
	if iconSize > maxIconSize {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Icon can't be larger than 500kb"))
		return
	}

	if err := os.MkdirAll(iconUploadDirectory, os.ModePerm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create upload directory"))
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while decoding image"))
		return
	}

	resizedImg := imaging.Resize(img, 100, 100, imaging.Lanczos)

	userIdString := userId.String()
	filename := fmt.Sprintf("icon-%s.png", userIdString)
	filePath := filepath.Join(iconUploadDirectory, filename)

	outFile, err := os.Create(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to save file"))
		return
	}
	defer outFile.Close()

	jpegOptions := &jpeg.Options{Quality: 100}
	if err := jpeg.Encode(outFile, resizedImg, jpegOptions); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode image"))
		return
	}

	responseJson, err := json.Marshal(db.Success{Success: true})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)

}

func (u *UserHandler) GetIcon(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	userIdString := userId.String()
	filename := fmt.Sprintf("icon-%s.png", userIdString)
	filePath := filepath.Join(iconUploadDirectory, filename)

	if _, err := os.Stat(filePath); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Profile picture not found"))
		return
	}

	http.ServeFile(w, r, filePath)
}
