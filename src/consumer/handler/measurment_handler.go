package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Measurement struct {
	Temperature float32 `json:"t"`
	Humidity    float32 `json:"h"`
	Timestamp   int64   `json:"ts"`
}

type MeasurementStorage interface {
	Store(
		ctx context.Context,
		temperature float32,
		humidity float32,
		time time.Time,
	) error
}

type MeasurementHandler struct {
	storage MeasurementStorage
}

func New(storage MeasurementStorage) *MeasurementHandler {
	return &MeasurementHandler{
		storage: storage,
	}
}

func (h *MeasurementHandler) Handle(ctx context.Context, message amqp.Delivery) error {
	var data Measurement

	if err := json.Unmarshal(message.Body[:], &data); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	if err := h.storage.Store(ctx, data.Temperature, data.Humidity, time.Unix(data.Timestamp, 0)); err != nil {
		return fmt.Errorf("store measurement: %w", err)
	}

	return nil
}
