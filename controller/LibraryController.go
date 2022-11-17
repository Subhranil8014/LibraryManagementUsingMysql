package controller

import (
	"encoding/json"
	"fmt"
	"librarymanagement/database"
	"librarymanagement/entities"
	"librarymanagement/interfaces"
	"net/http"
	"strconv"
	"strings"
	"github.com/gorilla/mux"
)
 type LibraryController struct{
	interfaces.ILibraryManager
 }

func AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book entities.BookList
	json.NewDecoder(r.Body).Decode(&book)
	if getUserRating(book.LenderName)!=0 {
		book.LenderRating=getUserRating(book.LenderName)
	}else{
		book.LenderRating=0
	}
	book.Rating=0
	database.DbInstance.Create(&book)
	json.NewEncoder(w).Encode(book)
}

func getUserRating(name string) (float32){
	var book entities.BookList
	database.DbInstance.Where("lender_name=?",name).First(&book)
	if strings.EqualFold(book.LenderName,name) {
		return book.LenderRating
	}
	return 0
}

func (controller LibraryController) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	var books []entities.BookList
	database.DbInstance.Find(&books)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)

}

func GetBookByNameOrAuthorOrGenre(w http.ResponseWriter, r *http.Request) {
	criteria := mux.Vars(r)["criteria"]
	exr:="doesn't exist"
	var book entities.BookList
	var books []entities.BookList
	database.DbInstance.Where("book_name = ? OR author = ? OR genre=?", criteria ,criteria ,criteria).First(&book)
	if strings.EqualFold(book.BookName,criteria) {
		w.Header().Set("Content-Type", "application/json")
	    json.NewEncoder(w).Encode(book)
	} else if strings.EqualFold(book.Author,criteria){
		database.DbInstance.Where("author = ?",criteria).Find(&books)
		w.Header().Set("Content-Type", "application/json")
	    json.NewEncoder(w).Encode(books)
	}else if strings.EqualFold(book.Genre,criteria){
		database.DbInstance.Where("genre = ?",criteria).Find(&books)
		w.Header().Set("Content-Type", "application/json")
	    json.NewEncoder(w).Encode(books)
	}else {
	w.Write([]byte(exr))
	}
}

func GetBooksAddedByUser(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	exr:="doesn't exist"
	var books []entities.BookList
	database.DbInstance.Where("lender_name = ?", name ).Find(&books)
	if len(books)==0 {
		json.NewEncoder(w).Encode(exr)
		return
	}else{
	    w.Header().Set("Content-Type", "application/json")
	    w.WriteHeader(http.StatusOK)
	    json.NewEncoder(w).Encode(books)
	}
	
}

func BorrowBook(w http.ResponseWriter, r *http.Request) {
    bookId := mux.Vars(r)["id"]
	borrowername:=mux.Vars(r)["name"]
	exr:="doesn't exist"
	var borrowRecord entities.BorrowRecord
	var borrowToken entities.BorrowUpdate
	var lendToken entities.LendingRecord
	var book entities.BookList
	database.DbInstance.First(&book,"id=?", bookId)
	
	if book.ID==bookId{
		if strings.EqualFold(book.Status,"no") {
			json.NewEncoder(w).Encode("Book is already borrowed")
			return
		}
		database.DbInstance.Model(&book).Where("id = ?", bookId).Update("status","no")
        w.Header().Set("Content-Type", "application/json")
		borrowRecord.BorrowerName=borrowername
		borrowRecord.BookID=bookId
		borrowRecord.BookName=book.BookName
		borrowRecord.Author=book.Author
		borrowRecord.Genre=book.Genre
		borrowRecord.Rating=book.Rating
		borrowRecord.Status=book.Status
		database.DbInstance.Save(&borrowRecord)

        database.DbInstance.Where("borrower_name=?",borrowername).First(&borrowToken)
		if strings.EqualFold(borrowToken.BorrowerName,borrowername) {
			database.DbInstance.Model(&borrowToken).Where("borrower_name = ?", borrowername).Update("token",borrowToken.Token-1)
		}else{
			borrowToken.BorrowerName=borrowername
			borrowToken.Token=-1
			database.DbInstance.Save(&borrowToken)
		}

		database.DbInstance.Where("lender_name=?",book.LenderName).First(&lendToken)
		if strings.EqualFold(lendToken.LenderName,book.LenderName) {
			database.DbInstance.Model(&lendToken).Where("lender_name=?",book.LenderName).Update("token",lendToken.Token+1)
		}else{
			lendToken.LenderName=book.LenderName
			lendToken.Token=1
			database.DbInstance.Save(&lendToken)
		}
	    json.NewEncoder(w).Encode("successfully book issued")
	}else{
		json.NewEncoder(w).Encode(exr)
	}
}

func GetBooksBorrowedByUser(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	exr:="doesn't exist"
	var books []entities.BorrowRecord
	database.DbInstance.Where("borrower_name = ?", name ).Find(&books)
	if len(books)==0 {
		json.NewEncoder(w).Encode(exr)
		return
	}else{
	    w.Header().Set("Content-Type", "application/json")
	    w.WriteHeader(http.StatusOK)
	    json.NewEncoder(w).Encode(books)
	}
}

func RateBook(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["id"]
	name := mux.Vars(r)["name"]
	rating:=mux.Vars(r)["rating"]
	exr:="doesn't exist"
	var book entities.BookList
	var borrowRecords []entities.BorrowRecord
	database.DbInstance.First(&book,"id=?", bookId)
	database.DbInstance.Where("book_id=?", bookId).Find(&borrowRecords)

	if book.ID==bookId{
		value, _ := strconv.ParseFloat(rating, 32)
		if strings.EqualFold(book.LenderName,name) {
			json.NewEncoder(w).Encode("You can't rate your book")
			return
		}
		if float32(value)>10{
			json.NewEncoder(w).Encode("Rating should be out of 10")
			return
		}
		for i := 0; i < len(borrowRecords); i++ {
		if strings.EqualFold(borrowRecords[i].BorrowerName,name){
			if borrowRecords[i].Rating==0 {
				database.DbInstance.Model(&borrowRecords).Where("book_id=? AND borrower_name = ?",book.ID,borrowRecords[i].BorrowerName).Update("rating",value)
				database.DbInstance.Model(&book).Where("id=?",book.ID).Update("rating",value)
				json.NewEncoder(w).Encode("your have rated the book successfully")
			    return
			}else{
				newRating:=((borrowRecords[i].Rating+float32(value))/2)
				database.DbInstance.Model(&borrowRecords).Where("book_id=? AND borrower_name = ?",book.ID,borrowRecords[i].BorrowerName).Update("rating",newRating)
				database.DbInstance.Model(&book).Where("id=?",book.ID).Update("rating",newRating)
				json.NewEncoder(w).Encode("your have rated the book successfully")
			    return
			}	
		}
	}
	json.NewEncoder(w).Encode("Can't rate a book you have not borrowed")
	return
	}else{
		json.NewEncoder(w).Encode(exr)
	}
}

func RateUser(w http.ResponseWriter, r *http.Request) {
	bname := mux.Vars(r)["bname"]
	lname := mux.Vars(r)["lname"]
	rating:=mux.Vars(r)["rating"]
	value, _ := strconv.ParseFloat(rating, 32)
	exr:="doesn't exist"
	var book entities.BookList
	var borrowRecords []entities.BorrowRecord
	database.DbInstance.Where("borrower_name=?", bname).Find(&borrowRecords)
	//database.DbInstance.Find(&books)

    if len(borrowRecords)!=0 {
	   for i := 0; i < len(borrowRecords); i++ {
	//	database.DbInstance.Where("id=? ",borrowRecords[i].BookID).Find(&book)
	//	database.DbInstance.First(&book,"id <> ?",borrowRecords[i].BookID)
		database.DbInstance.Raw("SELECT * FROM `BookList` WHERE id=?",borrowRecords[i].BookID).Scan(&book)
		if strings.EqualFold(book.LenderName,lname){
			if float32(value)>10{
				json.NewEncoder(w).Encode("Rating should be out of 10")
				return
			}
			if book.LenderRating==0 {
				database.DbInstance.Raw("UPDATE BookList SET lender_rating=? WHERE lender_name = ?",value,book.LenderName).Scan(&book)
				json.NewEncoder(w).Encode("your have rated the user successfully")
			    return
			}else{
				newRating:=((book.LenderRating+float32(value))/2)
				database.DbInstance.Raw("UPDATE BookList SET lender_rating=? WHERE lender_name = ?",newRating,book.LenderName).Scan(&book)
				json.NewEncoder(w).Encode("your have rated the user successfully")
			    return
			}

		}
	}
	json.NewEncoder(w).Encode("Can't rate a User,whose book you haven't read yet")
	return
}else{
	json.NewEncoder(w).Encode(exr)
	return
}

}

func ReturnBook(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["id"]
	borrowername:=mux.Vars(r)["name"]
	exr:="doesn't exist"
	var borrowRecord entities.BorrowRecord
	var borrowToken entities.BorrowUpdate
	var lendToken entities.LendingRecord
	var book entities.BookList
	database.DbInstance.First(&book,"id=?", bookId)
	if book.ID==bookId{
		if strings.EqualFold(book.Status,"yes") {
			fmt.Println("enter1")
			json.NewEncoder(w).Encode("Book is not issued")
			return
		}else{
			fmt.Println("enter2")
		database.DbInstance.Model(&book).Where("id = ?", bookId).Update("status","yes")
		database.DbInstance.Model(&borrowRecord).Where("book_id = ?", bookId).Update("status","yes")
		database.DbInstance.Where("book_id = ?", bookId).First(&borrowRecord)
			if strings.EqualFold(borrowRecord.BorrowerName,borrowername) {
				database.DbInstance.Where("borrower_name=?",borrowername).First(&borrowToken)
				database.DbInstance.Where("lender_name=?",book.LenderName).First(&lendToken)
				fmt.Println("enter3",borrowToken.Token,lendToken.Token)
				database.DbInstance.Raw("UPDATE borrowtokenupdate SET token=? WHERE borrower_name = ?",borrowToken.Token+1,borrowername).Scan(&borrowToken)
				database.DbInstance.Raw("UPDATE lendertokenupdate SET token=? WHERE lender_name = ?",lendToken.Token-1,book.LenderName).Scan(&lendToken)
			}else{
				fmt.Println("enter4")
				json.NewEncoder(w).Encode("book wasn't issued by you")
				return
			}
			fmt.Println("enter5")
	    json.NewEncoder(w).Encode("successfully book returned")
		return
	}
	}else{
		fmt.Println("enter6")
		json.NewEncoder(w).Encode(exr)
	}
}

func RemoveBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookId := mux.Vars(r)["id"]
	name:=mux.Vars(r)["name"]
	exr:="doesn't exist"
	var book entities.BookList
    database.DbInstance.First(&book,"id=?", bookId)
	if book.ID==bookId {
		if book.LenderName==name {
			database.DbInstance.Delete(&book,"id=?", bookId)
	        json.NewEncoder(w).Encode("Book Deleted Successfully!")
			return
		}else{
			json.NewEncoder(w).Encode("You can't delete a book which you have not added!")
			return
		}
	}else{
		json.NewEncoder(w).Encode(exr)
		return
	}
	
}
