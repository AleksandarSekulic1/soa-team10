// src/app/services/blog.service.ts

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Blog } from '../models/blog.model'; // Uvozimo naš novi model

@Injectable({
  providedIn: 'root'
})
export class BlogService {
  // URL do tvog blog-service-a. Port je 8082 kao što smo definisali na backendu.
  private apiUrl = 'http://localhost:8082/api/blogs';

  constructor(private http: HttpClient) { }

  /**
   * Dobavlja sve blogove. (GET /api/blogs)
   */
  getAllBlogs(): Observable<Blog[]> {
    // AuthInterceptor će automatski dodati token
    return this.http.get<Blog[]>(this.apiUrl);
  }

  /**
   * Dobavlja jedan specifičan blog po ID-u. (GET /api/blogs/:id)
   * @param id ID bloga koji se traži
   */
  getBlogById(id: string): Observable<Blog> {
    return this.http.get<Blog>(`${this.apiUrl}/${id}`);
  }

  /**
   * Kreira novi blog. (POST /api/blogs)
   * @param blogData Podaci za novi blog (samo title i content su obavezni)
   */
  createBlog(blogData: { Title: string; Content: string; Images?: string[] }): Observable<Blog> {
    return this.http.post<Blog>(this.apiUrl, blogData);
  }

  /**
   * Ažurira postojeći blog. (PUT /api/blogs/:id)
   * @param id ID bloga koji se ažurira
   * @param blogData Podaci koji se menjaju
   */
  updateBlog(id: string, blogData: Partial<Blog>): Observable<Blog> {
    return this.http.put<Blog>(`${this.apiUrl}/${id}`, blogData);
  }

  /**
   * Dodaje komentar na blog. (POST /api/blogs/:id/comments)
   * @param blogId ID bloga na koji se dodaje komentar
   * @param commentText Tekst komentara
   */
  addComment(blogId: string, commentText: { Text: string }): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/${blogId}/comments`, commentText);
  }

    updateComment(blogId: string, commentId: string, commentData: { Text: string }): Observable<any> {
    return this.http.put<any>(`${this.apiUrl}/${blogId}/comments/${commentId}`, commentData);
  }

  /**
   * Lajkuje ili unlajkuje blog. (POST /api/blogs/:id/likes)
   * @param blogId ID bloga za lajkovanje
   */
  toggleLike(blogId: string): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/${blogId}/likes`, {});
  }
}