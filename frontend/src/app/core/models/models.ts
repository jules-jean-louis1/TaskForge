// ==========================================
// Enums & Types
// ==========================================

export type Role = 'admin' | 'tech' | 'standard';

export type Status = 'open' | 'in_progress' | 'resolved' | 'closed';

export type Priority = 'low' | 'mid' | 'high' | 'critical';

// ==========================================
// Entities / Models
// ==========================================

export interface User {
  id: string;
  firstname: string;
  lastname: string;
  email: string;
  role: Role;
  createdAt: string;
  updatedAt: string;
}

export interface Category {
  id: number;
  name: string;
}

export interface Ticket {
  id: string;
  title: string;
  description: string;
  status: Status;
  priority: Priority;
  categoryId?: number;
  createdBy?: string;
  assignedTo?: string;
  createdAt: string;
  updatedAt: string;
  resolvedAt?: string;

  category?: Category;
  creator?: User;
  assignee?: User;
}

export interface TicketAssignmentHistory {
  id: number;
  ticketId: string;
  assignedToUserId?: string;
  assignedByUserId?: string;
  assignedAt: string;
  endedAt?: string;
  statusAtAssignment: Status;

  ticket?: Ticket;
  assignedToUser?: User;
  assignedByUser?: User;
}

export interface RefreshToken {
  id: number;
  userId: string;
  revoked: boolean;
  expiresAt: string;
  createdAt: string;

  user?: User;
}