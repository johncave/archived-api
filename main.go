package main

import (
  "encoding/json"
  "net/http"
  "fmt"
  "github.com/nu7hatch/gouuid"
)

func main() {
  http.HandleFunc("/save", save)
  http.HandleFunc("/progress", progress)
  http.HandleFunc("/find", find)

  http.ListenAndServe("127.0.0.1:8085", nil)
  fmt.Println("Listening for requests")
}



type GrabRequest struct{
	Url		string
}

//If the call to action was successful, return one of these.
type SuccessSaving struct {
  Status    string
  Code		int
  Additional	string
  UUID		string
}

//Otherwise return one of these.
type ErrorSaving struct {
  Status    string
  Code		int
  Additional	string
}

func save(w http.ResponseWriter, r *http.Request) {
	//Queues a URL up for saving to disk. Once this has been confirmed ok, the file can be assumed saved.
	/*
	
	
	*/
	
	decoder := json.NewDecoder(r.Body)
	var gr GrabRequest
	err := decoder.Decode(&gr)
	
	if err != nil{
		//If there was an error deconding the JSON
		response := ErrorSaving{"error", 400, "Invalid json!"}
		js, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
 	} else{
		//JSON decoded, now add it to the queue.
		
		//Generate a new UUID for this job.
		rid, err := uuid.NewV4()
		
		//Put the job on the redis pub/sub queue.
		
		
		
		
	  //Tell the waiting client all about it.
	  response := SuccessSaving{"ok", 200, "Fetching "+gr.Url+" shortly...", rid.String()}

	  js, err := json.Marshal(response)
	  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	  }
	  
	    w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func progress(w http.ResponseWriter, r *http.Request) {
	//Returns the progress of an archiving job to the waiting API consumer. 
	/*
	Satuses can be one of: 
	in-queue: The file is waiting in the queue to be saved.
	fetching: The file is being downloaded if image, or being captured if webpage.
	saving: The file is in the process of being saved to disk if png, encoded if jpg / web page, or being converted to webm if gif.
	waiting-deep-freezing: This png file is saved and accessable, but is waiting to be deep-freezed (zopfli).
	deep-freezing: The png file is in the process of being deep freezed with zopfli.
	*/
	
}

func find(w http.ResponseWriter, r *http.Request) {
	
}