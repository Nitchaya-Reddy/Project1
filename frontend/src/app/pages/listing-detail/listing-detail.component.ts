import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { ListingService } from '../../services/listing.service';
import { ChatService } from '../../services/chat.service';
import { AuthService } from '../../services/auth.service';
import { Listing } from '../../models/listing.model';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-listing-detail',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './listing-detail.component.html',
  styleUrl: './listing-detail.component.scss'
})
export class ListingDetailComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private listingService = inject(ListingService);
  private chatService = inject(ChatService);
  authService = inject(AuthService);

  listing = signal<Listing | null>(null);
  isLoading = signal(true);
  selectedImageIndex = signal(0);
  
  messageText = '';
  isSendingMessage = false;
  showMessageForm = false;

  ngOnInit(): void {
    const id = Number(this.route.snapshot.paramMap.get('id'));
    if (id) {
      this.loadListing(id);
    }
  }

  loadListing(id: number): void {
    this.isLoading.set(true);
    this.listingService.getListing(id).subscribe({
      next: (listing) => {
        this.listing.set(listing);
        this.isLoading.set(false);
      },
      error: (err) => {
        console.error('Error loading listing:', err);
        this.isLoading.set(false);
        this.router.navigate(['/']);
      }
    });
  }

  selectImage(index: number): void {
    this.selectedImageIndex.set(index);
  }

  getImageUrl(url: string): string {
    if (url.startsWith('/uploads/')) {
      return environment.apiUrl.replace('/api', '') + url;
    }
    return url;
  }

  get currentImage(): string {
    const images = this.listing()?.images || [];
    if (images.length === 0) return 'assets/placeholder.svg';
    const url = images[this.selectedImageIndex()]?.image_url || images[0]?.image_url;
    return this.getImageUrl(url);
  }

  get isOwner(): boolean {
    return this.listing()?.seller_id === this.authService.currentUser()?.id;
  }

  toggleMessageForm(): void {
    if (!this.authService.isLoggedIn()) {
      this.router.navigate(['/login'], { queryParams: { returnUrl: this.router.url } });
      return;
    }
    this.showMessageForm = !this.showMessageForm;
  }

  sendMessage(): void {
    if (!this.messageText.trim() || !this.listing()) return;

    this.isSendingMessage = true;
    this.chatService.createChat(this.listing()!.id, this.messageText).subscribe({
      next: () => {
        this.isSendingMessage = false;
        this.showMessageForm = false;
        this.messageText = '';
        alert('Message sent! Check your messages.');
      },
      error: (err) => {
        console.error('Error sending message:', err);
        this.isSendingMessage = false;
        alert('Error sending message');
      }
    });
  }

  getConditionLabel(condition: string): string {
    const labels: { [key: string]: string } = {
      'new': 'New',
      'like_new': 'Like New',
      'good': 'Good',
      'fair': 'Fair',
      'poor': 'Poor'
    };
    return labels[condition] || condition;
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }
}
