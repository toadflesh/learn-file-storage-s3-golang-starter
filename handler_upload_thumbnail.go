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

	// #2
	const maxMemory = 10 << 20
	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "multipart form error", err)
		return
	}

	// #3
	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "form file error", err)
		return
	}

	defer file.Close()
	mediaType := header.Header.Get("Content-Type")

	// #4
	imageData, err := io.ReadAll(file)
	if err != nil {
		return
	}

	// #5
	metadata, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "userid for video doesnt match", err)
		return
	}

	// #6
	var tn = thumbnail{imageData, mediaType}
	videoThumbnails[videoID] = tn

	// #7
	thumbnailURL := fmt.Sprintf("http://localhost:8091/api/thumbnails/%s", videoID)
	metadata.ThumbnailURL = &thumbnailURL

	err = cfg.db.UpdateVideo(metadata)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update video with thumbnail URL", err)
		return
	}

	respondWithJSON(w, http.StatusOK, metadata)
}
