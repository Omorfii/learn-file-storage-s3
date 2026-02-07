package main

import (
	"fmt"
	"io"
	"net/http"

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
	mediaType := fileHeader.Header.Get("Content-Type")

	imageData, err := io.ReadAll(fileData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't return imageData", err)
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

	thumbnail := thumbnail{
		data:      imageData,
		mediaType: mediaType,
	}
	videoThumbnails[videoDb.ID] = thumbnail

	url := "http://localhost:" + cfg.port + "/api/thumbnails/" + videoDb.ID.String()
	videoDb.ThumbnailURL = &url
	cfg.db.UpdateVideo(videoDb)

	respondWithJSON(w, http.StatusOK, videoDb)
}
