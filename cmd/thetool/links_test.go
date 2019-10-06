package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckLink(t *testing.T) {
	fixtures := map[string]bool{
		"https://fa.wikipedia.org/wiki/صفحهٔ_اصلی":                                    true,
		"https://fa.wikipedia.org/wiki/سیقفبغلعاتهنخمثسقیفبغلعاتهنمک۴۵۶۷غعهخحیبلاذتد": false,
	}

	client := &http.Client{
		Timeout: time.Second,
	}

	for lnk, valid := range fixtures {
		err := checkSingle(client, lnk)
		if valid {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
