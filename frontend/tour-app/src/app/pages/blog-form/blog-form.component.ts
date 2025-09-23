// src/app/pages/blog-form/blog-form.component.ts

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { BlogService } from '../../services/blog.service';
import { switchMap } from 'rxjs/operators';
import { of } from 'rxjs';
import { Blog } from '../../models/blog.model';

// Potrebno je importovati FormsModule da bi [(ngModel)] radio u standalone komponenti
import { FormsModule } from '@angular/forms'; 
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-blog-form',
  standalone: true, // Dodajemo jer je verovatno standalone
  imports: [CommonModule, FormsModule], // Importujemo CommonModule i FormsModule
  templateUrl: './blog-form.component.html',
  styleUrls: ['./blog-form.component.scss']
})
export class BlogFormComponent implements OnInit {
  // Umesto FormGroup, sada imamo objekat koji direktno vezujemo za formu
  blogData: Partial<Blog> = { title: '', content: '', images: [] };
  isEditMode = false;
  blogId: string | null = null;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private blogService: BlogService
  ) {}

  ngOnInit(): void {
    this.route.paramMap.pipe(
      switchMap(params => {
        this.blogId = params.get('id');
        if (this.blogId) {
          this.isEditMode = true;
          return this.blogService.getBlogById(this.blogId);
        }
        return of(null);
      })
    ).subscribe(blog => {
      if (blog) {
        // U edit modu, jednostavno postavimo podatke bloga
        this.blogData = blog;
      }
    });
  }

  // Nova metoda za selekciju VIŠE fajlova
  onFilesSelected(event: any): void {
    const files: FileList = event.target.files;
    if (files && files.length > 0) {
      // Počinjemo sa praznim nizom da ne bi dodavali duplikate
      this.blogData.images = [];
      for (let i = 0; i < files.length; i++) {
        const file = files[i];
        // Pravimo putanju kao i za profilnu sliku
        const imagePath = 'assets/images/' + file.name;
        this.blogData.images.push(imagePath);
      }
    }
  }

  onSubmit(): void {
    // Proveravamo da li su osnovna polja popunjena
    if (!this.blogData.title || !this.blogData.content) {
      return;
    }

    if (this.isEditMode && this.blogId) {
      this.blogService.updateBlog(this.blogId, this.blogData).subscribe(() => {
        this.router.navigate(['/blogs', this.blogId]);
      });
    } else {
      // Pravimo payload sa title, content, images
      const payload = {
        title: this.blogData.title || '',
        content: this.blogData.content || '',
        images: this.blogData.images || []
      };
      this.blogService.createBlog(payload).subscribe(() => {
        this.router.navigate(['/blogs']); 
      });
    }
  }
}