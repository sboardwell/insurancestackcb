import { Shield, Car, Home, Heart, Activity, Calendar } from 'lucide-react';
import type { Policy } from '../types';
import { format } from 'date-fns';
import useRoxFlag from '../hooks/useRoxFlag';

interface PolicyCardProps {
  policy: Policy;
  onViewDetails?: () => void;
}

export default function PolicyCard({ policy, onViewDetails }: PolicyCardProps) {
  const enhancedPolicyView = useRoxFlag('enhancedPolicyView');
  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const getPolicyIcon = (type: string) => {
    switch (type) {
      case 'auto':
        return Car;
      case 'home':
        return Home;
      case 'life':
        return Heart;
      case 'health':
        return Activity;
      default:
        return Shield;
    }
  };

  const getPolicyTypeColor = (type: string) => {
    switch (type) {
      case 'auto':
        return 'bg-blue-100 text-blue-600';
      case 'home':
        return 'bg-green-100 text-green-600';
      case 'life':
        return 'bg-purple-100 text-purple-600';
      case 'health':
        return 'bg-red-100 text-red-600';
      default:
        return 'bg-gray-100 text-gray-600';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'badge-success';
      case 'pending':
        return 'badge-warning';
      case 'expired':
        return 'badge-error';
      case 'cancelled':
        return 'badge-default';
      default:
        return 'badge-default';
    }
  };

  const Icon = getPolicyIcon(policy.type);
  const typeColorClass = getPolicyTypeColor(policy.type);

  return (
    <div className="card p-6 hover:shadow-lg transition-shadow cursor-pointer">
      <div className="flex items-start justify-between mb-4">
        <div className={`${typeColorClass} rounded-lg p-3`}>
          <Icon className="w-6 h-6" />
        </div>
        <span className={`${getStatusColor(policy.status)}`}>
          {policy.status.charAt(0).toUpperCase() + policy.status.slice(1)}
        </span>
      </div>

      <div className="space-y-3">
        <div>
          <h3 className="text-lg font-semibold text-gray-900 capitalize">
            {policy.type} Insurance
          </h3>
          <p className="text-sm text-gray-500">Policy #{policy.policyNumber}</p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-600">Coverage</span>
            <span className="font-semibold text-gray-900">
              {formatCurrency(policy.coverage)}
            </span>
          </div>

          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-600">Premium</span>
            <span className="font-semibold text-gray-900">
              {formatCurrency(policy.premium)}/mo
            </span>
          </div>

          {policy.deductible && (
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-600">Deductible</span>
              <span className="font-semibold text-gray-900">
                {formatCurrency(policy.deductible)}
              </span>
            </div>
          )}

          {enhancedPolicyView && (
            <>
              {policy.renewalDate && (
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600">Renewal Date</span>
                  <span className="font-semibold text-gray-900">
                    {format(new Date(policy.renewalDate), 'MMM dd, yyyy')}
                  </span>
                </div>
              )}
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-600">Customer ID</span>
                <span className="font-mono text-xs text-gray-900">
                  {policy.customerId.slice(0, 8)}...
                </span>
              </div>
              {policy.currency && (
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600">Currency</span>
                  <span className="font-semibold text-gray-900">
                    {policy.currency}
                  </span>
                </div>
              )}
            </>
          )}
        </div>

        <div className="pt-3 border-t border-gray-200">
          <div className="flex items-center text-xs text-gray-500">
            <Calendar className="w-3 h-3 mr-1" />
            <span>
              {format(new Date(policy.startDate), 'MMM dd, yyyy')} - {format(new Date(policy.endDate), 'MMM dd, yyyy')}
            </span>
          </div>
        </div>
      </div>

      <div className="mt-4 pt-4 border-t border-gray-200">
        <button
          onClick={(e) => {
            e.stopPropagation();
            onViewDetails?.();
          }}
          className="w-full btn-secondary text-sm"
        >
          View Details
        </button>
      </div>
    </div>
  );
}
