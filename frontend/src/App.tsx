import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { UserProvider } from './context/UserContext';
import { TenantProvider } from './context/TenantContext';
import { ThemeProvider } from './context/ThemeContext';
import { Login, Register, Dashboard, ProtectedRoute } from './components';
import DataSourceManagement from './components/DataSourceManagement';
import DataSourceUpload from './components/DataSourceUpload';
import MainLayout from './components/MainLayout';
import UserProfile from './components/UserProfile';
import MatchsetConfiguration from './components/MatchsetConfiguration';
import MatchedTransactions from './components/MatchedTransactions';
import UnmatchedTransactions from './components/UnmatchedTransactions';
import AdminTroubleshooter from './components/AdminTroubleshooter';
import ImportRecords from './components/ImportRecords';
import RawTransactionsList from './components/RawTransactionsList';

function App() {
  return (
    <ThemeProvider>
      <UserProvider>
        <TenantProvider>
          <Router>
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              
              {/* Protected routes with layout */}
              <Route
                path="/"
                element={
                  <ProtectedRoute>
                    <MainLayout />
                  </ProtectedRoute>
                }
              >
                <Route path="dashboard" element={<Dashboard />} />
                <Route path="datasources" element={<DataSourceManagement />} />
                <Route path="datasources/upload" element={<DataSourceUpload />} />
                <Route path="datasources/:dataSourceId/imports" element={<ImportRecords />} />
                <Route path="imports/:importId/transactions" element={<RawTransactionsList />} />
                <Route path="transaction-match/matchset" element={<MatchsetConfiguration />} />
                <Route path="transaction-match/matched" element={<MatchedTransactions />} />
                <Route path="transaction-match/unmatched" element={<UnmatchedTransactions />} />
                <Route path="profile" element={<UserProfile />} />
                <Route path="admin-troubleshooter" element={<AdminTroubleshooter />} />
                {/* Redirect root to dashboard */}
                <Route index element={<Navigate to="/dashboard" replace />} />
              </Route>
              
              {/* Redirect any unknown routes to dashboard if authenticated */}
              <Route path="*" element={<Navigate to="/dashboard" replace />} />
            </Routes>
          </Router>
        </TenantProvider>
      </UserProvider>
    </ThemeProvider>
  );
}

export default App; 