import { Component, OnInit, OnDestroy } from '@angular/core';
import { BlogService } from '../../services/blog.service';
import { FollowerService } from '../../services/follower.service';
import { AuthService } from '../../services/auth.service';
import { Blog } from '../../models/blog.model';
import { Observable, combineLatest, map, of, Subscription } from 'rxjs';
import { BlogReloadService } from '../../services/blog-reload.service';
import { Router } from '@angular/router';
// ======== 1. UVEZI CommonModule ========
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-blog-list',
  standalone: true,
  // ======== 2. DODAJ CommonModule U IMPORTS NIZ ========
  imports: [CommonModule],
  templateUrl: './blog-list.component.html',
  styleUrls: ['./blog-list.component.scss']
})
export class BlogListComponent implements OnInit, OnDestroy {
  blogs$!: Observable<Blog[]>;
  showFollowedOnly: boolean = true;
  private reloadSub?: Subscription;

  constructor(
    private blogService: BlogService, 
    private followerService: FollowerService,
    private authService: AuthService,
    private router: Router,
    private blogReloadService: BlogReloadService
  ) {}

  ngOnInit(): void {
    this.loadBlogs();
    this.reloadSub = this.blogReloadService.reload$.subscribe(() => this.loadBlogs());
  }

  ngOnDestroy(): void {
    this.reloadSub?.unsubscribe();
  }

  loadBlogs(): void {
    if (this.showFollowedOnly) {
      // Koristi novi endpoint za blogove od praćenih korisnika
      this.blogs$ = this.blogService.getBlogsFromFollowing();
    } else {
      // Učitaj sve blogove
      this.blogs$ = this.blogService.getAllBlogs();
    }
  }

  toggleFilter(): void {
    this.showFollowedOnly = !this.showFollowedOnly;
    this.loadBlogs();
  }

  viewBlog(id: string): void {
    this.router.navigate(['/blogs', id]);
  }
  
  createBlog(): void {
    this.router.navigate(['/blogs/create']);
  }
}