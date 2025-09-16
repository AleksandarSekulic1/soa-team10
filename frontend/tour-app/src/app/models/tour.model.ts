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

// NOVO: Definišemo oblik informacije o transportu
export interface TourTransport {
  type: 'walking' | 'bicycle' | 'car';
  timeInMinutes: number;
}

// Definišemo glavni oblik Ture
export interface Tour {
  id: string;
  authorId: string;
  name: string;
  description: string;
  difficulty: number;
  tags: string[];
  status: 'draft' | 'published' | 'archived'; // Status je sada striktno definisan
  price: number;
  reviews: TourReview[];
  keyPoints: TourKeyPoint[];
  
  // --- NOVA POLJA ---
  distance: number; // Distanca u kilometrima
  transportInfo: TourTransport[]; // Lista vremena putovanja
  publishedAt?: Date; // Opciono polje za vreme objave
  archivedAt?: Date;  // Opciono polje za vreme arhiviranja
}