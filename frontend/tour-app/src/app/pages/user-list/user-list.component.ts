import { Component, OnInit } from '@angular/core';
import { UserService } from '../../services/user.service';
import { CommonModule } from '@angular/common'; // <-- DODAJEMO IMPORT

@Component({
  selector: 'app-user-list',
  standalone: true, // <-- ČINIMO GA SAMOSTALNIM
  imports: [CommonModule], // <-- DODAJEMO IMPORTE (ovo rešava *ngIf i *ngFor)
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.scss']
})
export class UserListComponent implements OnInit {
  users: any[] = [];

  constructor(private userService: UserService) { }

  ngOnInit(): void {
    this.userService.getAllUsers().subscribe({
      next: (data) => {
        this.users = data;
      },
      error: (err) => {
        console.error('Greška pri preuzimanju korisnika', err);
      }
    });
  }
}
