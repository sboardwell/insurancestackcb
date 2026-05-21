// Type definitions for InsuranceStack application

export interface User {
  id: string;
  email: string;
  name: string;
  createdAt: string;
}

export interface Customer {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  phone?: string;
  dateOfBirth?: string;
  address?: {
    street: string;
    city: string;
    state: string;
    zipCode: string;
    country: string;
  };
  riskScore?: number;
  createdAt: string;
  updatedAt: string;
}

export interface Policy {
  id: string;
  customerId: string;
  policyNumber: string;
  type: 'auto' | 'home' | 'life' | 'health';
  status: 'active' | 'lapsed' | 'cancelled';
  premium: number;
  coverage: number;
  deductible?: number;
  currency?: string;
  startDate: string;
  endDate: string;
  renewalDate?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Claim {
  id: string;
  policyId: string;
  customerId: string;
  claimNumber: string;
  type: string;
  status: 'submitted' | 'under_review' | 'approved' | 'rejected';
  amount: number;
  description: string;
  submittedDate: string;
  reviewedDate?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Payment {
  id: string;
  type: 'premium' | 'payout';
  policyId?: string;
  claimId?: string;
  customerId: string;
  amount: number;
  status: 'pending' | 'completed' | 'failed';
  paymentMethod?: string;
  processedDate?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Quote {
  id: string;
  customerId?: string;
  policyType: 'auto' | 'home' | 'life' | 'health';
  premium: number;
  coverage: number;
  deductible?: number;
  status: 'draft' | 'provided' | 'accepted' | 'expired';
  expiryDate?: string;
  createdAt: string;
}

export interface ApiResponse<T> {
  data: T;
  message?: string;
  timestamp: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
  hasMore: boolean;
}
