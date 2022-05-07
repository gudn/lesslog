package logging

import (
	"github.com/rs/zerolog/log"
)

func LogRequest(method string, err error) {
	if err != nil {
		log.
			Error().
			Err(err).
			Str("method", method).
			Msg("error processing request")
	} else {
		log.
			Info().
			Str("method", method).
			Msg("success processing request")
	}
}
