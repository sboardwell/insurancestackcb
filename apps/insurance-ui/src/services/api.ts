// API service using Axios for InsuranceStack
import axios, { AxiosInstance, AxiosError } from 'axios';
import type {
  User,
  Customer,
  Policy,
  Claim,
  Payment,
  Quote,
} from '../types';
import { isDebugModeEnabled } from '../features/flags';

// Create axios instance with base configuration
const apiClient: AxiosInstance = axios.create({
  baseURL: 'api',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for adding auth tokens, logging, etc.
apiClient.interceptors.request.use(
  (config) => {
    // Add auth token if available
    const token = localStorage.getItem('authToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Debug logging for requests
    if (isDebugModeEnabled()) {
      console.log('[DEBUG] API Request:', {
        method: config.method?.toUpperCase(),
        url: config.url,
        params: config.params,
        data: config.data,
      });
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => {
    // Debug logging for successful responses
    if (isDebugModeEnabled()) {
      console.log('[DEBUG] API Response:', {
        status: response.status,
        url: response.config.url,
        data: response.data,
      });
    }
    return response;
  },
  (error: AxiosError) => {
    if (error.response) {
      // Server responded with error status
      console.error('[API Error]', {
        status: error.response.status,
        data: error.response.data,
        url: error.config?.url,
      });
    } else if (error.request) {
      // Request made but no response
      console.error('[Network Error]', error.message);
    } else {
      // Error in request setup
      console.error('[Request Error]', error.message);
    }
    return Promise.reject(error);
  }
);

// API Service class
class ApiService {
  // User endpoints
  async getCurrentUser(): Promise<User> {
    const response = await apiClient.get<User>('customers/me');
    return response.data;
  }

  // Customer endpoints
  async getCustomers(): Promise<Customer[]> {
    const response = await apiClient.get<Customer[]>('customers');
    return response.data;
  }

  async getCustomer(customerId: string): Promise<Customer> {
    const response = await apiClient.get<Customer>(`customers/${customerId}`);
    return response.data;
  }

  async createCustomer(customerData: Partial<Customer>): Promise<Customer> {
    const response = await apiClient.post<Customer>('customers', customerData);
    return response.data;
  }

  async updateCustomer(customerId: string, customerData: Partial<Customer>): Promise<Customer> {
    const response = await apiClient.put<Customer>(
      `customers/${customerId}`,
      customerData
    );
    return response.data;
  }

  async deleteCustomer(customerId: string): Promise<void> {
    await apiClient.delete(`customers/${customerId}`);
  }

  // Policy endpoints
  async getPolicies(params?: {
    customerId?: string;
    policyType?: string;
    status?: string;
  }): Promise<Policy[]> {
    const response = await apiClient.get<Policy[]>('policies', {
      params,
    });
    return response.data;
  }

  async getPolicy(policyId: string): Promise<Policy> {
    const response = await apiClient.get<Policy>(`policies/${policyId}`);
    return response.data;
  }

  async createPolicy(policyData: Partial<Policy>): Promise<Policy> {
    const response = await apiClient.post<Policy>('policies', policyData);
    return response.data;
  }

  async updatePolicy(policyId: string, policyData: Partial<Policy>): Promise<Policy> {
    const response = await apiClient.put<Policy>(
      `policies/${policyId}`,
      policyData
    );
    return response.data;
  }

  async deletePolicy(policyId: string): Promise<void> {
    await apiClient.delete(`policies/${policyId}`);
  }

  // Claim endpoints
  async getClaims(params?: {
    policyId?: string;
    status?: string;
    startDate?: string;
    endDate?: string;
  }): Promise<Claim[]> {
    const response = await apiClient.get<Claim[]>('claims', {
      params,
    });
    return response.data;
  }

  async getClaim(claimId: string): Promise<Claim> {
    const response = await apiClient.get<Claim>(`claims/${claimId}`);
    return response.data;
  }

  async createClaim(claimData: Partial<Claim>): Promise<Claim> {
    const response = await apiClient.post<Claim>('claims', claimData);
    return response.data;
  }

  async updateClaim(claimId: string, claimData: Partial<Claim>): Promise<Claim> {
    const response = await apiClient.put<Claim>(
      `claims/${claimId}`,
      claimData
    );
    return response.data;
  }

  // Payment endpoints
  async getPayments(params?: {
    policyId?: string;
    paymentType?: string;
    status?: string;
  }): Promise<Payment[]> {
    const response = await apiClient.get<Payment[]>('payments', {
      params,
    });
    return response.data;
  }

  async getPayment(paymentId: string): Promise<Payment> {
    const response = await apiClient.get<Payment>(`payments/${paymentId}`);
    return response.data;
  }

  async createPayment(paymentData: Partial<Payment>): Promise<Payment> {
    const response = await apiClient.post<Payment>('payments', paymentData);
    return response.data;
  }

  // Quote endpoints
  async getQuotes(params?: {
    customerId?: string;
    policyType?: string;
    status?: string;
  }): Promise<Quote[]> {
    const response = await apiClient.get<Quote[]>('quotes', {
      params,
    });
    return response.data;
  }

  async getQuote(quoteId: string): Promise<Quote> {
    const response = await apiClient.get<Quote>(`quotes/${quoteId}`);
    return response.data;
  }

  async createQuote(quoteData: Partial<Quote>): Promise<Quote> {
    const response = await apiClient.post<Quote>('quotes', quoteData);
    return response.data;
  }

  async acceptQuote(quoteId: string): Promise<Policy> {
    const response = await apiClient.post<Policy>(`quotes/${quoteId}/accept`);
    return response.data;
  }
}

// Export singleton instance
export const api = new ApiService();
export default apiClient;
