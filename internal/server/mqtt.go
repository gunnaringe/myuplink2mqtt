package server

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gunnaringe/myuplink2mqtt/internal/homeassistant"
	"time"
)

var client mqtt.Client

type Update struct {
	TargetTemp float64 `json:"temperature,omitempty"`
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
	if token := client.Subscribe("myuplink2mqtt/+/set", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func CreateMqtt(server string, callback mqtt.MessageHandler) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)
	opts.SetDefaultPublishHandler(callback)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetResumeSubs(true)

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func PublishStatus(deviceId string, status *Status) {
	topic := fmt.Sprintf("myuplink2mqtt/%s", deviceId)
	Publish(topic, status, false)
}

func PublishHvac(deviceId string, config *homeassistant.Hvac) {
	topic := fmt.Sprintf("homeassistant/climate/%s/config", deviceId)
	Publish(topic, config, true)
}

func PublishSensor(deviceId string, config *homeassistant.Sensor, sensor string) {
	topic := fmt.Sprintf("homeassistant/sensor/%s/%s/config", deviceId, sensor)
	Publish(topic, config, true)
}

func Publish(topic string, v any, retained bool) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	if token := client.Publish(topic, 1, retained, b); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func PublishState() {
	if token := client.Publish("myuplink2mqtt/state", 0, false, "online"); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
