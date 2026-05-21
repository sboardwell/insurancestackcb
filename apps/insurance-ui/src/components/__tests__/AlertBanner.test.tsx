import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import AlertBanner from '../AlertBanner';

// Mock the feature flags module for useRoxFlag
vi.mock('../../features/flags', () => ({
  getFlagsSnapshot: () => ({
    alertsBanner: true,
    claimsFilters: true,
    paymentsFilters: true,
    enhancedPolicyView: false,
    quickClaimFiling: true,
  }),
  subscribeFlags: () => () => {}, // Return unsubscribe function
}));

describe('AlertBanner', () => {
  it('renders info alert correctly', () => {
    render(
      <AlertBanner
        type="info"
        title="Test Title"
        message="Test message"
      />
    );

    expect(screen.getByText('Test Title')).toBeInTheDocument();
    expect(screen.getByText('Test message')).toBeInTheDocument();
  });

  it('renders warning alert correctly', () => {
    render(
      <AlertBanner
        type="warning"
        title="Warning Title"
        message="Warning message"
      />
    );

    expect(screen.getByText('Warning Title')).toBeInTheDocument();
    expect(screen.getByText('Warning message')).toBeInTheDocument();
  });

  it('renders critical alert correctly', () => {
    render(
      <AlertBanner
        type="critical"
        title="Critical Title"
        message="Critical message"
      />
    );

    expect(screen.getByText('Critical Title')).toBeInTheDocument();
    expect(screen.getByText('Critical message')).toBeInTheDocument();
  });

  it('can be dismissed when dismissible is true', () => {
    render(
      <AlertBanner
        type="info"
        title="Test Title"
        message="Test message"
        dismissible={true}
      />
    );

    const dismissButton = screen.getByLabelText('Dismiss alert');
    expect(dismissButton).toBeInTheDocument();

    fireEvent.click(dismissButton);

    // Alert should no longer be visible after dismissing
    expect(screen.queryByText('Test Title')).not.toBeInTheDocument();
  });

  it('does not show dismiss button when dismissible is false', () => {
    render(
      <AlertBanner
        type="info"
        title="Test Title"
        message="Test message"
        dismissible={false}
      />
    );

    expect(screen.queryByLabelText('Dismiss alert')).not.toBeInTheDocument();
  });
});
