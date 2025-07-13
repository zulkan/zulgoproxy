import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext.jsx'
import ProtectedRoute from './components/ProtectedRoute.jsx'
import Layout from './components/Layout.jsx'
import Login from './components/Login.jsx'
import Dashboard from './components/Dashboard.jsx'
import Users from './components/Users.jsx'
import Logs from './components/Logs.jsx'
import Health from './components/Health.jsx'
import Settings from './components/Settings.jsx'

function App() {
  return (
    <AuthProvider>
      <Router>
        <div className="App">
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/" element={
              <ProtectedRoute>
                <Layout />
              </ProtectedRoute>
            }>
              <Route index element={<Navigate to="/dashboard" replace />} />
              <Route path="dashboard" element={<Dashboard />} />
              <Route path="users" element={
                <ProtectedRoute adminOnly>
                  <Users />
                </ProtectedRoute>
              } />
              <Route path="logs" element={
                <ProtectedRoute adminOnly>
                  <Logs />
                </ProtectedRoute>
              } />
              <Route path="health" element={
                <ProtectedRoute adminOnly>
                  <Health />
                </ProtectedRoute>
              } />
              <Route path="settings" element={<Settings />} />
            </Route>
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App