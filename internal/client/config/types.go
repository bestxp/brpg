package config

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type Resolution string

const resSplitter = "x"

func (r Resolution) Width() int {
	parts := strings.Split(string(r), resSplitter)
	val, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msgf("parsing width: %s", r)
	}
	return int(val)
}

func (r Resolution) Height() int {
	parts := strings.Split(string(r), resSplitter)
	val, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msgf("parsing height: %s", r)
	}
	return int(val)
}
