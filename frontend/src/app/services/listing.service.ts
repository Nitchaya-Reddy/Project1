import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import { Listing, ListingsResponse, Category, CreateListingRequest } from '../models/listing.model';

export interface ListingFilters {
  search?: string;
  category_id?: number;
  min_price?: number;
  max_price?: number;
  condition?: string;
  sort?: string;
  order?: string;
  page?: number;
  limit?: number;
}

@Injectable({
  providedIn: 'root'
})
export class ListingService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  getListings(filters: ListingFilters = {}): Observable<ListingsResponse> {
    let params = new HttpParams();
    
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params = params.set(key, value.toString());
      }
    });

    return this.http.get<ListingsResponse>(`${this.apiUrl}/listings`, { params });
  }

  getListing(id: number): Observable<Listing> {
    return this.http.get<Listing>(`${this.apiUrl}/listings/${id}`);
  }

  createListing(data: CreateListingRequest): Observable<Listing> {
    return this.http.post<Listing>(`${this.apiUrl}/listings`, data);
  }

  updateListing(id: number, data: Partial<CreateListingRequest & { status: string }>): Observable<Listing> {
    return this.http.put<Listing>(`${this.apiUrl}/listings/${id}`, data);
  }

  deleteListing(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/listings/${id}`);
  }

  getCategories(): Observable<Category[]> {
    return this.http.get<Category[]>(`${this.apiUrl}/categories`);
  }

  getMyListings(status?: string): Observable<Listing[]> {
    let params = new HttpParams();
    if (status) {
      params = params.set('status', status);
    }
    return this.http.get<Listing[]>(`${this.apiUrl}/users/me/listings`, { params });
  }

  getUserListings(userId: number): Observable<Listing[]> {
    return this.http.get<Listing[]>(`${this.apiUrl}/users/${userId}/listings`);
  }

  uploadImage(file: File): Observable<{ url: string; filename: string }> {
    const formData = new FormData();
    formData.append('image', file);
    return this.http.post<{ url: string; filename: string }>(`${this.apiUrl}/upload`, formData);
  }
}
