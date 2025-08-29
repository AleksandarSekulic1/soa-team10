// src/app/models/tour-execution.model.ts

export interface CompletedKeyPoint {
  KeyPointId: string;
  CompletionTime: Date;
}

export interface TourExecution {
  ID: string;
  TourId: string;
  UserId: string;
  Status: 'Active' | 'Completed' | 'Abandoned';
  CompletedKeyPoints: CompletedKeyPoint[];
  LastActivity: Date;
  StartTime: Date;
  EndTime?: Date;
}