package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlekSi/pointer"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gunnaringe/myuplink2mqtt/internal/homeassistant"
	"github.com/gunnaringe/myuplink2mqtt/pkg/auth"
	"github.com/gunnaringe/myuplink2mqtt/pkg/myuplink"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var DefaultClient = http.Client{
	Timeout: 10 * time.Second,
}
var CommandTopicRegex = regexp.MustCompile(".*/(?P<deviceId>\\S+)/set")

type Server struct {
	logger       *zap.SugaredLogger
	clientId     string
	clientSecret string
	client       *myuplink.ClientWithResponses
	mqttClient   mqtt.Client
	mqttServer   string
	Auth         *auth.Auth
}

func New(clientId, clientSecret, mqttServer string, logger *zap.Logger) *Server {
	return &Server{
		clientId:     clientId,
		clientSecret: clientSecret,
		mqttServer:   mqttServer,
		logger:       logger.Sugar(),
	}
}

func (r *Server) Run() error {
	r.logger.Info("Starting server")

	r.Auth = auth.New(r.clientId, r.clientSecret)

	c, err := myuplink.NewClientWithResponses(
		"https://api.myuplink.com",
		myuplink.WithRequestEditorFn(r.Auth.Intercept),
		myuplink.WithHTTPClient(&DefaultClient),
	)
	if err != nil {
		return err
	}
	r.client = c

	cbck := func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

		update := &Update{}
		err := json.Unmarshal(msg.Payload(), update)
		if err != nil {
			r.logger.Warnf("Could not unmarshal payload: topic=%s", msg.Topic())
			return
		}

		matches := findNamedMatches(CommandTopicRegex, msg.Topic())
		deviceId := matches["deviceId"]
		r.SetTargetTemp(deviceId, update.TargetTemp)
	}

	CreateMqtt(r.mqttServer, cbck)

	deviceIds := r.getDeviceIds()
	r.logger.Infof("Found %v devices", len(deviceIds))
	for _, deviceId := range deviceIds {
		r.logger.Infof("Found device: %v", deviceId)
	}

	for _, deviceId := range deviceIds {
		name := "Høiax Connected"
		stateTopic := fmt.Sprintf("myuplink2mqtt/%s", deviceId)
		commandTopic := fmt.Sprintf("myuplink2mqtt/%s/set", deviceId)
		availability := &homeassistant.Availability{
			Topic: "myuplink2mqtt/state",
		}
		hvac := &homeassistant.Hvac{
			Availability:               []homeassistant.Availability{*availability},
			UniqueID:                   fmt.Sprintf("%s_%s_%s", name, "energy", "myuplink2mqtt"),
			Name:                       name,
			ActionTopic:                stateTopic,
			ActionTemplate:             "{{ value_json.action }}",
			TemperatureCommandTopic:    commandTopic,
			TemperatureCommandTemplate: "{\"temperature\": {{value}} }",
			TemperatureStateTopic:      stateTopic,
			TemperatureStateTemplate:   "{{ value_json.target_temp }}",
			CurrentTemperatureTopic:    stateTopic,
			CurrentTemperatureTemplate: "{{ value_json.current_temp }}",
			MinTemp:                    20,
			MaxTemp:                    85,
			TempStep:                   1,
			TemperatureUnit:            "C",
			ModeStateTopic:             stateTopic,
			ModeStateTemplate:          "{{ value_json.mode }}",
			Modes:                      []string{"off", "heat"},
		}
		PublishHvac(deviceId, hvac)

		energy := &homeassistant.Sensor{
			Availability:        []homeassistant.Availability{*availability},
			DeviceClass:         "energy",
			EnabledByDefault:    true,
			JSONAttributesTopic: stateTopic,
			Name:                fmt.Sprintf("%s energy", name),
			StateClass:          "total_increasing",
			StateTopic:          stateTopic,
			UniqueID:            fmt.Sprintf("%s_%s_%s", name, "energy", "myuplink2mqtt"),
			UnitOfMeasurement:   "kWh",
			ValueTemplate:       "{{ value_json.energy }}",
		}
		PublishSensor(deviceId, energy, "energy")

		power := &homeassistant.Sensor{
			Availability:        []homeassistant.Availability{*availability},
			DeviceClass:         "power",
			EnabledByDefault:    true,
			JSONAttributesTopic: stateTopic,
			Name:                fmt.Sprintf("%s estimated power", name),
			StateClass:          "measurement",
			StateTopic:          stateTopic,
			UniqueID:            fmt.Sprintf("%s_%s_%s", name, "power", "myuplink2mqtt"),
			UnitOfMeasurement:   "W",
			ValueTemplate:       "{{ value_json.power }}",
		}
		PublishSensor(deviceId, power, "power")

		storedEnergy := &homeassistant.Sensor{
			Availability:        []homeassistant.Availability{*availability},
			DeviceClass:         "energy",
			EnabledByDefault:    true,
			JSONAttributesTopic: stateTopic,
			Name:                fmt.Sprintf("%s stored energy", name),
			StateClass:          "measurement",
			StateTopic:          stateTopic,
			UniqueID:            fmt.Sprintf("%s_%s_%s", name, "stored_energy", "myuplink2mqtt"),
			UnitOfMeasurement:   "kWh",
			ValueTemplate:       "{{ value_json.stored_energy }}",
		}
		PublishSensor(deviceId, storedEnergy, "stored-energy")
	}

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, deviceId := range deviceIds {
					r.reportStatus(deviceId)
				}
			}
		}
	}()

	select {}
}

func (r *Server) getDeviceIds() []string {
	params := &myuplink.GetV2SystemsMeParams{
		Page:         pointer.ToInt32(1),
		ItemsPerPage: pointer.ToInt32(100),
	}

	resp, err := r.client.GetV2SystemsMe(context.Background(), params)
	if err != nil {
		r.logger.Fatalw(
			"Could not load systems - Shutting down",
			zap.Error(err),
		)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	var result myuplink.PagedSystemResult
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		r.logger.Fatalw(
			"Could not load systems - Shutting down",
			zap.String("response", string(bodyBytes)),
			zap.Error(err),
		)
	}

	var deviceIds []string
	for _, systems := range *result.Systems {
		for _, device := range *systems.Devices {
			id := *device.Id
			deviceIds = append(deviceIds, id)
		}
	}
	return deviceIds
}

func (r *Server) reportStatus(deviceId string) {
	params := &myuplink.GetV2DevicesDeviceIdPointsParams{}
	resp, err := r.client.GetV2DevicesDeviceIdPointsWithResponse(context.Background(), deviceId, params)
	if err != nil {
		r.logger.Warn("Error fetching status", zap.String("device", deviceId))
	}
	result := *resp.JSON200
	parameters := make(map[string]myuplink.ParameterData)
	for _, data := range result {
		parameters[*data.ParameterId] = data
	}

	var action string
	if *parameters["505"].Value == 1 || *parameters["506"].Value == 1 {
		action = "heating"
	} else {
		action = "idle"
	}

	maxPower := *parameters["517"].Value
	var mode string
	if maxPower > 0 {
		mode = "heat"
	} else {
		mode = "off"
	}

	status := &Status{
		Energy:       *parameters["303"].Value,
		StoredEnergy: *parameters["302"].Value,
		TargetTemp:   *parameters["527"].Value,
		CurrentTemp:  *parameters["528"].Value,
		Power:        *parameters["400"].Value,
		MaxPower:     maxPower,
		Mode:         mode,
		Action:       action,
	}
	PublishStatus(deviceId, status)
	PublishState()

	fmt.Println("")
	fmt.Println("======================================================")
	fmt.Printf("Timestamp: %v kWh\n", *parameters["303"].Timestamp)
	fmt.Printf("Energy used: %v kWh\n", *parameters["303"].Value)
	fmt.Printf("Stored used: %v kWh\n", *parameters["302"].Value)
	fmt.Printf("Current power (estimated): %v W\n", *parameters["400"].Value)
	fmt.Printf("Element 1: %v\n", *parameters["505"].StrVal)
	fmt.Printf("Element 2: %v\n", *parameters["506"].StrVal)
	fmt.Printf("Max effect: %v\n", *parameters["517"].StrVal)
	fmt.Printf("Temperature set: %v\n", *parameters["527"].StrVal)
	fmt.Printf("Temperature current: %v\n", *parameters["528"].StrVal)
}

func (r *Server) SetTargetTemp(deviceId string, targetTemp float64) {
	r.logger.Infof("Setting target temperature: %g °C", targetTemp)
	body, _ := json.Marshal(map[string]string{
		"527": fmt.Sprintf("%g", targetTemp),
	})
	url := fmt.Sprintf("https://api.myuplink.com/v2/devices/%s/points", deviceId)
	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json-patch+json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.Auth.Token()))
	if err != nil {
		r.logger.Fatalf("Error while making request to update heater: %v", err)
	}

	response, err := DefaultClient.Do(request)
	if err != nil {
		r.logger.Warn("Could not update temperature", zap.String("device", deviceId))
		return
	}

	if response.StatusCode != 200 {
		r.logger.Warn("Could not update temperature",
			zap.String("device", deviceId),
			zap.String("http_status", response.Status),
		)
	}
}

type Status struct {
	Energy       float64 `json:"energy,omitempty"`
	StoredEnergy float64 `json:"stored_energy,omitempty"`
	TargetTemp   float64 `json:"target_temp,omitempty"`
	CurrentTemp  float64 `json:"current_temp,omitempty"`
	Power        float64 `json:"power,omitempty"`
	MaxPower     float64 `json:"max_power,omitempty"`
	Mode         string  `json:"mode,omitempty"`
	Action       string  `json:"action,omitempty"`
}

func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}
