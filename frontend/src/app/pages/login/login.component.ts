import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule, ActivatedRoute } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss'
})
export class LoginComponent {
  private authService = inject(AuthService);
  private router = inject(Router);
  private route = inject(ActivatedRoute);

  email = '';
  password = '';
  error = '';
  isLoading = false;

  login(): void {
    this.error = '';

    if (!this.email.trim()) {
      this.error = 'Email is required';
      return;
    }

    if (!this.password) {
      this.error = 'Password is required';
      return;
    }

    this.isLoading = true;

    this.authService.login({ email: this.email.trim().toLowerCase(), password: this.password }).subscribe({
      next: () => {
        const returnUrl = this.route.snapshot.queryParams['returnUrl'] || '/';
        this.router.navigateByUrl(returnUrl);
      },
      error: (err) => {
        this.isLoading = false;
        const errorMessage = err.error?.error;
        if (errorMessage) {
          this.error = errorMessage;
        } else if (err.status === 0) {
          this.error = 'Unable to connect to server. Please try again.';
        } else if (err.status === 401) {
          this.error = 'Invalid email or password.';
        } else {
          this.error = 'Login failed. Please try again.';
        }
      }
    });
  }
}
