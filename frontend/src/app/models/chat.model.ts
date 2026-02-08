import { User } from './user.model';
import { Listing } from './listing.model';

export interface Message {
  id: number;
  chat_id: number;
  sender_id: number;
  sender: User;
  content: string;
  is_read: boolean;
  read_at?: string;
  created_at: string;
}

export interface Chat {
  id: number;
  listing_id: number;
  listing: Listing;
  buyer_id: number;
  buyer: User;
  seller_id: number;
  seller: User;
  last_message?: Message;
  unread_count: number;
  created_at: string;
  updated_at: string;
}

export interface Notification {
  id: number;
  user_id: number;
  type: 'new_message' | 'new_offer' | 'listing_sold' | 'price_dropped';
  title: string;
  message: string;
  link: string;
  is_read: boolean;
  read_at?: string;
  created_at: string;
}
