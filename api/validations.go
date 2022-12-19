package api

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
)

func validateRequestParameters(sortKey, limit string) []BadRequestError {
	errors := []BadRequestError{}
	if sortKey != sortKeyRelevanceScore && sortKey != sortKeyViews {
		errors = append(errors, createBadRequestError("sortKey", fmt.Sprintf(
			"sortKey is invalid. It should be either %s or %s", sortKeyRelevanceScore, sortKeyViews)))
	}

	limitInt, parseIntErr := strconv.Atoi(limit)
	if parseIntErr != nil {
		log.Warn().Err(parseIntErr).Msg("error parsing limit string to int")
	}

	if limitInt < 2 || limitInt > 199 {
		errors = append(errors, createBadRequestError("limit", "limit should be a number greater than 1 and less than 200"))
	}

	return errors
}
