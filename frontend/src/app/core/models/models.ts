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

export interface TicketAssignmentHistory {
  id: number;
  ticketId: string;
  assignedToUserId?: string;
  assignedByUserId?: string;
  assignedAt: string;
  endedAt?: string; 
  statusAtAssignment: 'open' | 'in_progress' | 'resolved' | 'closed';

  assignedToUser?: User;
  assignedByUser?: User;
}

export interface Ticket {
  id: string;
  title: string;
  description: string;
  status: 'open' | 'in_progress' | 'resolved' | 'closed';
  priority: 'low' | 'mid' | 'high' | 'critical';
  categoryId?: number;
  createdBy?: string;
  assignedTo?: string;
  createdAt: string;
  updatedAt: string;
  resolvedAt?: string;

  category?: Category;
  creator?: User;
  assignee?: User;

  // ➕ AJOUTER CETTE LIGNE
  assignmentHistories?: TicketAssignmentHistory[];
}

export interface RefreshToken {
  id: number;
  userId: string;
  revoked: boolean;
  expiresAt: string;
  createdAt: string;

  user?: User;
}
