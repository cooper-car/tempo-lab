package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.16.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
)

const paymentServiceURL = "http://payment_service:8081/v1/payments"

type OrderRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

func handleOrder(w http.ResponseWriter, r *http.Request) {

	ctx, span := otel.Tracer("order-tracer").Start(r.Context(), "handleOrder", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	span.SetAttributes(
		semconv.HTTPMethodKey.String(r.Method),
		semconv.HTTPTargetKey.String(r.URL.Path),
	)
	log.Print("Received order request")

	var orderReq OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(orderReq)
	if err != nil {
		span.RecordError(err)
		http.Error(w, "Failed to marshal order request", http.StatusInternalServerError)
		return
	}

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, err := http.NewRequestWithContext(ctx, "POST", paymentServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		span.RecordError(err)
		http.Error(w, "Failed to create payment request", http.StatusInternalServerError)
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	req.Header.Set("Content-Type", "application/json")
	log.Printf("Injected Headers: %v", req.Header) // Log injected headers for verification

	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		http.Error(w, "Failed to process payment", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	//w.WriteHeader(resp.StatusCode)
	//body, _ := ioutil.ReadAll(resp.Body)
	//w.Write(body)
}

// entry
func main() {
	initTracer()
	http.HandleFunc("/v1/orders", handleOrder)
	fmt.Println("Order Service running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
