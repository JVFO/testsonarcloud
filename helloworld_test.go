package main

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"testsonarcloud/src/main/schema"
)

func TestName(t *testing.T) {

	nameTest := "name test go"
	eventTest := cloudevents.NewEvent()
	eventTest.SetID(uuid.New().String())
	eventTest.SetSource("knative/eventing/samples/hello-world")
	eventTest.SetType("dev.knative.samples.hifromknative")
	eventTest.SetData(cloudevents.ApplicationJSON, schema.HiFromKnative{Msg: nameTest})

	response, _ := Receive(context.TODO(), eventTest)
	log.Printf("------------------------------------------------")

	log.Printf(string(response.Data()))

	dataevent := &schema.HelloWorld{Msg: string(response.Data())}
	response.DataAs(dataevent)
	log.Printf(dataevent.Msg)
	assert.Equal(t, "Hi "+nameTest, dataevent.Msg)
}
