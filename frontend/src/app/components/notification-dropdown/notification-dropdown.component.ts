import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { NotificationService } from '../../services/notification.service';

@Component({
  selector: 'app-notification-dropdown',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './notification-dropdown.component.html',
  styleUrl: './notification-dropdown.component.scss'
})
export class NotificationDropdownComponent implements OnInit {
  notificationService = inject(NotificationService);
  isOpen = signal(false);

  ngOnInit(): void {
    this.loadNotifications();
  }

  loadNotifications(): void {
    this.notificationService.getNotifications().subscribe();
    this.notificationService.getUnreadCount().subscribe();
  }

  toggle(): void {
    this.isOpen.update(v => !v);
  }

  close(): void {
    this.isOpen.set(false);
  }

  markAsRead(id: number, event: Event): void {
    event.stopPropagation();
    this.notificationService.markAsRead(id).subscribe();
  }

  markAllAsRead(): void {
    this.notificationService.markAllAsRead().subscribe();
  }

  deleteNotification(id: number, event: Event): void {
    event.stopPropagation();
    event.preventDefault();
    this.notificationService.deleteNotification(id).subscribe();
  }

  getIcon(type: string): string {
    switch (type) {
      case 'new_message': return '💬';
      case 'new_offer': return '🏷️';
      case 'listing_sold': return '🎉';
      case 'price_dropped': return '📉';
      default: return '🔔';
    }
  }
}
