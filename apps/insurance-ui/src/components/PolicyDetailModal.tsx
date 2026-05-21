import { X, Shield, Car, Home, Heart, Activity, Calendar, DollarSign, FileText, CheckCircle } from 'lucide-react';
import type { Policy } from '../types';
import { format } from 'date-fns';

interface PolicyDetailModalProps {
  policy: Policy;
  onClose: () => void;
}

export default function PolicyDetailModal({ policy, onClose }: PolicyDetailModalProps) {
  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
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
        return 'bg-green-100 text-green-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'expired':
        return 'bg-red-100 text-red-800';
      case 'cancelled':
        return 'bg-gray-100 text-gray-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const Icon = getPolicyIcon(policy.type);
  const typeColorClass = getPolicyTypeColor(policy.type);

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:p-0">
        {/* Background overlay */}
        <div
          className="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75"
          onClick={onClose}
        ></div>

        {/* Modal panel */}
        <div className="relative inline-block w-full max-w-2xl my-8 overflow-hidden text-left align-middle transition-all transform bg-white rounded-lg shadow-xl">
          {/* Header */}
          <div className="flex items-start justify-between p-6 border-b border-gray-200">
            <div className="flex items-center space-x-4">
              <div className={`${typeColorClass} rounded-lg p-3`}>
                <Icon className="w-8 h-8" />
              </div>
              <div>
                <h2 className="text-2xl font-bold text-gray-900 capitalize">
                  {policy.type} Insurance Policy
                </h2>
                <p className="text-sm text-gray-500">Policy #{policy.policyNumber}</p>
              </div>
            </div>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-500 transition-colors"
            >
              <X className="w-6 h-6" />
            </button>
          </div>

          {/* Body */}
          <div className="p-6 space-y-6">
            {/* Status */}
            <div>
              <label className="text-sm font-medium text-gray-500">Status</label>
              <div className="mt-1">
                <span className={`px-3 py-1 text-sm font-medium rounded-full ${getStatusColor(policy.status)}`}>
                  {policy.status.charAt(0).toUpperCase() + policy.status.slice(1)}
                </span>
              </div>
            </div>

            {/* Coverage Details */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="card p-4 bg-gray-50">
                <div className="flex items-center space-x-2 mb-2">
                  <Shield className="w-4 h-4 text-brand-600" />
                  <p className="text-sm font-medium text-gray-600">Coverage Amount</p>
                </div>
                <p className="text-2xl font-bold text-gray-900">
                  {formatCurrency(policy.coverage)}
                </p>
              </div>

              <div className="card p-4 bg-gray-50">
                <div className="flex items-center space-x-2 mb-2">
                  <DollarSign className="w-4 h-4 text-brand-600" />
                  <p className="text-sm font-medium text-gray-600">Monthly Premium</p>
                </div>
                <p className="text-2xl font-bold text-gray-900">
                  {formatCurrency(policy.premium)}
                </p>
                <p className="text-xs text-gray-500 mt-1">per month</p>
              </div>

              <div className="card p-4 bg-gray-50">
                <div className="flex items-center space-x-2 mb-2">
                  <FileText className="w-4 h-4 text-brand-600" />
                  <p className="text-sm font-medium text-gray-600">Deductible</p>
                </div>
                <p className="text-2xl font-bold text-gray-900">
                  {policy.deductible ? formatCurrency(policy.deductible) : 'N/A'}
                </p>
                <p className="text-xs text-gray-500 mt-1">per claim</p>
              </div>
            </div>

            {/* Policy Dates */}
            <div className="space-y-3">
              <h3 className="text-lg font-semibold text-gray-900">Policy Period</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="flex items-start space-x-3">
                  <Calendar className="w-5 h-5 text-gray-400 mt-0.5" />
                  <div>
                    <p className="text-sm font-medium text-gray-600">Start Date</p>
                    <p className="text-base text-gray-900">
                      {format(new Date(policy.startDate), 'MMMM dd, yyyy')}
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-3">
                  <Calendar className="w-5 h-5 text-gray-400 mt-0.5" />
                  <div>
                    <p className="text-sm font-medium text-gray-600">End Date</p>
                    <p className="text-base text-gray-900">
                      {format(new Date(policy.endDate), 'MMMM dd, yyyy')}
                    </p>
                  </div>
                </div>
              </div>
              {policy.renewalDate && new Date(policy.renewalDate).getFullYear() > 2000 && (
                <div className="flex items-start space-x-3">
                  <Calendar className="w-5 h-5 text-gray-400 mt-0.5" />
                  <div>
                    <p className="text-sm font-medium text-gray-600">Renewal Date</p>
                    <p className="text-base text-gray-900">
                      {format(new Date(policy.renewalDate), 'MMMM dd, yyyy')}
                    </p>
                  </div>
                </div>
              )}
            </div>

            {/* Policy Information */}
            <div className="space-y-3">
              <h3 className="text-lg font-semibold text-gray-900">Policy Information</h3>
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Policy ID</span>
                  <span className="font-medium text-gray-900">{policy.id}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Customer ID</span>
                  <span className="font-medium text-gray-900">{policy.customerId}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Currency</span>
                  <span className="font-medium text-gray-900">{policy.currency}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Created</span>
                  <span className="font-medium text-gray-900">
                    {format(new Date(policy.createdAt), 'MMM dd, yyyy')}
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Last Updated</span>
                  <span className="font-medium text-gray-900">
                    {format(new Date(policy.updatedAt), 'MMM dd, yyyy')}
                  </span>
                </div>
              </div>
            </div>

            {/* Coverage Includes */}
            <div className="space-y-3">
              <h3 className="text-lg font-semibold text-gray-900">Coverage Includes</h3>
              <ul className="space-y-2">
                <li className="flex items-start">
                  <CheckCircle className="w-5 h-5 text-green-600 mr-3 mt-0.5 flex-shrink-0" />
                  <span className="text-gray-700">Comprehensive coverage up to policy limits</span>
                </li>
                <li className="flex items-start">
                  <CheckCircle className="w-5 h-5 text-green-600 mr-3 mt-0.5 flex-shrink-0" />
                  <span className="text-gray-700">24/7 customer support and claims assistance</span>
                </li>
                <li className="flex items-start">
                  <CheckCircle className="w-5 h-5 text-green-600 mr-3 mt-0.5 flex-shrink-0" />
                  <span className="text-gray-700">Fast claim processing and payment</span>
                </li>
                <li className="flex items-start">
                  <CheckCircle className="w-5 h-5 text-green-600 mr-3 mt-0.5 flex-shrink-0" />
                  <span className="text-gray-700">No hidden fees or surprise charges</span>
                </li>
              </ul>
            </div>
          </div>

          {/* Footer */}
          <div className="flex justify-end px-6 py-4 bg-gray-50 border-t border-gray-200">
            <button
              onClick={onClose}
              className="btn-secondary mr-3"
            >
              Close
            </button>
            <button className="btn-primary">
              File a Claim
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
