import { Component, OnInit } from '@angular/core';
import { BlogService } from '../../services/blog.service';
import { Blog } from '../../models/blog.model';
import { Observable } from 'rxjs';
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
export class BlogListComponent implements OnInit {
  blogs$!: Observable<Blog[]>;

  constructor(private blogService: BlogService, private router: Router) {}

  ngOnInit(): void {
    this.blogs$ = this.blogService.getAllBlogs();
  }

  viewBlog(id: string): void {
    this.router.navigate(['/blogs', id]);
  }
  
  createBlog(): void {
    this.router.navigate(['/blogs/create']);
  }
}