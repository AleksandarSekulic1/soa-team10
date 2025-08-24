import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { BlogService } from '../../services/blog.service';
import { Blog, Comment } from '../../models/blog.model';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { AuthService } from '../../services/auth.service';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms'; // <-- Potrebno za [(ngModel)]
import { RouterModule } from '@angular/router'; // <-- 1. Uvezite RouterModule


@Component({
  selector: 'app-blog-detail',
  // Potrebno je da bude standalone ako i ostale komponente jesu
  standalone: true,
  // U imports dodajemo CommonModule za async pipe i *ngIf/*ngFor
  imports: [CommonModule,FormsModule,RouterModule],
  templateUrl: './blog-detail.component.html',
  styleUrls: ['./blog-detail.component.scss']
})
export class BlogDetailComponent implements OnInit {
  blog$!: Observable<Blog>;
  blogId!: string;
  isAuthor = false;
  currentUser!: string;

  // Svojstva za praćenje editovanja komentara
  editingCommentId: string | null = null;
  editingCommentText = '';

  constructor(
    private route: ActivatedRoute,
    private blogService: BlogService,
    private authService: AuthService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.blogId = this.route.snapshot.paramMap.get('id')!;
    this.currentUser = this.authService.getUsername()!;
    this.loadBlog();
  }

  loadBlog(): void {
    this.blog$ = this.blogService.getBlogById(this.blogId).pipe(
      tap(blog => {
        this.isAuthor = this.currentUser === blog.AuthorID;
      })
    );
  }

  toggleLike(): void {
    this.blogService.toggleLike(this.blogId).subscribe(() => {
      this.loadBlog();
    });
  }
  
  addComment(commentInput: HTMLTextAreaElement): void {
    const text = commentInput.value;
    if (!text) return;
    
    // ISPRAVKA: Šaljemo objekat sa velikim slovom 'Text'
    this.blogService.addComment(this.blogId, { Text: text }).subscribe(() => {
      commentInput.value = '';
      this.loadBlog();
    });
  }

  editBlog(): void {
    this.router.navigate(['/blogs/edit', this.blogId]);
  }

  startEditComment(comment: Comment): void {
    this.editingCommentId = comment.ID!;
    this.editingCommentText = comment.Text;
  }

  cancelEdit(): void {
    this.editingCommentId = null;
    this.editingCommentText = '';
  }

  saveComment(comment: Comment): void {
    if (!this.editingCommentText.trim()) return;

    const commentData = { Text: this.editingCommentText };
    
    this.blogService.updateComment(this.blogId, comment.ID!, commentData).subscribe(() => {
      this.cancelEdit();
      this.loadBlog();
    });
  }
}