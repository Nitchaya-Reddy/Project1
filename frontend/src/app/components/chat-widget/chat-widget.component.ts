import { Component, inject, OnInit, signal, ElementRef, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { ChatService } from '../../services/chat.service';
import { AuthService } from '../../services/auth.service';
import { Chat, Message } from '../../models/chat.model';

@Component({
  selector: 'app-chat-widget',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './chat-widget.component.html',
  styleUrl: './chat-widget.component.scss'
})
export class ChatWidgetComponent implements OnInit {
  @ViewChild('messagesContainer') messagesContainer!: ElementRef;
  
  chatService = inject(ChatService);
  authService = inject(AuthService);
  
  isOpen = signal(false);
  activeChat = signal<Chat | null>(null);
  newMessage = '';
  isLoading = signal(false);

  ngOnInit(): void {
    if (this.authService.isLoggedIn()) {
      this.loadChats();
    }
  }

  loadChats(): void {
    this.chatService.getChats().subscribe();
  }

  toggle(): void {
    this.isOpen.update(v => !v);
  }

  close(): void {
    this.isOpen.set(false);
    this.activeChat.set(null);
  }

  openChat(chat: Chat): void {
    this.activeChat.set(chat);
    this.chatService.getChatMessages(chat.id).subscribe(() => {
      setTimeout(() => this.scrollToBottom(), 100);
    });
  }

  backToList(): void {
    this.activeChat.set(null);
    this.loadChats();
  }

  sendMessage(): void {
    if (!this.newMessage.trim() || !this.activeChat()) return;
    
    const chat = this.activeChat();
    if (!chat) return;

    this.isLoading.set(true);
    this.chatService.sendMessage(chat.id, this.newMessage).subscribe({
      next: () => {
        this.newMessage = '';
        this.isLoading.set(false);
        setTimeout(() => this.scrollToBottom(), 100);
      },
      error: () => {
        this.isLoading.set(false);
      }
    });
  }

  scrollToBottom(): void {
    if (this.messagesContainer) {
      const el = this.messagesContainer.nativeElement;
      el.scrollTop = el.scrollHeight;
    }
  }

  getOtherUser(chat: Chat): any {
    const currentUserId = this.authService.currentUser()?.id;
    return chat.buyer_id === currentUserId ? chat.seller : chat.buyer;
  }

  isOwnMessage(message: Message): boolean {
    return message.sender_id === this.authService.currentUser()?.id;
  }

  getTotalUnread(): number {
    return this.chatService.getTotalUnreadCount();
  }
}
