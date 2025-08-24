import { Component, OnInit } from '@angular/core';
import { UserService } from '../../services/user.service';
import { User } from '../../models/user.model';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss']
})
export class ProfileComponent implements OnInit {
  user: User | null = null;
  originalUser: User | null = null;
  isEditMode = false;
  successMessage = '';

  constructor(private userService: UserService) { }

  ngOnInit(): void {
    this.loadUserProfile();
  }

  loadUserProfile(): void {
    this.userService.getProfile().subscribe({
      next: (data) => {
        this.user = data;
        this.originalUser = JSON.parse(JSON.stringify(data));
      },
      error: (err) => console.error('Error loading profile', err)
    });
  }

  toggleEditMode(): void {
    this.isEditMode = true;
  }

  onFileSelected(event: any): void {
    const file: File = event.target.files?.[0];
    if (file && this.user) {
      this.user.profilePicture = 'assets/images/' + file.name;
    }
  }

  onCancel(): void {
    this.user = JSON.parse(JSON.stringify(this.originalUser));
    this.isEditMode = false;
  }

  onSave(): void {
    if (this.user) {
      // ISPRAVKA:
      // Više ne pravimo novi 'profileDataToUpdate' objekat.
      // Šaljemo ceo 'this.user' objekat, koji je već ažuriran
      // podacima iz forme zahvaljujući [(ngModel)].
      this.userService.updateProfile(this.user).subscribe({
        next: (updatedUser) => {
          this.user = updatedUser;
          this.originalUser = JSON.parse(JSON.stringify(updatedUser));
          this.successMessage = 'Profile successfully updated!';
          this.isEditMode = false;
          setTimeout(() => this.successMessage = '', 3000);
        },
        error: (err) => console.error('Error updating profile', err)
      });
    }
  }
}