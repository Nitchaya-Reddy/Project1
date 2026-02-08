import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { environment } from '../../environments/environment';
import { Chat, Message } from '../models/chat.model';

@Injectable({
  providedIn: 'root'
})
export class ChatService {
  private apiUrl = environment.apiUrl;
  
  chats = signal<Chat[]>([]);
  activeChat = signal<Chat | null>(null);
  messages = signal<Message[]>([]);

  constructor(private http: HttpClient) {}

  getChats(): Observable<Chat[]> {
    return this.http.get<Chat[]>(`${this.apiUrl}/chats`).pipe(
      tap(chats => this.chats.set(chats || []))
    );
  }

  getChat(id: number): Observable<Chat> {
    return this.http.get<Chat>(`${this.apiUrl}/chats/${id}`).pipe(
      tap(chat => this.activeChat.set(chat))
    );
  }

  getChatMessages(chatId: number): Observable<Message[]> {
    return this.http.get<Message[]>(`${this.apiUrl}/chats/${chatId}/messages`).pipe(
      tap(messages => this.messages.set(messages))
    );
  }

  createChat(listingId: number, message: string): Observable<{ chat_id: number; message: Message }> {
    return this.http.post<{ chat_id: number; message: Message }>(`${this.apiUrl}/chats`, {
      listing_id: listingId,
      message: message
    }).pipe(
      tap(() => this.refreshChats())
    );
  }

  refreshChats(): void {
    this.getChats().subscribe();
  }

  sendMessage(chatId: number, content: string): Observable<Message> {
    return this.http.post<Message>(`${this.apiUrl}/chats/${chatId}/messages`, { content }).pipe(
      tap(message => {
        this.messages.update(messages => [...messages, message]);
      })
    );
  }

  getTotalUnreadCount(): number {
    return this.chats().reduce((total, chat) => total + chat.unread_count, 0);
  }
}
