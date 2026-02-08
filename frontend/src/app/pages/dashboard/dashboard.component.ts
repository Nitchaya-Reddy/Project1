import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { ListingService, ListingFilters } from '../../services/listing.service';
import { ListingCardComponent } from '../../components/listing-card/listing-card.component';
import { Listing, Category } from '../../models/listing.model';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule, ListingCardComponent],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss'
})
export class DashboardComponent implements OnInit {
  private listingService = inject(ListingService);

  listings = signal<Listing[]>([]);
  categories = signal<Category[]>([]);
  isLoading = signal(true);
  
  // Filters
  searchQuery = '';
  selectedCategory: number | null = null;
  minPrice: number | null = null;
  maxPrice: number | null = null;
  selectedCondition = '';
  sortBy = 'created_at';
  sortOrder = 'desc';

  // Pagination
  currentPage = 1;
  totalPages = 1;
  totalItems = 0;

  conditions = [
    { value: '', label: 'All Conditions' },
    { value: 'new', label: 'New' },
    { value: 'like_new', label: 'Like New' },
    { value: 'good', label: 'Good' },
    { value: 'fair', label: 'Fair' },
    { value: 'poor', label: 'Poor' }
  ];

  ngOnInit(): void {
    this.loadCategories();
    this.loadListings();
  }

  loadCategories(): void {
    this.listingService.getCategories().subscribe({
      next: (categories) => this.categories.set(categories),
      error: (err) => console.error('Error loading categories:', err)
    });
  }

  loadListings(): void {
    this.isLoading.set(true);
    
    const filters: ListingFilters = {
      page: this.currentPage,
      limit: 20,
      sort: this.sortBy,
      order: this.sortOrder
    };

    if (this.searchQuery) filters.search = this.searchQuery;
    if (this.selectedCategory) filters.category_id = this.selectedCategory;
    if (this.minPrice) filters.min_price = this.minPrice;
    if (this.maxPrice) filters.max_price = this.maxPrice;
    if (this.selectedCondition) filters.condition = this.selectedCondition;

    this.listingService.getListings(filters).subscribe({
      next: (response) => {
        this.listings.set(response.listings || []);
        this.totalPages = response.pages;
        this.totalItems = response.total;
        this.isLoading.set(false);
      },
      error: (err) => {
        console.error('Error loading listings:', err);
        this.isLoading.set(false);
      }
    });
  }

  search(): void {
    this.currentPage = 1;
    this.loadListings();
  }

  selectCategory(categoryId: number | null): void {
    this.selectedCategory = categoryId;
    this.currentPage = 1;
    this.loadListings();
  }

  applyFilters(): void {
    this.currentPage = 1;
    this.loadListings();
  }

  clearFilters(): void {
    this.searchQuery = '';
    this.selectedCategory = null;
    this.minPrice = null;
    this.maxPrice = null;
    this.selectedCondition = '';
    this.sortBy = 'created_at';
    this.sortOrder = 'desc';
    this.currentPage = 1;
    this.loadListings();
  }

  changePage(page: number): void {
    if (page >= 1 && page <= this.totalPages) {
      this.currentPage = page;
      this.loadListings();
      window.scrollTo({ top: 0, behavior: 'smooth' });
    }
  }

  getPageNumbers(): number[] {
    const pages: number[] = [];
    const start = Math.max(1, this.currentPage - 2);
    const end = Math.min(this.totalPages, this.currentPage + 2);
    for (let i = start; i <= end; i++) {
      pages.push(i);
    }
    return pages;
  }
}
