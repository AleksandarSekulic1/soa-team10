// src/app/models/blog.model.ts


export interface Like {
  userId?: string;
  createdAt?: string;
  UserID?: string;
  CreatedAt?: string;
}


export interface Comment {
  _id?: string;
  authorId?: string;
  text?: string;
  createdAt?: string;
  lastUpdatedAt?: string;
  // legacy/pascal
  ID?: string;
  AuthorID?: string;
  Text?: string;
  CreatedAt?: string;
  LastUpdatedAt?: string;
}


export interface Blog {
  _id?: string;
  title?: string;
  content?: string;
  authorId?: string;
  createdAt?: string;
  images?: string[];
  comments?: Comment[];
  likes?: Like[];
  // legacy/pascal
  ID?: string;
  Title?: string;
  Content?: string;
  AuthorID?: string;
  CreatedAt?: string;
  Images?: string[];
  Comments?: Comment[];
  Likes?: Like[];
}