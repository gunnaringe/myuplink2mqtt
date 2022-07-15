package homeassistant

type Availability struct {
	Topic string `json:"topic,omitempty"`
}

type Sensor struct {
	Availability        []Availability `json:"availability,omitempty"`
	DeviceClass         string         `json:"device_class,omitempty"`
	EnabledByDefault    bool           `json:"enabled_by_default,omitempty"`
	JSONAttributesTopic string         `json:"json_attributes_topic,omitempty"`
	Name                string         `json:"name,omitempty"`
	StateClass          string         `json:"state_class,omitempty"`
	StateTopic          string         `json:"state_topic,omitempty"`
	UniqueID            string         `json:"unique_id,omitempty"`
	UnitOfMeasurement   string         `json:"unit_of_measurement,omitempty"`
	ValueTemplate       string         `json:"value_template,omitempty"`
}

type Selection struct {
	Availability        []Availability `json:"availability,omitempty"`
	UniqueID            string         `json:"unique_id,omitempty"`
	Name                string         `json:"name,omitempty"`
	EnabledByDefault    bool           `json:"enabled_by_default,omitempty"`
	Icon                string         `json:"icon,omitempty"`
	CommandTemplate     string         `json:"command_template,omitempty"`
	CommandTopic        string         `json:"command_topic,omitempty"`
	JSONAttributesTopic string         `json:"json_attributes_topic,omitempty"`
	Optimistic          bool           `json:"optimistic,omitempty"`
	Options             []string       `json:"options,omitempty"`
	StateTopic          string         `json:"state_topic,omitempty"`
	ValueTemplate       string         `json:"value_template,omitempty"`
}

type Hvac struct {
	Availability               []Availability `json:"availability,omitempty"`
	UniqueID                   string         `json:"unique_id,omitempty"`
	Name                       string         `json:"name,omitempty"`
	ActionTopic                string         `json:"action_topic,omitempty"`
	ActionTemplate             string         `json:"action_template,omitempty"`
	TemperatureCommandTopic    string         `json:"temperature_command_topic,omitempty"`
	TemperatureCommandTemplate string         `json:"temperature_command_template,omitempty"`
	TemperatureStateTopic      string         `json:"temperature_state_topic,omitempty"`
	TemperatureStateTemplate   string         `json:"temperature_state_template,omitempty"`
	CurrentTemperatureTopic    string         `json:"current_temperature_topic,omitempty"`
	CurrentTemperatureTemplate string         `json:"current_temperature_template,omitempty"`
	MinTemp                    float64        `json:"min_temp,omitempty"`
	MaxTemp                    float64        `json:"max_temp,omitempty"`
	TempStep                   int            `json:"temp_step,omitempty"`
	TemperatureUnit            string         `json:"temperature_unit,omitempty"`
	ModeStateTopic             string         `json:"mode_state_topic,omitempty"`
	ModeStateTemplate          string         `json:"mode_state_template,omitempty"`
	Modes                      []string       `json:"modes,omitempty"`
}
