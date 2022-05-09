package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
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

func NewTlsConfig() *tls.Config {
    certpool := x509.NewCertPool()
    ca, err := ioutil.ReadFile("ca-cert.pem")
    if err != nil {
        log.Fatalln(err.Error())
    }
    certpool.AppendCertsFromPEM(ca)
	
    // Import client certificate/key pair
    clientKeyPair, err := tls.LoadX509KeyPair("client.crt", "key.unencrypted.pem")
    if err != nil {
        panic(err)
    }
    return &tls.Config{
        RootCAs: certpool,
        ClientAuth: tls.NoClientCert,
        ClientCAs: nil,
        InsecureSkipVerify: true,
        Certificates: []tls.Certificate{clientKeyPair},
    }
}


func main() {
	var broker = "broker.hivemq.com"
    var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetUsername("drogue-public-temperature")
	opts.SetPassword("public")
	opts.SetTLSConfig(NewTlsConfig())
	
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
	num := 5
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
