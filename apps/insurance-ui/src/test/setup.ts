// Vitest setup file
import { expect, afterEach } from 'vitest';
import { cleanup } from '@testing-library/react';
import '@testing-library/jest-dom';

// Cleanup after each test
afterEach(() => {
  cleanup();
});

// Mock environment variables
process.env.VITE_API_BASE_URL = 'http://localhost:3000/api';
process.env.VITE_ROX_API_KEY = 'test-api-key';

// Mock window.matchMedia for responsive design tests
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: (query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {}, // deprecated
    removeListener: () => {}, // deprecated
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => {},
  }),
});

// Mock IntersectionObserver
global.IntersectionObserver = class IntersectionObserver {
  constructor() {}
  disconnect() {}
  observe() {}
  takeRecords() {
    return [];
  }
  unobserve() {}
} as any;

// Suppress console warnings in tests
const originalWarn = console.warn;
const originalError = console.error;

beforeEach(() => {
  console.warn = (...args: any[]) => {
    // Suppress specific warnings
    if (
      typeof args[0] === 'string' &&
      (args[0].includes('ReactDOM.render') ||
        args[0].includes('Not implemented: HTMLFormElement.prototype.submit'))
    ) {
      return;
    }
    originalWarn.call(console, ...args);
  };

  console.error = (...args: any[]) => {
    // Suppress specific errors
    if (
      typeof args[0] === 'string' &&
      (args[0].includes('Not implemented: HTMLFormElement.prototype.submit') ||
        args[0].includes('Could not parse CSS stylesheet'))
    ) {
      return;
    }
    originalError.call(console, ...args);
  };
});

afterEach(() => {
  console.warn = originalWarn;
  console.error = originalError;
});
