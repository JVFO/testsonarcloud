package main

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"testsonarcloud/src/main/schema"
)

var (
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
	DebugLogger *log.Logger
)

func init() {
	file, err := os.OpenFile("helloworld-go.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Receive(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	// Aquí es donde irá su código para procesar el evento.
	// En este ejemplo registraremos el mensaje de evento
	log.Printf("Event received. \n%s\n", event.Context)
	log.Printf("Event data. \n%s\n", event)

	dataevent := &schema.HelloWorld{Msg: string(event.Data())}
	if err := event.DataAs(dataevent); err != nil {
		log.Printf("Error while extracting cloudevent Data: %s\n", err.Error())
		return nil, cloudevents.NewHTTPResult(400, "failed to convert data: %s", err)
	}
	log.Printf("Hello World Message from received event %q", dataevent.Msg)

	ErrorLogger.Println("Something went wrong")
	InfoLogger.Println("Something noteworthy happened")
	DebugLogger.Println("Useful debugging information.")

	// Responder con otro evento (opcional)
	// Esto es opcional y pretende mostrar cómo responder con otro evento después del procesamiento.
	// La respuesta volverá al sistema de eventos knative como cualquier otro evento.
	newEvent := cloudevents.NewEvent()
	newEvent.SetID(uuid.New().String())
	newEvent.SetSource("knative/eventing/samples/hello-world")
	newEvent.SetType("dev.knative.samples.hifromknative")
	if err := newEvent.SetData(cloudevents.ApplicationJSON, schema.HiFromKnative{Msg: "Hi " + dataevent.Msg}); err != nil {
		return nil, cloudevents.NewHTTPResult(500, "failed to set response data: %s", err)
	}
	log.Printf("Responding with event\n%s\n", newEvent)
	return &newEvent, nil
}

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

func main() {
	log.Print("Hello world sample started.")

	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}

	ctx := context.Background()

	p, err := cloudevents.NewHTTP(cloudevents.WithPort(env.Port), cloudevents.WithPath(env.Path))
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}
	c, err := cloudevents.NewClient(p,
		cloudevents.WithUUIDs(),
		cloudevents.WithTimeNow(),
	)
	if err != nil {
		log.Fatalf("failed to create client: %s", err.Error())
	}

	if err := c.StartReceiver(ctx, Receive); err != nil {
		log.Fatalf("failed to start receiver: %s", err.Error())
	}

	log.Printf("listening on :%d%s\n", env.Port, env.Path)
	<-ctx.Done()
}
