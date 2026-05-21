import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { CreditCard, CheckCircle, Clock, XCircle, DollarSign, TrendingUp, Filter, X } from 'lucide-react';
import { api } from '../services/api';
import AlertBanner from '../components/AlertBanner';
import type { Payment, Policy } from '../types';
import { useState } from 'react';
import { format } from 'date-fns';
import useRoxFlag from '../hooks/useRoxFlag';

export default function Payments() {
  const queryClient = useQueryClient();
  const paymentsFilters = useRoxFlag('paymentsFilters');
  const [typeFilter, setTypeFilter] = useState<string>('');
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [showMakePaymentModal, setShowMakePaymentModal] = useState(false);

  // Form state for new payment
  const [paymentForm, setPaymentForm] = useState({
    policyId: '',
    type: 'premium',
    amount: '',
    paymentMethod: 'credit_card',
  });

  // Fetch active policies for payment selection
  const { data: policies } = useQuery<Policy[]>({
    queryKey: ['policies'],
    queryFn: () => api.getPolicies({ status: 'active' }),
  });

  // Fetch payments data
  const {
    data: payments,
    isLoading,
    isError,
    error,
  } = useQuery<Payment[]>({
    queryKey: ['payments'],
    queryFn: () => api.getPayments(),
    refetchInterval: 30000,
  });

  // Mutation for making new payment
  const makePaymentMutation = useMutation({
    mutationFn: (paymentData: Partial<Payment>) => api.createPayment(paymentData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payments'] });
      setShowMakePaymentModal(false);
      setPaymentForm({ policyId: '', type: 'premium', amount: '', paymentMethod: 'credit_card' });
    },
  });

  const handleMakePaymentSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const policy = policies?.find(p => p.id === paymentForm.policyId);
    if (!policy) return;

    makePaymentMutation.mutate({
      policyId: paymentForm.policyId,
      customerId: policy.customerId,
      type: paymentForm.type as 'premium' | 'payout',
      amount: parseFloat(paymentForm.amount),
      status: 'pending',
    });
  };

  // Filter payments
  const filteredPayments = payments?.filter((payment) => {
    const matchesType = !typeFilter || payment.type === typeFilter;
    const matchesStatus = !statusFilter || payment.status === statusFilter;
    return matchesType && matchesStatus;
  });

  // Calculate summary statistics
  const summary = payments?.reduce(
    (acc, payment) => {
      acc.totalPayments += 1;
      acc.totalAmount += payment.amount;

      if (payment.status === 'completed') {
        acc.completedAmount += payment.amount;
      }
      if (payment.status === 'pending') {
        acc.pendingAmount += payment.amount;
      }
      if (payment.type === 'premium') {
        acc.premiumPayments += payment.amount;
      }

      return acc;
    },
    {
      totalPayments: 0,
      totalAmount: 0,
      completedAmount: 0,
      pendingAmount: 0,
      premiumPayments: 0,
    }
  );

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount);
  };

  const formatDate = (date: string) => {
    return format(new Date(date), 'MMM dd, yyyy');
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <CheckCircle className="w-5 h-5 text-green-600" />;
      case 'pending':
        return <Clock className="w-5 h-5 text-yellow-600" />;
      case 'failed':
        return <XCircle className="w-5 h-5 text-red-600" />;
      case 'refunded':
        return <DollarSign className="w-5 h-5 text-blue-600" />;
      default:
        return <CreditCard className="w-5 h-5 text-gray-600" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'badge-success';
      case 'pending':
        return 'badge-warning';
      case 'failed':
        return 'badge-error';
      case 'refunded':
        return 'badge-info';
      default:
        return 'badge-default';
    }
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'premium':
        return 'bg-blue-100 text-blue-600';
      case 'claim':
        return 'bg-green-100 text-green-600';
      case 'refund':
        return 'bg-purple-100 text-purple-600';
      default:
        return 'bg-gray-100 text-gray-600';
    }
  };

  // Loading state
  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="spinner border-brand-500 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading payments...</p>
          </div>
        </div>
      </div>
    );
  }

  // Error state
  if (isError) {
    return (
      <div className="space-y-6">
        <AlertBanner
          type="critical"
          title="Error Loading Payments"
          message={error instanceof Error ? error.message : 'Failed to load payment data. Please try again.'}
          dismissible={false}
        />
        <div className="card p-8 text-center">
          <p className="text-gray-600 mb-4">Unable to load payments</p>
          <button
            onClick={() => window.location.reload()}
            className="btn-primary"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Payments</h1>
          <p className="text-gray-600 mt-1">Track and manage your insurance payments.</p>
        </div>
        <button onClick={() => setShowMakePaymentModal(true)} className="btn-primary">
          Make Payment
        </button>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Total Payments</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {summary?.totalPayments || 0}
              </p>
            </div>
            <div className="bg-brand-100 rounded-full p-3">
              <CreditCard className="w-6 h-6 text-brand-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Total Amount</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {formatCurrency(summary?.totalAmount || 0)}
              </p>
            </div>
            <div className="bg-blue-100 rounded-full p-3">
              <DollarSign className="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Completed</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {formatCurrency(summary?.completedAmount || 0)}
              </p>
            </div>
            <div className="bg-green-100 rounded-full p-3">
              <CheckCircle className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Premium Paid</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {formatCurrency(summary?.premiumPayments || 0)}
              </p>
            </div>
            <div className="bg-purple-100 rounded-full p-3">
              <TrendingUp className="w-6 h-6 text-purple-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Filters */}
      {paymentsFilters && (
        <div className="card p-6">
          <div className="flex items-center space-x-4">
            <Filter className="w-5 h-5 text-gray-400" />
            <select
              value={typeFilter}
              onChange={(e) => setTypeFilter(e.target.value)}
              className="input min-w-[150px]"
            >
              <option value="">All Types</option>
              <option value="premium">Premium</option>
              <option value="claim">Claim</option>
              <option value="refund">Refund</option>
            </select>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="input min-w-[150px]"
            >
              <option value="">All Status</option>
              <option value="pending">Pending</option>
              <option value="completed">Completed</option>
              <option value="failed">Failed</option>
              <option value="refunded">Refunded</option>
            </select>
          </div>
        </div>
      )}

      {/* Payments List */}
      <div className="card">
        {filteredPayments && filteredPayments.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Payment
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Type
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Amount
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Date
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Method
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {filteredPayments.map((payment) => (
                  <tr key={payment.id} className="hover:bg-gray-50 transition-colors">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <div className="mr-3">{getStatusIcon(payment.status)}</div>
                        <div>
                          <div className="text-sm font-medium text-gray-900">
                            Payment #{payment.id.slice(0, 8)}
                          </div>
                          <div className="text-sm text-gray-500">
                            Policy: {payment.policyId ? payment.policyId.slice(0, 8) : 'N/A'}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${getTypeColor(payment.type)}`}>
                        {payment.type.charAt(0).toUpperCase() + payment.type.slice(1)}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm font-semibold text-gray-900">
                        {formatCurrency(payment.amount)}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`${getStatusColor(payment.status)}`}>
                        {payment.status.charAt(0).toUpperCase() + payment.status.slice(1)}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">{formatDate(payment.processedDate || payment.createdAt)}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900 capitalize">
                        {payment.paymentMethod || 'N/A'}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <button className="text-brand-600 hover:text-brand-900 transition-colors">
                        View Receipt
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="p-12 text-center">
            <div className="bg-gray-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
              <CreditCard className="w-8 h-8 text-gray-400" />
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              {payments?.length === 0 ? 'No Payments Yet' : 'No Payments Found'}
            </h3>
            <p className="text-gray-600 mb-6">
              {payments?.length === 0
                ? 'You have not made any payments yet'
                : 'Try adjusting your filter criteria'}
            </p>
          </div>
        )}
      </div>

      {/* Make Payment Modal */}
      {showMakePaymentModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-2xl w-full p-6">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold text-gray-900">Make Payment</h2>
              <button
                onClick={() => setShowMakePaymentModal(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            <form onSubmit={handleMakePaymentSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Policy <span className="text-red-500">*</span>
                </label>
                <select
                  required
                  value={paymentForm.policyId}
                  onChange={(e) => setPaymentForm({ ...paymentForm, policyId: e.target.value })}
                  className="input w-full"
                >
                  <option value="">Select a policy</option>
                  {policies?.map((policy) => (
                    <option key={policy.id} value={policy.id}>
                      {policy.policyNumber} - {policy.type.toUpperCase()} Insurance (Premium: {formatCurrency(policy.premium)})
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Payment Type <span className="text-red-500">*</span>
                </label>
                <select
                  required
                  value={paymentForm.type}
                  onChange={(e) => setPaymentForm({ ...paymentForm, type: e.target.value })}
                  className="input w-full"
                >
                  <option value="premium">Premium Payment</option>
                  <option value="payout">Payout</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Amount ($) <span className="text-red-500">*</span>
                </label>
                <input
                  type="number"
                  required
                  min="0"
                  step="0.01"
                  placeholder="0.00"
                  value={paymentForm.amount}
                  onChange={(e) => setPaymentForm({ ...paymentForm, amount: e.target.value })}
                  className="input w-full"
                />
                {paymentForm.policyId && paymentForm.type === 'premium' && (() => {
                  const selectedPolicy = policies?.find(p => p.id === paymentForm.policyId);
                  return selectedPolicy ? (
                    <p className="text-sm text-gray-500 mt-1">
                      Monthly premium: {formatCurrency(selectedPolicy.premium)}
                    </p>
                  ) : null;
                })()}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Payment Method <span className="text-red-500">*</span>
                </label>
                <select
                  required
                  value={paymentForm.paymentMethod}
                  onChange={(e) => setPaymentForm({ ...paymentForm, paymentMethod: e.target.value })}
                  className="input w-full"
                >
                  <option value="credit_card">Credit Card</option>
                  <option value="debit_card">Debit Card</option>
                  <option value="bank_transfer">Bank Transfer</option>
                  <option value="check">Check</option>
                </select>
              </div>

              <div className="flex justify-end space-x-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowMakePaymentModal(false)}
                  className="btn-secondary"
                  disabled={makePaymentMutation.isPending}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="btn-primary"
                  disabled={makePaymentMutation.isPending}
                >
                  {makePaymentMutation.isPending ? 'Processing...' : 'Make Payment'}
                </button>
              </div>

              {makePaymentMutation.isError && (
                <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-md">
                  <p className="text-sm text-red-600">
                    Failed to process payment. Please try again.
                  </p>
                </div>
              )}
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
