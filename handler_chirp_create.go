package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/maxwell7774/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const chirpMaxLength = 140
	if len(params.Body) > chirpMaxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := cleanString(
		params.Body,
		" ",
		[]string{"kerfuffle", "sharbert", "fornax"},
		"****",
	)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:  cleanedBody,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func cleanString(str string, delimeter string, wordsToClean []string, replaceWith string) string {
	words := strings.Split(str, delimeter)

	wordsToCleanMap := make(map[string]bool, len(wordsToClean))
	for _, word := range wordsToClean {
		wordsToCleanMap[word] = true
	}

	for i, word := range words {
		word = strings.ToLower(word)
		if wordsToCleanMap[word] {
			words[i] = replaceWith
		}
	}

	return strings.Join(words, delimeter)
}
