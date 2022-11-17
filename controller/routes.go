package controller

import "github.com/gorilla/mux"

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/books",LibraryController{}.GetAllBooks).Methods("GET")
	router.HandleFunc("/books/{criteria}",GetBookByNameOrAuthorOrGenre).Methods("GET")
	router.HandleFunc("/books/added/{name}",GetBooksAddedByUser).Methods("GET")

	router.HandleFunc("/books",AddBook).Methods("POST")

	router.HandleFunc("/books/borrow/{id}/{name}",BorrowBook).Methods("GET")
	router.HandleFunc("/books/borrowed/{name}",GetBooksBorrowedByUser).Methods("GET")

	router.HandleFunc("/books/rating/{name}/{id}/{rating}",RateBook).Methods("GET")
	router.HandleFunc("/user/rating/{bname}/{lname}/{rating}",RateUser).Methods("GET")

	router.HandleFunc("/books/return/{id}/{name}",ReturnBook).Methods("GET")
	
	router.HandleFunc("/books/remove/{name}/{id}",RemoveBook).Methods("DELETE")


}