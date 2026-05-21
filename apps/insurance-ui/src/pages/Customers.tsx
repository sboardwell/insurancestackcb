import { useQuery } from '@tanstack/react-query';
import { Users, UserCheck, UserX, Search, Mail, Phone } from 'lucide-react';
import { api } from '../services/api';
import AlertBanner from '../components/AlertBanner';
import type { Customer } from '../types';
import { useState } from 'react';
import { format } from 'date-fns';

export default function Customers() {
  const [searchQuery, setSearchQuery] = useState('');

  // Fetch customers data
  const {
    data: customers,
    isLoading,
    isError,
    error,
  } = useQuery<Customer[]>({
    queryKey: ['customers'],
    queryFn: () => api.getCustomers(),
    refetchInterval: 30000,
  });

  // Filter customers based on search
  const filteredCustomers = customers?.filter((customer) => {
    const fullName = `${customer.firstName} ${customer.lastName}`.toLowerCase();
    const query = searchQuery.toLowerCase();
    return (
      fullName.includes(query) ||
      customer.email.toLowerCase().includes(query) ||
      customer.phone?.toLowerCase().includes(query)
    );
  });

  // Calculate summary statistics
  const summary = customers?.reduce(
    (acc, customer) => {
      acc.totalCustomers += 1;
      // Count customers by risk score
      if (customer.riskScore && customer.riskScore <= 2) {
        acc.activeCustomers += 1; // Low risk
      } else if (customer.riskScore && customer.riskScore >= 4) {
        acc.inactiveCustomers += 1; // High risk
      }
      return acc;
    },
    { totalCustomers: 0, activeCustomers: 0, inactiveCustomers: 0 }
  );

  // Loading state
  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="spinner border-brand-500 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading customers...</p>
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
          title="Error Loading Customers"
          message={error instanceof Error ? error.message : 'Failed to load customer data. Please try again.'}
          dismissible={false}
        />
        <div className="card p-8 text-center">
          <p className="text-gray-600 mb-4">Unable to load customers</p>
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
          <h1 className="text-3xl font-bold text-gray-900">Customers</h1>
          <p className="text-gray-600 mt-1">Manage your customer database.</p>
        </div>
        <button className="btn-primary">
          Add Customer
        </button>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Total Customers</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {summary?.totalCustomers || 0}
              </p>
            </div>
            <div className="bg-brand-100 rounded-full p-3">
              <Users className="w-6 h-6 text-brand-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Active</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {summary?.activeCustomers || 0}
              </p>
            </div>
            <div className="bg-green-100 rounded-full p-3">
              <UserCheck className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Inactive</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {summary?.inactiveCustomers || 0}
              </p>
            </div>
            <div className="bg-gray-100 rounded-full p-3">
              <UserX className="w-6 h-6 text-gray-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Search */}
      <div className="card p-6">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
          <input
            type="text"
            placeholder="Search customers by name, email, or phone..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="input pl-10 w-full"
          />
        </div>
      </div>

      {/* Customers List */}
      <div className="card">
        {filteredCustomers && filteredCustomers.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Customer
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Contact
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Risk Score
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Joined
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {filteredCustomers.map((customer) => (
                  <tr key={customer.id} className="hover:bg-gray-50 transition-colors">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <div className="w-10 h-10 bg-brand-100 rounded-full flex items-center justify-center mr-3">
                          <span className="text-brand-600 font-medium text-sm">
                            {customer.firstName[0]}{customer.lastName[0]}
                          </span>
                        </div>
                        <div>
                          <div className="text-sm font-medium text-gray-900">
                            {customer.firstName} {customer.lastName}
                          </div>
                          <div className="text-sm text-gray-500">ID: {customer.id.slice(0, 8)}</div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="space-y-1">
                        <div className="flex items-center text-sm text-gray-900">
                          <Mail className="w-4 h-4 mr-2 text-gray-400" />
                          {customer.email}
                        </div>
                        {customer.phone && (
                          <div className="flex items-center text-sm text-gray-500">
                            <Phone className="w-4 h-4 mr-2 text-gray-400" />
                            {customer.phone}
                          </div>
                        )}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`badge ${customer.riskScore && customer.riskScore <= 2 ? 'badge-success' : customer.riskScore && customer.riskScore >= 4 ? 'badge-error' : 'badge-warning'}`}>
                        Risk: {customer.riskScore || 'N/A'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">
                        {format(new Date(customer.createdAt), 'MMM dd, yyyy')}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <button className="text-brand-600 hover:text-brand-900 transition-colors">
                        View Details
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
              <Users className="w-8 h-8 text-gray-400" />
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              {customers?.length === 0 ? 'No Customers Yet' : 'No Customers Found'}
            </h3>
            <p className="text-gray-600 mb-6">
              {customers?.length === 0
                ? 'Get started by adding your first customer'
                : 'Try adjusting your search criteria'}
            </p>
            {customers?.length === 0 && (
              <button className="btn-primary">
                Add Your First Customer
              </button>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
