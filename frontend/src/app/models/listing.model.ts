import { User } from './user.model';

export interface Category {
  id: number;
  name: string;
  description: string;
  icon: string;
}

export interface ListingImage {
  id: number;
  listing_id: number;
  image_url: string;
  is_primary: boolean;
}

export interface Listing {
  id: number;
  title: string;
  description: string;
  price: number;
  category_id: number;
  category: Category;
  seller_id: number;
  seller: User;
  images: ListingImage[];
  status: 'active' | 'sold' | 'inactive';
  condition: string;
  location: string;
  views: number;
  created_at: string;
  updated_at: string;
}

export interface ListingsResponse {
  listings: Listing[];
  total: number;
  page: number;
  limit: number;
  pages: number;
}

export interface CreateListingRequest {
  title: string;
  description: string;
  price: number;
  category_id: number;
  condition: string;
  location: string;
  images: string[];
}
