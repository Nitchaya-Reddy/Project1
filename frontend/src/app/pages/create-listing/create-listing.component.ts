import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { ListingService } from '../../services/listing.service';
import { Category } from '../../models/listing.model';

@Component({
  selector: 'app-create-listing',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './create-listing.component.html',
  styleUrl: './create-listing.component.scss'
})
export class CreateListingComponent implements OnInit {
  private listingService = inject(ListingService);
  private router = inject(Router);

  categories = signal<Category[]>([]);
  isLoading = signal(false);
  error = '';

  // Form data
  title = '';
  description = '';
  price: number | null = null;
  categoryId: number | null = null;
  condition = '';
  location = '';
  images: string[] = [];
  uploadingImage = false;

  conditions = [
    { value: 'new', label: 'New' },
    { value: 'like_new', label: 'Like New' },
    { value: 'good', label: 'Good' },
    { value: 'fair', label: 'Fair' },
    { value: 'poor', label: 'Poor' }
  ];

  ngOnInit(): void {
    this.loadCategories();
  }

  loadCategories(): void {
    this.listingService.getCategories().subscribe({
      next: (categories) => this.categories.set(categories),
      error: (err) => console.error('Error loading categories:', err)
    });
  }

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) return;

    const file = input.files[0];
    if (!file.type.startsWith('image/')) {
      this.error = 'Please select an image file';
      return;
    }

    this.uploadingImage = true;
    this.listingService.uploadImage(file).subscribe({
      next: (response) => {
        this.images.push(response.url);
        this.uploadingImage = false;
      },
      error: (err) => {
        console.error('Error uploading image:', err);
        this.error = 'Error uploading image';
        this.uploadingImage = false;
      }
    });

    input.value = '';
  }

  removeImage(index: number): void {
    this.images.splice(index, 1);
  }

  submit(): void {
    this.error = '';

    if (!this.title) {
      this.error = 'Please enter a title';
      return;
    }
    if (!this.price || this.price <= 0) {
      this.error = 'Please enter a valid price';
      return;
    }
    if (!this.categoryId) {
      this.error = 'Please select a category';
      return;
    }

    this.isLoading.set(true);

    this.listingService.createListing({
      title: this.title,
      description: this.description,
      price: this.price,
      category_id: this.categoryId,
      condition: this.condition,
      location: this.location,
      images: this.images
    }).subscribe({
      next: (listing) => {
        this.router.navigate(['/listing', listing.id]);
      },
      error: (err) => {
        console.error('Error creating listing:', err);
        this.error = err.error?.error || 'Error creating listing';
        this.isLoading.set(false);
      }
    });
  }
}
