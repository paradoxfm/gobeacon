package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := createPhoneAdminApi()

	w := httptest.NewRecorder()
	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, _ := http.NewRequest("GET", "/api/v1/test/push", bytes.NewBuffer(jsonStr))
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIzMjYyNTY0MzYsIm9yaWdfaWF0IjoxNTM3ODU2NDM2LCJwcml2YXRlX2NsYWltX2lkIjoiNTRkOTYyM2MtMTFhZi0xMWU4LWEyN2QtNTQ1MjAwN2Q4NzdlIn0.WTK3UWQaja_6nKcuDu5QpxsGaVeIddjRegKwGMqviZ0")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test(t *testing.T) {
	g := Goblin(t)
	g.Describe("Набор тестов", func() {
		// Passing Test
		g.It("ДОлжно быть сложено 2 числа ", func() {
			g.Assert(1 + 1).Equal(2)
		})
		g.It("Еще один тест ", func() {
			g.Assert(1 + 1).Equal(2)
		})
		// Failing Test
		g.It("Should match equal numbers", func() {
			g.Assert(2).Equal(4)
		})
		// Pending Test
		g.It("Should substract two numbers")
		// Excluded Test
		g.Xit("Should add two numbers ", func() {
			g.Assert(3 + 1).Equal(4)
		})
	})
}
