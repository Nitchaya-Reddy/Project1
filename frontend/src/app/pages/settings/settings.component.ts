import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule, Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { User, getUserFullName } from '../../models/user.model';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './settings.component.html',
  styleUrl: './settings.component.scss'
})
export class SettingsComponent implements OnInit {
  private authService = inject(AuthService);
  private router = inject(Router);

  user = this.authService.user;
  
  formData = signal({
    name: '',
    phone: '',
    bio: ''
  });

  passwordData = signal({
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  });

  isLoading = signal(false);
  isSaving = signal(false);
  successMessage = signal('');
  errorMessage = signal('');
  activeSection = signal<'profile' | 'password' | 'account'>('profile');

  ngOnInit(): void {
    const user = this.user();
    if (user) {
      this.formData.set({
        name: getUserFullName(user),
        phone: user.phone || '',
        bio: user.bio || ''
      });
    }
  }

  setSection(section: 'profile' | 'password' | 'account'): void {
    this.activeSection.set(section);
    this.successMessage.set('');
    this.errorMessage.set('');
  }

  updateProfile(): void {
    this.isSaving.set(true);
    this.errorMessage.set('');
    this.successMessage.set('');

    this.authService.updateProfile(this.formData()).subscribe({
      next: (user) => {
        this.successMessage.set('Profile updated successfully!');
        this.isSaving.set(false);
      },
      error: (err) => {
        this.errorMessage.set(err.error?.error || 'Failed to update profile');
        this.isSaving.set(false);
      }
    });
  }

  changePassword(): void {
    const pwd = this.passwordData();
    
    if (pwd.newPassword !== pwd.confirmPassword) {
      this.errorMessage.set('Passwords do not match');
      return;
    }

    if (pwd.newPassword.length < 6) {
      this.errorMessage.set('Password must be at least 6 characters');
      return;
    }

    this.isSaving.set(true);
    this.errorMessage.set('');
    this.successMessage.set('');

    this.authService.changePassword(pwd.currentPassword, pwd.newPassword).subscribe({
      next: () => {
        this.successMessage.set('Password changed successfully!');
        this.passwordData.set({
          currentPassword: '',
          newPassword: '',
          confirmPassword: ''
        });
        this.isSaving.set(false);
      },
      error: (err) => {
        this.errorMessage.set(err.error?.error || 'Failed to change password');
        this.isSaving.set(false);
      }
    });
  }

  logout(): void {
    this.authService.logout();
    this.router.navigate(['/login']);
  }

  updateFormField(field: string, event: Event): void {
    const value = (event.target as HTMLInputElement).value;
    this.formData.update(data => ({ ...data, [field]: value }));
  }

  updatePasswordField(field: string, event: Event): void {
    const value = (event.target as HTMLInputElement).value;
    this.passwordData.update(data => ({ ...data, [field]: value }));
  }
}
