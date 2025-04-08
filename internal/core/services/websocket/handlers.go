package websocket

import (
	"ApiSmart/internal/core/domain/websocket"
	"log"
)

// handleSensorEvent procesa los eventos de sensores
func (s *Server) handleSensorEvent(event websocket.SensorEvent) {
	switch event.Type {
	case websocket.HumidityEvent:
		s.handleHumidityEvent(event)
	case websocket.TemperatureEvent:
		s.handleTemperatureEvent(event)
	case websocket.LightEvent:
		s.handleLightEvent(event)
	case websocket.PHEvent:
		s.handlePHEvent(event)
	case websocket.SystemEvent:
		s.handleSystemStatusEvent(event)
	default:
		log.Printf("Evento desconocido: %v", event.Type)
	}
}

// handleHumidityEvent procesa eventos de humedad
func (s *Server) handleHumidityEvent(event websocket.SensorEvent) {
	if event.Value < 30.0 {
		s.BroadcastEvent(map[string]interface{}{
			"type":    "alert",
			"message": "Humedad baja",
			"value":   event.Value,
			"unit":    event.Unit,
		})
	} else if event.Value > 70.0 {
		s.BroadcastEvent(map[string]interface{}{
			"type":    "alert",
			"message": "Humedad alta",
			"value":   event.Value,
			"unit":    event.Unit,
		})
	}
	s.BroadcastEvent(event)
}

// handleTemperatureEvent procesa eventos de temperatura
func (s *Server) handleTemperatureEvent(event websocket.SensorEvent) {
	if event.Value < 15.0 {
		s.BroadcastEvent(map[string]interface{}{
			"type":    "alert",
			"message": "Temperatura baja",
			"value":   event.Value,
			"unit":    event.Unit,
		})
	} else if event.Value > 30.0 {
		s.BroadcastEvent(map[string]interface{}{
			"type":    "alert",
			"message": "Temperatura alta",
			"value":   event.Value,
			"unit":    event.Unit,
		})
	}
	s.BroadcastEvent(event)
}

// handleLightEvent procesa eventos de luz
func (s *Server) handleLightEvent(event websocket.SensorEvent) {
	if event.Value < 100.0 {
		s.BroadcastEvent(map[string]interface{}{
			"type":    "alert",
			"message": "Luz insuficiente",
			"value":   event.Value,
			"unit":    event.Unit,
		})
	}
	s.BroadcastEvent(event)
}

// handlePHEvent procesa eventos de pH
func (s *Server) handlePHEvent(event websocket.SensorEvent) {
	if event.Value < 5.5 || event.Value > 7.5 {
		s.BroadcastEvent(map[string]interface{}{
			"type":    "alert",
			"message": "pH fuera de rango",
			"value":   event.Value,
			"unit":    event.Unit,
		})
	}
	s.BroadcastEvent(event)
}

// handleSystemStatusEvent procesa eventos de estado del sistema
func (s *Server) handleSystemStatusEvent(event websocket.SensorEvent) {
	// Verificar si Details existe y es del tipo correcto
	if event.Details == nil {
		log.Printf("Error: Details es nil en evento de sistema")
		return
	}

	// Acceder directamente a los valores del mapa
	sensorsOnline, ok1 := event.Details["sensors_online"].(float64)
	totalSensors, ok2 := event.Details["total_sensors"].(float64)

	if !ok1 || !ok2 {
		log.Printf("Error: No se pudieron convertir los valores de sensores")
		return
	}

	if sensorsOnline < totalSensors {
		s.BroadcastEvent(map[string]interface{}{
			"type":    "alert",
			"message": "Sensores desconectados",
			"details": map[string]interface{}{
				"online": sensorsOnline,
				"total":  totalSensors,
			},
		})
	}
	s.BroadcastEvent(event)
}
