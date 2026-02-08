export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  profile_image: string;
  phone: string;
  bio: string;
  is_admin: boolean;
  created_at: string;
  CreatedAt?: string; // Alternative casing from backend
}

// Helper function to get full name
export function getUserFullName(user: User | null | undefined): string {
  if (!user) return '';
  return `${user.first_name || ''} ${user.last_name || ''}`.trim();
}

// Helper function to get initials
export function getUserInitials(user: User | null | undefined): string {
  if (!user) return 'U';
  const first = user.first_name?.charAt(0) || '';
  const last = user.last_name?.charAt(0) || '';
  return (first + last).toUpperCase() || 'U';
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
}
