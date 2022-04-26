package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)


var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	options := client.OptionsReader()
	fmt.Println("Connected to: v", options.Servers())
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}


func main() {
	var broker = "mqtt-integration.sandbox.drogue.cloud"
	var port = 443
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetUsername("drogue-public-temperature")
	opts.SetPassword("public")
	
	//callback functions
	opts.SetDefaultPublishHandler(messageHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectionLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	
	sub(client)
	publish(client)
	
	client.Disconnect(250)
}



//publish function
func publish(client mqtt.Client) {
	num := 1
	for i := 1; i <= num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish("drogue-public-temperature/go", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}


//Subscribe function
func sub(client mqtt.Client) {
	topic := "drogue-public-temperature/#"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}