package main

import (
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.16.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
)

type PaymentRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

func handlePayment(w http.ResponseWriter, r *http.Request) {

	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	log.Printf("Extracted Trace Context: %v", trace.SpanContextFromContext(ctx).TraceID())

	_, span := otel.Tracer("payment-tracer").Start(ctx, "handlePayment", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	span.SetAttributes(
		semconv.HTTPMethodKey.String(r.Method),
		semconv.HTTPTargetKey.String(r.URL.Path),
	)

	log.Print("Received payment request")

	var paymentReq PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentReq); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	fmt.Printf("Processing payment for Order ID: %s, Amount: %.2f\n", paymentReq.OrderID, paymentReq.Amount)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payment processed successfully"))
}

func main() {
	initTracer()
	http.Handle("/v1/payments", otelhttp.NewHandler(http.HandlerFunc(handlePayment), "handlePayment"))

	fmt.Println("Payment Service running on port 8081...")
	http.ListenAndServe(":8081", nil)
}
