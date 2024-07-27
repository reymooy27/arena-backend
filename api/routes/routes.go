package routes

import (
	"fmt"
	"net/http"

	"github.com/reymooy27/arena-backend/api/controllers"
)

func Routes() {
	http.HandleFunc("/", controllers.GetHello)
	http.HandleFunc("/post", controllers.PostData)

	fmt.Println("Server running on port 8000")
	http.ListenAndServe(":8000", nil)

}
