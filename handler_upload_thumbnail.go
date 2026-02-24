package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here

	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)

	fileData, fileHeader, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't return file for thumbnail", err)
		return
	}

	mediaType, _, err := mime.ParseMediaType(fileHeader.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid content-type", err)
		return
	}
	if mediaType != "image/png" && mediaType != "image/jpeg" {
		respondWithError(w, http.StatusBadRequest, "not right mediatype", nil)
		return
	}

	videoDb, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "didnt find the video", err)
		return
	}
	if videoDb.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "not the owner of the video", err)
	}

	key := make([]byte, 32)
	rand.Read(key)
	urlKey := base64.RawURLEncoding.EncodeToString(key)

	assetPath := getAssetPath(urlKey, mediaType)

	fileDst, err := os.Create(cfg.getAssetDiskPath(assetPath))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt create file for thumbnail", err)
		return
	}
	defer fileDst.Close()
	defer fileData.Close()

	_, err = io.Copy(fileDst, fileData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt write data to the thumbnail file", err)
		return
	}

	url := cfg.getAssetUrl(assetPath)
	videoDb.ThumbnailURL = &url
	cfg.db.UpdateVideo(videoDb)

	respondWithJSON(w, http.StatusOK, videoDb)
}
