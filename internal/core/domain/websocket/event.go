package websocket

import "time"

// EventType representa el tipo de evento WebSocket
type EventType string

const (
	// Tipos de eventos de sensores
	HumidityEvent    EventType = "humidity"
	TemperatureEvent EventType = "temperature"
	LightEvent       EventType = "light"
	PHEvent          EventType = "ph"
	SystemEvent      EventType = "system_status"
)

// SensorEvent representa un evento de sensor
type SensorEvent struct {
	Type      EventType              `json:"type"`
	Value     float64                `json:"value"`
	Unit      string                 `json:"unit"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// SystemStatus representa el estado del sistema
type SystemStatus struct {
	Status         string `json:"status"`
	SensorsOnline  int    `json:"sensors_online"`
	TotalSensors   int    `json:"total_sensors"`
	LastUpdateTime string `json:"last_update_time"`
}
