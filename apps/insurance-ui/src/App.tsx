import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AuthProvider } from './contexts/AuthContext';
import Layout from './components/Layout';
import ProtectedRoute from './components/ProtectedRoute';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Policies from './pages/Policies';
import Claims from './pages/Claims';
import Customers from './pages/Customers';
import Payments from './pages/Payments';
import GetQuote from './pages/GetQuote';

// Create a client for React Query
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 2,
      refetchOnWindowFocus: false,
      staleTime: 5 * 60 * 1000, // 5 minutes
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            {/* Public Route */}
            <Route path="/login" element={<Login />} />

            {/* Protected Routes */}
            <Route
              path="/*"
              element={
                <ProtectedRoute>
                  <Layout>
                    <Routes>
                      <Route path="/" element={<Dashboard />} />
                      <Route path="dashboard" element={<Dashboard />} />
                      <Route path="policies" element={<Policies />} />
                      <Route path="claims" element={<Claims />} />
                      <Route path="customers" element={<Customers />} />
                      <Route path="payments" element={<Payments />} />
                      <Route path="quote" element={<GetQuote />} />
                      {/* 404 Route */}
                      <Route
                        path="*"
                        element={
                          <div className="flex items-center justify-center min-h-[400px]">
                            <div className="text-center">
                              <h1 className="text-4xl font-bold text-gray-900 mb-2">404</h1>
                              <p className="text-gray-600 mb-6">Page not found</p>
                              <a href="/" className="btn-primary">
                                Go to Dashboard
                              </a>
                            </div>
                          </div>
                        }
                      />
                    </Routes>
                  </Layout>
                </ProtectedRoute>
              }
            />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
