package model_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/sonda2208/grc/model"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	es := model.Config{}
	res, err := json.Marshal(es)
	require.NoError(t, err)

	log.Println(string(res))
}
