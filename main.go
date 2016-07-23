package main

import (
  "encoding/json"
  "net/http"
  "fmt"
  "github.com/nu7hatch/gouuid" //UUID Generator
  "menteslibres.net/gosexy/redis" // Redis client.
)

var rediscon *redis.Client

func main() {

	  //Connect to the Redis server
	  rediscon = redis.New()
	  err := rediscon.Connect("127.0.0.1", 6379)
	  if err != nil{
		panic("Could not connect to Redis")
	  }
	  
	  rediscon.Select(7)
	  
	  msg, err := rediscon.Ping()
	  if err != nil {
			panic("Could not ping Redis")
	  } else {
		fmt.Println("Pinged Redis!:"+msg)
	  }
	  
	  fmt.Println("Connected to Redis")
	  
	  //Register path handlers for http requests
	  http.HandleFunc("/save", save)
	  http.HandleFunc("/progress", progress)
	  http.HandleFunc("/find", find)

	  http.ListenAndServe("127.0.0.1:8085", nil)
	  fmt.Println("Listening for requests")
}





//The JSON we're expecting from the client.
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

type Job struct {
	UUID	string //The job's unique identifier.
	Url		string //The url to archive.
	IP		string //The IP address that requested the job.
}





func save(w http.ResponseWriter, r *http.Request) {
	//Queues a URL up for saving to disk. Once this has been confirmed ok, the file can be assumed saved.

	//Decode the request body.
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
		//JSON decoded and valid, now add it to the queue.
		
		//Generate a new UUID for this job.
		rid, err := uuid.NewV4()
		
		// //Put the job on the redis pub/sub queue and set its status in Redis.
		jobInfo := Job{rid.String(), gr.Url, r.RemoteAddr}
		pushthis, err := json.Marshal(jobInfo)
		rediscon.LPush("jobs", pushthis)
		fmt.Println("Pushed task to Redis!")
		
		//Set the item's status in Redis.
		rediscon.Set("status:"+rid.String(), "in-queue")
		fmt.Println("Set the job status for status:"+rid.String()+" in Redis")
		
		
		
	  //Tell the waiting client all about it.
	  response := SuccessSaving{"ok", 200, "The address "+gr.Url+" has been added to the queue.", rid.String()}

	  js, err := json.Marshal(response)
	  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	  }
	  
	    w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

type ProgressRequest struct{
	UUID string
}

type ProgressResponse struct{
	Status string
}

func progress(w http.ResponseWriter, r *http.Request) {
	//Returns the progress of an archiving job to the waiting API consumer. Success is only stored for one hour before the status is removed from memory.
	/*
	Satuses can be one of: 
	in-queue: The file is waiting in the queue to be saved.
	fetching: The file is being downloaded if image, or being captured if webpage.
	saving: The file is in the process of being saved to disk if png, encoded if jpg / web page, or being converted to webm if gif.
	waiting-deep-freezing: This png file is saved and accessable, but is waiting to be deep-freezed (zopfli).
	deep-freezing: The png file is in the process of being deep freezed with zopfli.
	
	unknown-job: The UUID is invalid or this job finished more than an hour ago. Consult Find to find it.
	saved: The image / webpage was saved successfully.
	*/
	fmt.Println("Recieved progress request.")
	//Decode the request body.
	decoder := json.NewDecoder(r.Body)
	var pr ProgressRequest
	err := decoder.Decode(&pr)
	
	fmt.Println("Decoded JSON")
	
	if err != nil{
		//If there was an error decoding the JSON
		response := ErrorSaving{"error", 400, "Invalid json!"}
		js, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Recieved invalid json")
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
 	} else{
		fmt.Println("Recieved valid progress request.")
		
		
		state, err := rediscon.Get("status:"+pr.UUID)
		if err != nil {
			state = "unknown-job"
		}
		fmt.Println("Status: "+state)
		
		
		response := ProgressResponse{state}
		js, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
	fmt.Println("Got to end of progress.")
	
}

func find(w http.ResponseWriter, r *http.Request) {
	
}