import { Component, OnInit, Inject, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser, CommonModule } from '@angular/common';
import { UserService } from '../../services/user.service';
import { User } from '../../models/user.model';

@Component({
  selector: 'app-user-list',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.scss']
})
export class UserListComponent implements OnInit {
  users: User[] = [];

  constructor(
    private userService: UserService,
    @Inject(PLATFORM_ID) private platformId: Object
  ) { }

  ngOnInit(): void {
    if (isPlatformBrowser(this.platformId)) {
      this.loadUsers();
    }
  }

  loadUsers(): void {
    this.userService.getAllUsers().subscribe({
      next: (data) => {
        this.users = data;
      },
      error: (err) => {
        console.error('Error fetching users', err);
      }
    });
  }

  toggleBlockStatus(user: User): void {
    const action = user.isBlocked
      ? this.userService.unblockUser(user.username)
      : this.userService.blockUser(user.username);

    action.subscribe({
      next: () => {
        user.isBlocked = !user.isBlocked;
      },
      error: (err) => console.error('Error changing user block status', err)
    });
  }

  // NOVA METODA KOJA ĆE BITI DOSTUPNA U HTML TEMPLEJTU
  public getTypeOf(value: any): string {
    return typeof value;
  }
}