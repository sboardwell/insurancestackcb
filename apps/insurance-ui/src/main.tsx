import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './styles/index.css';
import { initializeFeatureFlags } from './features/flags';

// Initialize CloudBees Feature Management
initializeFeatureFlags()
  .then(() => {
    console.log('[App] Feature flags initialized successfully');
  })
  .catch((error) => {
    console.error('[App] Failed to initialize feature flags:', error);
  });

// Render the app
ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
