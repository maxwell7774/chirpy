package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
    dbChirps, err := cfg.db.GetChirps(r.Context())
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
        return
    }

    chirps := make([]Chirp, len(dbChirps))
    for i, c := range dbChirps {
        chirps[i] = Chirp{
            ID: c.ID,
            CreatedAt: c.CreatedAt,
            UpdatedAt: c.UpdatedAt,
            Body: c.Body,
            UserID: c.UserID,
        }
    }

    respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {
    id, err := uuid.Parse(r.PathValue("chirpID"))
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Please provide a valid id", err)
        return
    }

    chirp, err := cfg.db.GetChirp(r.Context(), id)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirp with that id", err)
        return
    }

    respondWithJSON(w, http.StatusOK, Chirp{
        ID: chirp.ID,
        CreatedAt: chirp.CreatedAt,
        UpdatedAt: chirp.UpdatedAt,
        Body: chirp.Body,
        UserID: chirp.UserID,
    })
}
