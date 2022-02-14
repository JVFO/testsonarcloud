package testsonarcloud

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"testsonarcloud/src/main/schema"
)

func Receive(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	// Here is where your code to process the event will go.
	// In this example we will log the event msg
	log.Printf("Event received. \n%s\n", event)
	data := &schema.HelloWorld{Msg: string(event.Data())}

	if err := event.DataAs(data); err != nil {
		log.Printf("Error while extracting cloudevent Data: %s\n", err.Error())
		return nil, cloudevents.NewHTTPResult(400, "failed to convert data: %s", err)
	}
	log.Printf("Hello World Message from received event %q", data.Msg)

	// Respond with another event (optional)
	// This is optional and is intended to show how to respond back with another event after processing.
	// The response will go back into the knative eventing system just like any other event
	newEvent := cloudevents.NewEvent()
	// Setting the ID here is not necessary. When using NewDefaultClient the ID is set
	// automatically. We set the ID anyway so it appears in the log.
	newEvent.SetID(uuid.New().String())
	newEvent.SetSource("knative/eventing/samples/hello-world")
	newEvent.SetType("dev.knative.samples.hifromknative")
	if err := newEvent.SetData(cloudevents.ApplicationJSON, schema.HiFromKnative{Msg: "Hi from helloworld-go app!"}); err != nil {
		return nil, cloudevents.NewHTTPResult(500, "failed to set response data: %s", err)
	}
	log.Printf("Responding with event\n%s\n", newEvent)
	return &newEvent, nil
}

func main() {

	//	r := mux.NewRouter()
	//	// Route handles & endpoint
	//	//r.HandleFunc("/", getEvent)
	//	r.HandleFunc("/", getEvent).Methods("POST")
	//	log.Fatal(http.ListenAndServe(":8000", r))
	//
	//}
	r := mux.NewRouter()
	r.HandleFunc("/", getEvent).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))

}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println("one=" + r.FormValue("one"))
	log.Println("two=" + r.FormValue("two"))
	fmt.Fprintf(w, "Gorilla!\n")
}

func getEvent(w http.ResponseWriter, r *http.Request) {

	log.Println(r)
	fmt.Fprintf(w, "Gorilla!\n")

	//params := json.NewDecoder(request.Body)
	//fmt.Println("jkkkkkkkk")
	//fmt.Println(params)
	//fmt.Fprintln(writer, "test")
	//{
	//	"Content-Type":"application/json",
	//	"ce-specversion": "0.3",
	//	"ce-id":str(uuid.uuid4()),
	//	"ce-type": "dev.knative.samples.hifromknative",
	//	"ce-source": "knative/eventing/samples/hello-world"
	//},

	log.Print("Hello world sample started.")
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	log.Fatal(c.StartReceiver(context.Background(), Receive))
}
