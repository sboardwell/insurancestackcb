import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Car, Home, Heart, Shield, CheckCircle } from 'lucide-react';
import AlertBanner from '../components/AlertBanner';
import { api } from '../services/api';
import useRoxFlag from '../hooks/useRoxFlag';

export default function GetQuote() {
  const navigate = useNavigate();
  const killGetQuote = useRoxFlag('killGetQuote');
  const [step, setStep] = useState(1);
  const [policyType, setPolicyType] = useState<string>('');
  const [formData, setFormData] = useState({
    coverage: '',
    deductible: '',
    firstName: '',
    lastName: '',
    email: '',
    phone: '',
  });
  const [quoteResult, setQuoteResult] = useState<{
    premium: number;
    coverage: number;
    deductible: number;
  } | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitSuccess, setSubmitSuccess] = useState(false);

  const policyTypes = [
    {
      id: 'auto',
      name: 'Auto Insurance',
      icon: Car,
      description: 'Protect your vehicle and drive with confidence',
      color: 'bg-blue-100 text-blue-600',
    },
    {
      id: 'home',
      name: 'Home Insurance',
      icon: Home,
      description: 'Safeguard your home and belongings',
      color: 'bg-green-100 text-green-600',
    },
    {
      id: 'life',
      name: 'Life Insurance',
      icon: Heart,
      description: 'Secure your family\'s financial future',
      color: 'bg-purple-100 text-purple-600',
    },
  ];

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleGetQuote = () => {
    // Simulate quote calculation
    const coverage = parseFloat(formData.coverage) || 100000;
    const deductible = parseFloat(formData.deductible) || 1000;
    const premium = (coverage * 0.001) + (deductible * 0.01);

    setQuoteResult({
      premium: Math.round(premium * 100) / 100,
      coverage,
      deductible,
    });
    setStep(3);
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const handleAcceptQuote = async () => {
    if (!quoteResult) return;

    setIsSubmitting(true);
    setSubmitError(null);

    try {
      // Generate policy number
      const timestamp = Date.now();
      const typePrefix = policyType.toUpperCase();
      const policyNumber = `${typePrefix}-${new Date().getFullYear()}-${timestamp.toString().slice(-6)}`;

      // Calculate dates (policy starts today, ends in 1 year)
      const startDate = new Date().toISOString();
      const endDate = new Date();
      endDate.setFullYear(endDate.getFullYear() + 1);

      // Create the policy with the quote data
      const policyData = {
        policyNumber,
        type: policyType as 'auto' | 'home' | 'life' | 'health',
        premium: quoteResult.premium,
        coverage: quoteResult.coverage,
        deductible: quoteResult.deductible,
        startDate,
        endDate: endDate.toISOString(),
      };

      await api.createPolicy(policyData);

      // Show success message
      setSubmitSuccess(true);

      // Small delay to show success message, then navigate to policies page
      setTimeout(() => {
        navigate('/policies');
      }, 800);
    } catch (error) {
      console.error('Failed to create policy:', error);
      setSubmitError(error instanceof Error ? error.message : 'Failed to create policy. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // If killGetQuote flag is enabled, show maintenance message
  if (killGetQuote) {
    return (
      <div className="space-y-6">
        {/* Page Header */}
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Get a Quote</h1>
          <p className="text-gray-600 mt-1">Find the perfect insurance coverage for your needs.</p>
        </div>

        {/* Maintenance Message */}
        <AlertBanner
          type="critical"
          title="Temporarily Unavailable"
          message="The Get a Quote feature is currently down for maintenance. Please check back later."
          dismissible={false}
        />

        <div className="card p-12 text-center">
          <div className="bg-gray-100 rounded-full w-20 h-20 flex items-center justify-center mx-auto mb-6">
            <Shield className="w-10 h-10 text-gray-400" />
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-3">
            Down for Maintenance
          </h2>
          <p className="text-gray-600 mb-8 max-w-md mx-auto">
            We're working hard to improve our quote system. In the meantime, you can still view your existing policies or contact our support team for assistance.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <button
              onClick={() => navigate('/policies')}
              className="btn-primary"
            >
              View My Policies
            </button>
            <button
              onClick={() => navigate('/')}
              className="btn-secondary"
            >
              Return to Dashboard
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Get a Quote</h1>
        <p className="text-gray-600 mt-1">Find the perfect insurance coverage for your needs.</p>
      </div>

      {/* Progress Steps */}
      <div className="card p-6">
        <div className="flex items-center justify-between">
          {[
            { num: 1, label: 'Select Type' },
            { num: 2, label: 'Enter Details' },
            { num: 3, label: 'Get Quote' },
          ].map((s, idx) => (
            <div key={s.num} className="flex items-center">
              <div
                className={`flex items-center justify-center w-10 h-10 rounded-full font-semibold ${
                  step >= s.num
                    ? 'bg-brand-500 text-white'
                    : 'bg-gray-200 text-gray-600'
                }`}
              >
                {step > s.num ? <CheckCircle className="w-6 h-6" /> : s.num}
              </div>
              <span
                className={`ml-2 font-medium ${
                  step >= s.num ? 'text-gray-900' : 'text-gray-500'
                }`}
              >
                {s.label}
              </span>
              {idx < 2 && (
                <div
                  className={`w-16 h-1 mx-4 ${
                    step > s.num ? 'bg-brand-500' : 'bg-gray-200'
                  }`}
                />
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Step 1: Select Policy Type */}
      {step === 1 && (
        <div className="space-y-6">
          <AlertBanner
            type="info"
            title="Choose Your Coverage"
            message="Click one of the cards below to select your insurance type."
          />
          <div>
            <h2 className="text-lg font-semibold text-gray-900 mb-4">
              Select Insurance Type <span className="text-sm font-normal text-gray-500">(Click a card)</span>
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {policyTypes.map((type) => {
              const Icon = type.icon;
              return (
                <div
                  key={type.id}
                  onClick={() => setPolicyType(type.id)}
                  className={`card p-8 cursor-pointer transition-all hover:shadow-lg hover:border-brand-300 hover:scale-105 ${
                    policyType === type.id ? 'ring-4 ring-brand-500 shadow-xl' : ''
                  }`}
                >
                  <div className={`${type.color} rounded-lg p-4 w-fit mb-4`}>
                    <Icon className="w-8 h-8" />
                  </div>
                  <h3 className="text-xl font-semibold text-gray-900 mb-2">{type.name}</h3>
                  <p className="text-gray-600">{type.description}</p>
                </div>
              );
            })}
            </div>
          </div>
          <div className="flex justify-end">
            <button
              onClick={() => setStep(2)}
              disabled={!policyType}
              className="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Continue
            </button>
          </div>
        </div>
      )}

      {/* Step 2: Enter Details */}
      {step === 2 && (
        <div className="space-y-6">
          <div className="card p-8">
            <h2 className="text-xl font-semibold text-gray-900 mb-6">Coverage Details</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Coverage Amount
                </label>
                <select
                  name="coverage"
                  value={formData.coverage}
                  onChange={handleInputChange}
                  className="input w-full"
                  required
                >
                  <option value="">Select coverage amount</option>
                  <option value="50000">$50,000</option>
                  <option value="100000">$100,000</option>
                  <option value="250000">$250,000</option>
                  <option value="500000">$500,000</option>
                  <option value="1000000">$1,000,000</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Deductible
                </label>
                <select
                  name="deductible"
                  value={formData.deductible}
                  onChange={handleInputChange}
                  className="input w-full"
                  required
                >
                  <option value="">Select deductible</option>
                  <option value="500">$500</option>
                  <option value="1000">$1,000</option>
                  <option value="2500">$2,500</option>
                  <option value="5000">$5,000</option>
                </select>
              </div>
            </div>

            <h2 className="text-xl font-semibold text-gray-900 mb-6 mt-8">Personal Information</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  First Name
                </label>
                <input
                  type="text"
                  name="firstName"
                  value={formData.firstName}
                  onChange={handleInputChange}
                  className="input w-full"
                  placeholder="John"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Last Name
                </label>
                <input
                  type="text"
                  name="lastName"
                  value={formData.lastName}
                  onChange={handleInputChange}
                  className="input w-full"
                  placeholder="Doe"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Email
                </label>
                <input
                  type="email"
                  name="email"
                  value={formData.email}
                  onChange={handleInputChange}
                  className="input w-full"
                  placeholder="john.doe@example.com"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Phone
                </label>
                <input
                  type="tel"
                  name="phone"
                  value={formData.phone}
                  onChange={handleInputChange}
                  className="input w-full"
                  placeholder="(555) 123-4567"
                  required
                />
              </div>
            </div>
          </div>

          <div className="flex justify-between">
            <button onClick={() => setStep(1)} className="btn-secondary">
              Back
            </button>
            <button
              onClick={handleGetQuote}
              disabled={
                !formData.coverage ||
                !formData.deductible ||
                !formData.firstName ||
                !formData.lastName ||
                !formData.email
              }
              className="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Get Quote
            </button>
          </div>
        </div>
      )}

      {/* Step 3: Quote Result */}
      {step === 3 && quoteResult && (
        <div className="space-y-6">
          {submitSuccess ? (
            <AlertBanner
              type="success"
              title="Policy Created Successfully!"
              message="Your new policy has been created. Redirecting to your policies..."
            />
          ) : submitError ? (
            <AlertBanner
              type="critical"
              title="Error Creating Policy"
              message={submitError}
            />
          ) : (
            <AlertBanner
              type="info"
              title="Your Quote is Ready!"
              message="Review your quote below and accept it to create a new policy."
            />
          )}

          <div className="card p-8">
            <div className="text-center mb-8">
              <div className="bg-green-100 rounded-full w-20 h-20 flex items-center justify-center mx-auto mb-4">
                <Shield className="w-10 h-10 text-green-600" />
              </div>
              <h2 className="text-2xl font-bold text-gray-900 mb-2">
                Your {policyTypes.find((t) => t.id === policyType)?.name} Quote
              </h2>
              <p className="text-gray-600">Quote valid for 30 days</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
              <div className="card p-6 bg-gray-50">
                <p className="text-sm text-gray-600 font-medium mb-2">Monthly Premium</p>
                <p className="text-3xl font-bold text-brand-600">
                  {formatCurrency(quoteResult.premium)}
                </p>
                <p className="text-xs text-gray-500 mt-2">per month</p>
              </div>

              <div className="card p-6 bg-gray-50">
                <p className="text-sm text-gray-600 font-medium mb-2">Coverage Amount</p>
                <p className="text-3xl font-bold text-gray-900">
                  {formatCurrency(quoteResult.coverage)}
                </p>
                <p className="text-xs text-gray-500 mt-2">protection</p>
              </div>

              <div className="card p-6 bg-gray-50">
                <p className="text-sm text-gray-600 font-medium mb-2">Deductible</p>
                <p className="text-3xl font-bold text-gray-900">
                  {formatCurrency(quoteResult.deductible)}
                </p>
                <p className="text-xs text-gray-500 mt-2">per claim</p>
              </div>
            </div>

            <div className="space-y-4 mb-8">
              <h3 className="text-lg font-semibold text-gray-900">Coverage Includes:</h3>
              <ul className="space-y-3">
                <li className="flex items-start">
                  <CheckCircle className="w-5 h-5 text-green-600 mr-3 mt-0.5" />
                  <span className="text-gray-700">Comprehensive coverage up to policy limits</span>
                </li>
                <li className="flex items-start">
                  <CheckCircle className="w-5 h-5 text-green-600 mr-3 mt-0.5" />
                  <span className="text-gray-700">24/7 customer support and claims assistance</span>
                </li>
                <li className="flex items-start">
                  <CheckCircle className="w-5 h-5 text-green-600 mr-3 mt-0.5" />
                  <span className="text-gray-700">Fast claim processing and payment</span>
                </li>
                <li className="flex items-start">
                  <CheckCircle className="w-5 h-5 text-green-600 mr-3 mt-0.5" />
                  <span className="text-gray-700">No hidden fees or surprise charges</span>
                </li>
              </ul>
            </div>

            <div className="flex flex-col sm:flex-row gap-4">
              <button
                onClick={() => {
                  setStep(1);
                  setPolicyType('');
                  setFormData({
                    coverage: '',
                    deductible: '',
                    firstName: '',
                    lastName: '',
                    email: '',
                    phone: '',
                  });
                  setQuoteResult(null);
                  setSubmitError(null);
                  setSubmitSuccess(false);
                }}
                className="btn-secondary flex-1"
                disabled={isSubmitting}
              >
                Get Another Quote
              </button>
              <button
                onClick={handleAcceptQuote}
                disabled={isSubmitting || submitSuccess}
                className="btn-primary flex-1 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSubmitting ? 'Creating Policy...' : submitSuccess ? 'Policy Created!' : 'Accept Quote & Create Policy'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
