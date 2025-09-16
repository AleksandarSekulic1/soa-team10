// src/app/models/blog.model.ts

export interface Like {
  UserID: string;
  CreatedAt: string;
}

export interface Comment {
  ID?: string; // <-- Mora biti veliko 'ID'
  AuthorID: string;
  Text: string; // <-- Mora biti veliko 'Text'
  CreatedAt: string;
  LastUpdatedAt: string;
}

export interface Blog {
  ID: string;
  Title: string;
  Content: string;
  AuthorID: string;
  CreatedAt: string;
  Images?: string[];
  Comments: Comment[];
  Likes: Like[];
}