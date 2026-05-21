import { useQuery } from '@tanstack/react-query';
import { Shield, FileText, DollarSign, AlertCircle, TrendingUp, Activity } from 'lucide-react';
import { api } from '../services/api';
import PolicyCard from '../components/PolicyCard';
import AlertBanner from '../components/AlertBanner';
import type { Policy, Claim, Payment } from '../types';
import { useNavigate } from 'react-router-dom';

export default function Dashboard() {
  const navigate = useNavigate();

  // Fetch policies data
  const {
    data: policies,
    isLoading: policiesLoading,
  } = useQuery<Policy[]>({
    queryKey: ['policies'],
    queryFn: () => api.getPolicies(),
    refetchInterval: 30000,
  });

  // Fetch claims data
  const {
    data: claims,
    isLoading: claimsLoading,
  } = useQuery<Claim[]>({
    queryKey: ['claims'],
    queryFn: () => api.getClaims(),
    refetchInterval: 30000,
  });

  // Fetch payments data
  const {
    isLoading: paymentsLoading,
  } = useQuery<Payment[]>({
    queryKey: ['payments'],
    queryFn: () => api.getPayments(),
    refetchInterval: 30000,
  });

  // Calculate summary statistics
  const summary = {
    activePolicies: policies?.filter(p => p.status === 'active').length || 0,
    totalCoverage: policies?.filter(p => p.status === 'active').reduce((sum, p) => sum + p.coverage, 0) || 0,
    monthlyPremium: policies?.filter(p => p.status === 'active').reduce((sum, p) => sum + p.premium, 0) || 0,
    activeClaims: claims?.filter(c => c.status === 'submitted' || c.status === 'under_review').length || 0,
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount);
  };

  const isLoading = policiesLoading || claimsLoading || paymentsLoading;

  // Loading state
  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="spinner border-brand-500 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading your dashboard...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-600 mt-1">Welcome back! Here's your insurance overview.</p>
      </div>

      {/* Alert Banner */}
      <AlertBanner
        type="info"
        title="Policy Renewal Reminder"
        message="Check your policies for upcoming renewals to ensure continuous coverage!"
      />

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Active Policies</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {summary.activePolicies}
              </p>
              <p className="text-xs text-gray-500 mt-2">Currently insured</p>
            </div>
            <div className="bg-brand-100 rounded-full p-3">
              <Shield className="w-6 h-6 text-brand-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Total Coverage</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {formatCurrency(summary.totalCoverage)}
              </p>
              <p className="text-xs text-gray-500 mt-2">Protection amount</p>
            </div>
            <div className="bg-green-100 rounded-full p-3">
              <TrendingUp className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Monthly Premium</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {formatCurrency(summary.monthlyPremium)}
              </p>
              <p className="text-xs text-gray-500 mt-2">Per month</p>
            </div>
            <div className="bg-blue-100 rounded-full p-3">
              <DollarSign className="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Active Claims</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {summary.activeClaims}
              </p>
              <p className="text-xs text-gray-500 mt-2">In process</p>
            </div>
            <div className="bg-purple-100 rounded-full p-3">
              <Activity className="w-6 h-6 text-purple-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Recent Policies Section */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-gray-900">Recent Policies</h2>
          <button onClick={() => navigate('/policies')} className="text-brand-600 hover:text-brand-700 text-sm font-medium">
            View All
          </button>
        </div>

        {policies && policies.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {policies.slice(0, 3).map((policy) => (
              <PolicyCard key={policy.id} policy={policy} onViewDetails={() => {}} />
            ))}
          </div>
        ) : (
          <div className="card p-12 text-center">
            <div className="bg-gray-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
              <Shield className="w-8 h-8 text-gray-400" />
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">No Policies Yet</h3>
            <p className="text-gray-600 mb-6">Get started by requesting a quote for your first policy</p>
            <button onClick={() => navigate('/quote')} className="btn-primary">
              Get a Quote
            </button>
          </div>
        )}
      </div>

      {/* Quick Actions */}
      <div className="card p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <button
            onClick={() => navigate('/claims')}
            className="flex items-center space-x-3 p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors"
          >
            <div className="bg-brand-100 rounded-lg p-2">
              <AlertCircle className="w-5 h-5 text-brand-600" />
            </div>
            <div className="text-left">
              <p className="text-sm font-semibold text-gray-900">File a Claim</p>
              <p className="text-xs text-gray-500">Report an incident</p>
            </div>
          </button>

          <button
            onClick={() => navigate('/quote')}
            className="flex items-center space-x-3 p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors"
          >
            <div className="bg-green-100 rounded-lg p-2">
              <FileText className="w-5 h-5 text-green-600" />
            </div>
            <div className="text-left">
              <p className="text-sm font-semibold text-gray-900">Get a Quote</p>
              <p className="text-xs text-gray-500">New policy</p>
            </div>
          </button>

          <button
            onClick={() => navigate('/payments')}
            className="flex items-center space-x-3 p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors"
          >
            <div className="bg-purple-100 rounded-lg p-2">
              <DollarSign className="w-5 h-5 text-purple-600" />
            </div>
            <div className="text-left">
              <p className="text-sm font-semibold text-gray-900">Make Payment</p>
              <p className="text-xs text-gray-500">Pay premium</p>
            </div>
          </button>
        </div>
      </div>
    </div>
  );
}
