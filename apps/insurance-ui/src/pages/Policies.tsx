import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { FileText, Shield, TrendingUp, AlertCircle } from 'lucide-react';
import { api } from '../services/api';
import PolicyCard from '../components/PolicyCard';
import PolicyDetailModal from '../components/PolicyDetailModal';
import AlertBanner from '../components/AlertBanner';
import type { Policy } from '../types';

export default function Policies() {
  const navigate = useNavigate();
  const [selectedPolicy, setSelectedPolicy] = useState<Policy | null>(null);
  // Fetch policies data
  const {
    data: policies,
    isLoading,
    isError,
    error,
  } = useQuery<Policy[]>({
    queryKey: ['policies'],
    queryFn: () => api.getPolicies(),
    refetchInterval: 30000, // Refetch every 30 seconds
    refetchOnMount: 'always', // Always refetch when component mounts
    staleTime: 0, // Consider data stale immediately
  });

  // Calculate summary statistics
  const summary = policies?.reduce(
    (acc, policy) => {
      if (policy.status === 'active') {
        acc.totalPolicies += 1;
        acc.totalPremium += policy.premium;
        acc.totalCoverage += policy.coverage;
      }
      return acc;
    },
    { totalPolicies: 0, totalPremium: 0, totalCoverage: 0 }
  );

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount);
  };

  // Loading state
  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="spinner border-brand-500 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading your policies...</p>
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
          title="Error Loading Policies"
          message={error instanceof Error ? error.message : 'Failed to load policy data. Please try again.'}
          dismissible={false}
        />
        <div className="card p-8 text-center">
          <p className="text-gray-600 mb-4">Unable to load your policies</p>
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
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Policies</h1>
        <p className="text-gray-600 mt-1">Manage and view your insurance policies.</p>
      </div>

      {/* Alert Banner */}
      <AlertBanner
        type="info"
        title="Policy Renewal Reminder"
        message="Check your policies for upcoming renewals and ensure continuous coverage!"
      />

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Active Policies</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {summary?.totalPolicies || 0}
              </p>
              <p className="text-xs text-gray-500 mt-2">Currently active</p>
            </div>
            <div className="bg-brand-100 rounded-full p-3">
              <FileText className="w-6 h-6 text-brand-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Total Coverage</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {formatCurrency(summary?.totalCoverage || 0)}
              </p>
              <p className="text-xs text-gray-500 mt-2">Protection amount</p>
            </div>
            <div className="bg-green-100 rounded-full p-3">
              <Shield className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Annual Premium</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {formatCurrency((summary?.totalPremium || 0) * 12)}
              </p>
              <p className="text-xs text-gray-500 mt-2">Yearly payment</p>
            </div>
            <div className="bg-blue-100 rounded-full p-3">
              <TrendingUp className="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Policies Section */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-gray-900">Your Policies</h2>
          <button className="btn-primary">
            Add Policy
          </button>
        </div>

        {policies && policies.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {policies.map((policy) => (
              <PolicyCard
                key={policy.id}
                policy={policy}
                onViewDetails={() => setSelectedPolicy(policy)}
              />
            ))}
          </div>
        ) : (
          <div className="card p-12 text-center">
            <div className="bg-gray-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
              <FileText className="w-8 h-8 text-gray-400" />
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">No Policies Yet</h3>
            <p className="text-gray-600 mb-6">Get started by requesting a quote or adding your first policy</p>
            <button className="btn-primary">
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
            onClick={() => navigate('/claims')}
            className="flex items-center space-x-3 p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors"
          >
            <div className="bg-green-100 rounded-lg p-2">
              <FileText className="w-5 h-5 text-green-600" />
            </div>
            <div className="text-left">
              <p className="text-sm font-semibold text-gray-900">View Claims</p>
              <p className="text-xs text-gray-500">Check claim status</p>
            </div>
          </button>

          <button
            onClick={() => navigate('/quote')}
            className="flex items-center space-x-3 p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors"
          >
            <div className="bg-purple-100 rounded-lg p-2">
              <TrendingUp className="w-5 h-5 text-purple-600" />
            </div>
            <div className="text-left">
              <p className="text-sm font-semibold text-gray-900">Get Quote</p>
              <p className="text-xs text-gray-500">New insurance</p>
            </div>
          </button>
        </div>
      </div>

      {/* Policy Detail Modal */}
      {selectedPolicy && (
        <PolicyDetailModal
          policy={selectedPolicy}
          onClose={() => setSelectedPolicy(null)}
        />
      )}
    </div>
  );
}
