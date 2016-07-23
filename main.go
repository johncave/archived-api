package main

import (
  "encoding/json"
  "net/http"
  "fmt"
)

type Response struct {
  Status    string
  Code		int
  Additional	string
}

func main() {
  http.HandleFunc("/grab", grab)
  http.HandleFunc("/status", status)

  http.ListenAndServe("127.0.0.1:8085", nil)
  fmt.Println("Listening for requests")
}



type GrabRequest struct{
	Url		string
}


func grab(w http.ResponseWriter, r *http.Request) {
	
	decoder := json.NewDecoder(r.Body)
	var gr GrabRequest
	err := decoder.Decode(&gr)
	if err != nil{
		response := Response{"error", 400, "Invalid json entered!"}
		js, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
 	} else{
		
		
	  //profile := Profile{"Alex", []string{"snowboarding", "programming"}}
	  response := Response{"ok", 200, "Fetching "+gr.Url+" shortly..."}

	  js, err := json.Marshal(response)
	  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	  }
	  
	    w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	
}