// src/app/models/tour.model.ts

// Definišemo oblik ključne tačke
export interface TourKeyPoint {
  id: string;
  tourId: string;
  name: string;
  description: string;
  latitude: number;
  longitude: number;
  imageUrl: string;
}

// Definišemo oblik recenzije
export interface TourReview {
  id: string;
  rating: number;
  comment: string;
  touristId: string;
  visitDate: Date;
  commentDate: Date;
  imageUrls: string[];
}

// Definišemo glavni oblik Ture
export interface Tour {
  id: string;
  authorId: string;
  name: string;
  description: string;
  difficulty: number;
  tags: string[];
  status: string;
  price: number;
  reviews: TourReview[];
  keyPoints: TourKeyPoint[];
}