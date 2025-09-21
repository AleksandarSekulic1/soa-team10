import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { AuthService } from './auth.service';

export interface User {
  id: string;
  username: string;
  email: string;
  firstName?: string;
  lastName?: string;
}

export interface FollowRequest {
  followerId: string;
  followingId: string;
}

export interface UserRecommendation {
  user: User;
  mutualFollowers: number;
  reason: string;
}

export interface FollowResponse {
  followers: User[];
  following: User[];
  isFollowing: boolean;
}

@Injectable({
  providedIn: 'root'
})
export class FollowerService {
  private apiUrl = 'http://localhost:8000/api/followers';

  constructor(
    private http: HttpClient,
    private authService: AuthService
  ) {}

  private getHeaders(): HttpHeaders {
    const token = localStorage.getItem('jwt_token');
    return new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    });
  }

  // Kreiraj korisnika u follower servisu
  createUser(user: User): Observable<any> {
    return this.http.post(`${this.apiUrl}/api/users`, user, { headers: this.getHeaders() });
  }

  // Prati korisnika
  followUser(followingId: string): Observable<any> {
    const currentUser = this.authService.currentUser.getValue();
    if (!currentUser) throw new Error('User not logged in');
    
    const followRequest: FollowRequest = {
      followerId: currentUser.id,
      followingId: followingId
    };
    return this.http.post(`${this.apiUrl}/api/follow`, followRequest, { headers: this.getHeaders() });
  }

  // Otprati korisnika
  unfollowUser(followingId: string): Observable<any> {
    const currentUser = this.authService.currentUser.getValue();
    if (!currentUser) throw new Error('User not logged in');
    
    const unfollowRequest: FollowRequest = {
      followerId: currentUser.id,
      followingId: followingId
    };
    return this.http.post(`${this.apiUrl}/api/unfollow`, unfollowRequest, { headers: this.getHeaders() });
  }

  // Dobij pratiloce korisnika
  getFollowers(userId: string): Observable<User[]> {
    return this.http.get<User[]>(`${this.apiUrl}/api/users/${userId}/followers`, { headers: this.getHeaders() });
  }

  // Dobij korisnike koje prati
  getFollowing(userId: string): Observable<User[]> {
    return this.http.get<User[]>(`${this.apiUrl}/api/users/${userId}/following`, { headers: this.getHeaders() });
  }

  // Proveri da li prati korisnika
  isFollowing(userId: string): Observable<{ isFollowing: boolean }> {
    const currentUser = this.authService.currentUser.getValue();
    if (!currentUser) throw new Error('User not logged in');
    
    return this.http.get<{ isFollowing: boolean }>(
      `${this.apiUrl}/api/users/${userId}/is-following?followerId=${currentUser.id}`, 
      { headers: this.getHeaders() }
    );
  }

  // Dobij preporuke korisnika
  getRecommendations(): Observable<UserRecommendation[]> {
    const currentUser = this.authService.currentUser.getValue();
    if (!currentUser) throw new Error('User not logged in');
    
    return this.http.get<UserRecommendation[]>(
      `${this.apiUrl}/api/recommendations?userId=${currentUser.id}`, 
      { headers: this.getHeaders() }
    );
  }

  // Dobij ID-jeve praćenih korisnika (za blog filtriranje)
  getFollowedUserIds(): Observable<string[]> {
    const currentUser = this.authService.currentUser.getValue();
    if (!currentUser) throw new Error('User not logged in');
    
    return this.http.get<string[]>(
      `${this.apiUrl}/api/followed-users?userId=${currentUser.id}`, 
      { headers: this.getHeaders() }
    );
  }

  // Dobij sve korisnike iz follower servisa
  getAllUsers(): Observable<User[]> {
    return this.http.get<User[]>(`${this.apiUrl}/api/users`, { headers: this.getHeaders() });
  }
}
