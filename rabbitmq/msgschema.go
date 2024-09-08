package rabbitmq

import (
	"time"
)

type EntryEvent struct {
	ID            string    `json:"id"`
	VehiclePlate  string    `json:"vehicle_plate"`
	EntryDateTime time.Time `json:"entry_date_time"`
}

type ExitEvent struct {
	ID           string    `json:"id"`
	VehiclePlate string    `json:"vehicle_plate"`
	ExitDateTime time.Time `json:"exit_date_time"`
}

// // SchemaProcessor defines an interface for processing JSON schemas
// type SchemaProcessor interface {
// 	UnmarshalInto(data []byte) (interface{}, error)
// }

// type EntryEventProcessor struct{}

// func (p *EntryEventProcessor) UnmarshalInto(data []byte) (interface{}, error) {
// 	var payload EntryEvent
// 	if err := json.Unmarshal(data, &payload); err != nil {
// 		return nil, err
// 	}
// 	return payload, nil
// }

// type ExitEventProcessor struct{}

// func (p *ExitEventProcessor) UnmarshalInto(data []byte) (interface{}, error) {
// 	var payload ExitEvent
// 	if err := json.Unmarshal(data, &payload); err != nil {
// 		return nil, err
// 	}
// 	return payload, nil
// }
