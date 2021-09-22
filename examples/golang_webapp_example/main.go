package main
import (
   "net/http"
   "fmt"
   "time"
   "html/template"
)

//Create a struct that holds information to be displayed in our HTML file
type Welcome struct {
   Name string
   Time string
}

//Go application entrypoint
func main() {
   //Instantiate a Welcome struct object and pass in some random information. 
   //We shall get the name of the user as a query parameter from the URL
   welcome := Welcome{"Dharmesh", time.Now().Format(time.Stamp)}

   //We tell Go exactly where we can find our html file. We ask Go to parse the html file (Notice
   // the relative path). We wrap it in a call to template.Must() which handles any errors and halts if there are fatal errors
   
   templates := template.Must(template.ParseFiles("template/welcome-template.html"))

   

   //This method takes in the URL path "/" and a function that takes in a response writer, and a http request.
   http.HandleFunc("/" , func(w http.ResponseWriter, r *http.Request) {

      //Takes the name from the URL query e.g ?name=Martin, will set welcome.Name = Martin.
      if name := r.FormValue("name"); name != "" {
         welcome.Name = name;
      }
      //If errors show an internal server error message
      //I also pass the welcome struct to the welcome-template.html file.
      if err := templates.ExecuteTemplate(w, "welcome-template.html", welcome); err != nil {
         http.Error(w, err.Error(), http.StatusInternalServerError)
      }
   })

   //Start the web server, set the port to listen to 8080. Without a path it assumes localhost
   //Print any errors from starting the webserver using fmt
   fmt.Println("Listening");
   fmt.Println(http.ListenAndServe(":8080", nil));
}