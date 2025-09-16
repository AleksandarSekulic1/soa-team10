import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { FollowerService, User, UserRecommendation } from '../../services/follower.service';
import { AuthService } from '../../services/auth.service';
import { FollowButtonComponent } from '../../components/follow-button/follow-button.component';

@Component({
  selector: 'app-people',
  standalone: true,
  imports: [CommonModule, FollowButtonComponent],
  templateUrl: './people.component.html',
  styleUrls: ['./people.component.scss']
})
export class PeopleComponent implements OnInit {
  allUsers: User[] = [];
  recommendations: UserRecommendation[] = [];
  followers: User[] = [];
  following: User[] = [];
  
  currentTab: 'discover' | 'following' | 'followers' = 'discover';
  isLoading: boolean = false;
  currentUserId: string = '';

  constructor(
    private followerService: FollowerService,
    private authService: AuthService,
    private http: HttpClient
  ) {}

  ngOnInit() {
    const currentUser = this.authService.currentUser.getValue();
    if (currentUser) {
      this.currentUserId = currentUser.id;
      this.loadData();
    }
  }

  loadData() {
    this.loadAllUsers();
    this.loadRecommendations();
    this.loadFollowers();
    this.loadFollowing();
  }

  loadAllUsers() {
    this.isLoading = true;
    
    // Direktno koristi follower service umesto stakeholders API
    // jer stakeholders API zahteva admin prava
    this.loadUsersFromFollowerService();
  }

  loadUsersFromFollowerService() {
    this.followerService.getAllUsers().subscribe({
      next: (users: User[]) => {
        this.allUsers = users.filter(user => user.id !== this.currentUserId);
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Error loading users from follower service:', error);
        this.isLoading = false;
      }
    });
  }

  loadRecommendations() {
    this.followerService.getRecommendations().subscribe({
      next: (recommendations) => {
        this.recommendations = recommendations;
      },
      error: (error) => {
        console.error('Error loading recommendations:', error);
      }
    });
  }

  loadFollowers() {
    this.followerService.getFollowers(this.currentUserId).subscribe({
      next: (followers) => {
        this.followers = followers;
      },
      error: (error) => {
        console.error('Error loading followers:', error);
      }
    });
  }

  loadFollowing() {
    this.followerService.getFollowing(this.currentUserId).subscribe({
      next: (following) => {
        this.following = following;
      },
      error: (error) => {
        console.error('Error loading following:', error);
      }
    });
  }

  setTab(tab: 'discover' | 'following' | 'followers') {
    this.currentTab = tab;
    
    // Refresh data when tab changes to ensure fresh state
    if (tab === 'discover') {
      this.loadAllUsers();
      this.loadRecommendations();
    } else if (tab === 'following') {
      this.loadFollowing();
    } else if (tab === 'followers') {
      this.loadFollowers();
    }
  }

  onFollowStatusChanged(userId: string, isFollowing: boolean) {
    // Only refresh the current tab's data to avoid excessive API calls
    if (this.currentTab === 'following') {
      this.loadFollowing();
    } else if (this.currentTab === 'followers') {
      this.loadFollowers();
    } else if (this.currentTab === 'discover') {
      // For discover tab, just refresh recommendations
      this.loadRecommendations();
    }
    
    // Optional: refresh all data but with a delay to avoid conflicts
    // setTimeout(() => {
    //   this.loadData();
    // }, 500);
  }

  createUserInFollowerService(user: User) {
    this.followerService.createUser(user).subscribe({
      next: () => {
        console.log('User created in follower service');
      },
      error: (error) => {
        console.error('Error creating user in follower service:', error);
      }
    });
  }

  // TrackBy functions for better performance
  trackByUserId(index: number, user: any): string {
    return user.id || user.username || index.toString();
  }

  trackByRecommendationId(index: number, rec: any): string {
    return rec.user?.id || rec.user?.username || index.toString();
  }
}
