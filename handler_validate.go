package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
    type returnVals struct {
        CleanedBody string `json:"cleaned_body"`
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

    respondWithJSON(w, http.StatusOK, returnVals{
        CleanedBody: cleanString(params.Body, " ", []string{"kerfuffle", "sharbert", "fornax"}, "****"),
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
