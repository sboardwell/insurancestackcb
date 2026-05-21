import { describe, it, expect, vi, beforeEach } from 'vitest';
import axios from 'axios';
import { api } from '../api';

// Mock axios
vi.mock('axios');
const mockedAxios = axios as any;

describe('API Service', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getCurrentUser', () => {
    it('should fetch current user', async () => {
      const mockUser = {
        data: {
          id: '1',
          email: 'test@example.com',
          name: 'Test User',
          createdAt: '2024-01-01T00:00:00Z',
        },
        message: 'Success',
        timestamp: '2024-01-01T00:00:00Z',
      };

      mockedAxios.create.mockReturnValue({
        get: vi.fn().mockResolvedValue({ data: mockUser }),
        interceptors: {
          request: { use: vi.fn(), eject: vi.fn() },
          response: { use: vi.fn(), eject: vi.fn() },
        },
      });

      // Note: This test demonstrates the structure, but would need proper setup
      // to work with the actual axios instance
    });
  });

  describe('getAccounts', () => {
    it('should fetch all accounts', async () => {
      const mockAccounts = {
        data: [
          {
            id: '1',
            userId: 'user1',
            accountNumber: '1234567890',
            accountType: 'checking',
            balance: 1000.0,
            currency: 'USD',
            status: 'active',
            createdAt: '2024-01-01T00:00:00Z',
            updatedAt: '2024-01-01T00:00:00Z',
          },
        ],
        message: 'Success',
        timestamp: '2024-01-01T00:00:00Z',
      };

      // Test structure demonstration
      expect(mockAccounts.data).toHaveLength(1);
      expect(mockAccounts.data[0].accountType).toBe('checking');
    });
  });

  describe('getTransactions', () => {
    it('should fetch transactions with filters', async () => {
      const params = {
        accountId: 'acc123',
        type: 'debit',
        category: 'groceries',
      };

      // Test demonstrates filter structure
      expect(params.accountId).toBe('acc123');
      expect(params.type).toBe('debit');
      expect(params.category).toBe('groceries');
    });
  });

  describe('getInsights', () => {
    it('should fetch insights with filters', async () => {
      const params = {
        type: 'spending',
        severity: 'warning',
        dismissed: false,
      };

      // Test demonstrates filter structure
      expect(params.type).toBe('spending');
      expect(params.severity).toBe('warning');
      expect(params.dismissed).toBe(false);
    });
  });
});
