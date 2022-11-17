package main

import (
	"librarymanagement/controller"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllBooks(t *testing.T) {
	req, err := http.NewRequest("GET","/books/{criteria}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.GetBookByNameOrAuthorOrGenre)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
