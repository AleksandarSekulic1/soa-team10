import { Component, Input, Output, EventEmitter, OnInit, OnChanges, SimpleChanges } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FollowerService } from '../../services/follower.service';

@Component({
  selector: 'app-follow-button',
  standalone: true,
  imports: [CommonModule],
  template: `
    <button 
      [class]="buttonClass"
      [disabled]="isLoading"
      (click)="toggleFollow()"
    >
      <span *ngIf="isLoading">...</span>
      <span *ngIf="!isLoading">{{ buttonText }}</span>
    </button>
  `,
  styles: [`
    .follow-btn {
      background-color: #007bff;
      color: white;
      border: none;
      padding: 8px 16px;
      border-radius: 4px;
      cursor: pointer;
      transition: background-color 0.3s;
    }
    
    .follow-btn:hover:not(:disabled) {
      background-color: #0056b3;
    }
    
    .unfollow-btn {
      background-color: #dc3545;
      color: white;
      border: none;
      padding: 8px 16px;
      border-radius: 4px;
      cursor: pointer;
      transition: background-color 0.3s;
    }
    
    .unfollow-btn:hover:not(:disabled) {
      background-color: #c82333;
    }
    
    button:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
  `]
})
export class FollowButtonComponent implements OnInit, OnChanges {
  @Input() userId!: string;
  @Input() isCurrentUser: boolean = false;
  @Output() followStatusChanged = new EventEmitter<boolean>();

  isFollowing: boolean = false;
  isLoading: boolean = false;
  private lastClickTime: number = 0;
  private readonly clickDebounceMs = 1000; // 1 second debounce

  constructor(private followerService: FollowerService) {}

  ngOnInit() {
    this.refreshFollowStatus();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['userId'] && !changes['userId'].firstChange) {
      this.refreshFollowStatus();
    }
  }

  refreshFollowStatus() {
    if (!this.isCurrentUser && this.userId && this.userId.trim() !== '') {
      this.checkFollowStatus();
    }
  }

  get buttonText(): string {
    if (this.isCurrentUser) return 'You';
    return this.isFollowing ? 'Unfollow' : 'Follow';
  }

  get buttonClass(): string {
    if (this.isCurrentUser) return 'follow-btn';
    return this.isFollowing ? 'unfollow-btn' : 'follow-btn';
  }

  checkFollowStatus() {
    if (!this.userId || this.userId.trim() === '') {
      console.warn('Cannot check follow status: userId is empty');
      return;
    }
    
    this.isLoading = true;
    this.followerService.isFollowing(this.userId).subscribe({
      next: (response) => {
        const wasFollowing = this.isFollowing;
        this.isFollowing = response.isFollowing;
        this.isLoading = false;
        
        // Only emit if status actually changed
        if (wasFollowing !== this.isFollowing) {
          this.followStatusChanged.emit(this.isFollowing);
        }
      },
      error: (error) => {
        console.error('Error checking follow status:', error);
        this.isLoading = false;
      }
    });
  }

  toggleFollow() {
    const now = Date.now();
    if (now - this.lastClickTime < this.clickDebounceMs) {
      console.log('Click ignored due to debounce');
      return;
    }
    this.lastClickTime = now;

    if (this.isCurrentUser || this.isLoading || !this.userId || this.userId.trim() === '') {
      console.warn('Cannot toggle follow: invalid userId or user is current user');
      return;
    }

    // Prevent multiple clicks
    if (this.isLoading) {
      return;
    }

    this.isLoading = true;
    const previousState = this.isFollowing;
    
    if (this.isFollowing) {
      this.followerService.unfollowUser(this.userId).subscribe({
        next: () => {
          this.isFollowing = false;
          this.isLoading = false;
          this.followStatusChanged.emit(false);
        },
        error: (error) => {
          console.error('Error unfollowing user:', error);
          this.isFollowing = previousState; // Revert state on error
          this.isLoading = false;
        }
      });
    } else {
      this.followerService.followUser(this.userId).subscribe({
        next: () => {
          this.isFollowing = true;
          this.isLoading = false;
          this.followStatusChanged.emit(true);
        },
        error: (error) => {
          console.error('Error following user:', error);
          this.isFollowing = previousState; // Revert state on error
          this.isLoading = false;
        }
      });
    }
  }

  // Public method to refresh follow status from parent component
  public refreshStatus() {
    this.refreshFollowStatus();
  }
}
