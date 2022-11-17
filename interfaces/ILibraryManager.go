package interfaces

import(
	"net/http"
)
type ILibraryManager interface{
GetAllBooks(w http.ResponseWriter, r *http.Request) 
}
	
