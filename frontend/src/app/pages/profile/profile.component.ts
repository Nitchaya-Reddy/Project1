import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { ListingService } from '../../services/listing.service';
import { Listing } from '../../models/listing.model';
import { User, getUserFullName, getUserInitials } from '../../models/user.model';
import { ListingCardComponent } from '../../components/listing-card/listing-card.component';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, RouterModule, ListingCardComponent],
  templateUrl: './profile.component.html',
  styleUrl: './profile.component.scss'
})
export class ProfileComponent implements OnInit {
  private authService = inject(AuthService);
  private listingService = inject(ListingService);

  user = this.authService.user;
  activeListings = signal<Listing[]>([]);
  soldListings = signal<Listing[]>([]);
  isLoading = signal(true);
  activeTab = signal<'selling' | 'sold'>('selling');

  ngOnInit(): void {
    this.loadUserListings();
  }

  loadUserListings(): void {
    this.listingService.getMyListings().subscribe({
      next: (listings) => {
        this.activeListings.set(listings.filter(l => l.status === 'active'));
        this.soldListings.set(listings.filter(l => l.status === 'sold'));
        this.isLoading.set(false);
      },
      error: () => {
        this.isLoading.set(false);
      }
    });
  }

  setTab(tab: 'selling' | 'sold'): void {
    this.activeTab.set(tab);
  }

  getUserName(): string {
    return getUserFullName(this.user());
  }

  getUserInitials(): string {
    return getUserInitials(this.user());
  }

  getJoinedDate(): string {
    const user = this.user();
    const dateStr = user?.created_at || user?.CreatedAt;
    if (!dateStr) return '';
    return new Date(dateStr).toLocaleDateString('en-US', {
      month: 'long',
      year: 'numeric'
    });
  }

  getDisplayListings(): Listing[] {
    return this.activeTab() === 'selling' ? this.activeListings() : this.soldListings();
  }
}
