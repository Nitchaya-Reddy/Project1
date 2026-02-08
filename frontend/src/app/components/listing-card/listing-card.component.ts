import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { Listing } from '../../models/listing.model';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-listing-card',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './listing-card.component.html',
  styleUrl: './listing-card.component.scss'
})
export class ListingCardComponent {
  @Input({ required: true }) listing!: Listing;

  private getImageUrl(url: string): string {
    if (url.startsWith('/uploads/')) {
      return environment.apiUrl.replace('/api', '') + url;
    }
    return url;
  }

  get primaryImage(): string {
    const primary = this.listing.images?.find(img => img.is_primary);
    if (primary) return this.getImageUrl(primary.image_url);
    if (this.listing.images?.length > 0) return this.getImageUrl(this.listing.images[0].image_url);
    return 'assets/placeholder.svg';
  }

  get sellerInitials(): string {
    return `${this.listing.seller?.first_name?.charAt(0) || ''}${this.listing.seller?.last_name?.charAt(0) || ''}`;
  }
}
