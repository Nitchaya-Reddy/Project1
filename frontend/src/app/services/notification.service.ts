import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { environment } from '../../environments/environment';
import { Notification } from '../models/chat.model';

@Injectable({
  providedIn: 'root'
})
export class NotificationService {
  private apiUrl = environment.apiUrl;
  
  notifications = signal<Notification[]>([]);
  unreadCount = signal<number>(0);

  constructor(private http: HttpClient) {}

  getNotifications(unreadOnly: boolean = false): Observable<Notification[]> {
    const url = unreadOnly 
      ? `${this.apiUrl}/notifications?unread=true` 
      : `${this.apiUrl}/notifications`;
    return this.http.get<Notification[]>(url).pipe(
      tap(notifications => this.notifications.set(notifications))
    );
  }

  getUnreadCount(): Observable<{ count: number }> {
    return this.http.get<{ count: number }>(`${this.apiUrl}/notifications/unread-count`).pipe(
      tap(response => this.unreadCount.set(response.count))
    );
  }

  markAsRead(id: number): Observable<void> {
    return this.http.put<void>(`${this.apiUrl}/notifications/${id}/read`, {}).pipe(
      tap(() => {
        this.notifications.update(notifications => 
          notifications.map(n => n.id === id ? { ...n, is_read: true } : n)
        );
        this.unreadCount.update(count => Math.max(0, count - 1));
      })
    );
  }

  markAllAsRead(): Observable<void> {
    return this.http.put<void>(`${this.apiUrl}/notifications/read-all`, {}).pipe(
      tap(() => {
        this.notifications.update(notifications => 
          notifications.map(n => ({ ...n, is_read: true }))
        );
        this.unreadCount.set(0);
      })
    );
  }

  deleteNotification(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/notifications/${id}`).pipe(
      tap(() => {
        const notification = this.notifications().find(n => n.id === id);
        this.notifications.update(notifications => 
          notifications.filter(n => n.id !== id)
        );
        if (notification && !notification.is_read) {
          this.unreadCount.update(count => Math.max(0, count - 1));
        }
      })
    );
  }
}
