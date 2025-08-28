// services/tours-service/api/tours_grpc_handler.go
package api

import (
	"context"
	"tours-service/service"
	"tours-service/proto/tours" // Uvozimo generisani proto kod
)

type ToursGrpcHandler struct {
	tours.UnimplementedToursServiceServer // Obavezno za kompatibilnost
	service                             service.TourService
}

func NewToursGrpcHandler(service service.TourService) *ToursGrpcHandler {
	return &ToursGrpcHandler{service: service}
}

// Implementiramo GetTourById metodu definisanu u .proto fajlu
func (h *ToursGrpcHandler) GetTourById(ctx context.Context, req *tours.GetTourByIdRequest) (*tours.GetTourByIdResponse, error) {
	// Pozivamo postojeći servis da dobavimo turu iz baze
	tour, err := h.service.GetById(req.TourId)
	if err != nil {
		return nil, err
	}

	// Mapiramo naš domain.Tour na proto.GetTourByIdResponse
	response := &tours.GetTourByIdResponse{
		Id:          tour.ID.Hex(),
		AuthorId:    tour.AuthorId,
		Name:        tour.Name,
		Description: tour.Description,
		Difficulty:  int32(tour.Difficulty),
		Tags:        tour.Tags,
		Status:      string(tour.Status),
		Price:       tour.Price,
		Distance:    tour.Distance,
		KeyPoints:   []*tours.KeyPoint{},
		TransportInfo: []*tours.TourTransport{},
	}

	for _, kp := range tour.KeyPoints {
		response.KeyPoints = append(response.KeyPoints, &tours.KeyPoint{
			Id:          kp.ID.Hex(),
			TourId:      kp.TourId,
			Name:        kp.Name,
			Description: kp.Description,
			Latitude:    kp.Latitude,
			Longitude:   kp.Longitude,
			ImageUrl:    kp.ImageUrl,
		})
	}

	for _, ti := range tour.TransportInfo {
		response.TransportInfo = append(response.TransportInfo, &tours.TourTransport{
			Type:          string(ti.Type),
			TimeInMinutes: int32(ti.TimeInMinutes),
		})
	}
	
	return response, nil
}