import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './register.component.html',
  styleUrl: './register.component.scss'
})
export class RegisterComponent {
  private authService = inject(AuthService);
  private router = inject(Router);

  firstName = '';
  lastName = '';
  email = '';
  password = '';
  confirmPassword = '';
  error = '';
  isLoading = false;

  register(): void {
    // Clear previous errors
    this.error = '';

    // Validate all fields are filled
    if (!this.firstName.trim()) {
      this.error = 'First name is required';
      return;
    }

    if (!this.lastName.trim()) {
      this.error = 'Last name is required';
      return;
    }

    if (!this.email.trim()) {
      this.error = 'Email is required';
      return;
    }

    // Validate UF email
    if (!this.email.toLowerCase().endsWith('@ufl.edu')) {
      this.error = 'Must use a valid UF email (@ufl.edu)';
      return;
    }

    if (!this.password) {
      this.error = 'Password is required';
      return;
    }


      if (this.password.length < 6) {
        this.error = 'Password must be at least 6 characters';
        return;
      }

      // Password complexity: at least one symbol, one uppercase, one lowercase, and one number
      const complexityRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z\d]).+$/;
      if (!complexityRegex.test(this.password)) {
        this.error = 'Password must contain at least one uppercase letter, one lowercase letter, one number, and one symbol.';
        return;
      }

    if (!this.confirmPassword) {
      this.error = 'Please confirm your password';
      return;
    }

    if (this.password !== this.confirmPassword) {
      this.error = 'Passwords do not match. Please re-enter your password.';
      return;
    }

    this.isLoading = true;

    this.authService.register({
      email: this.email.trim().toLowerCase(),
      password: this.password,
      first_name: this.firstName.trim(),
      last_name: this.lastName.trim()
    }).subscribe({
      next: () => {
        this.router.navigate(['/']);
      },
      error: (err) => {
        this.isLoading = false;
        // Handle specific error messages from backend
        const errorMessage = err.error?.error;
        if (errorMessage) {
          this.error = errorMessage;
        } else if (err.status === 0) {
          this.error = 'Unable to connect to server. Please try again.';
        } else if (err.status === 409) {
          this.error = 'An account with this email already exists.';
        } else {
          this.error = 'Registration failed. Please try again.';
        }
      }
    });
  }
}
