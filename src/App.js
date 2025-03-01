import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Register from './components/Register';
import Login from './components/Login';
import GameMenu from './components/GameMenu';
import QuestionPage from './components/QuestionPage';
import Challenge from './components/Challenge';
import ChallengeComplete from './components/ChallengeComplete';
import './styles.css';

function ProtectedRoute({ children }) {
  const token = localStorage.getItem('auth_token');
  if (!token) {
    return <Navigate to="/" replace />;
  }
  return children;
}

function App() {
  return (
    <Router>
      <div className="app-container">
        <Routes>
          <Route path="/" element={<Register />} />
          <Route path="/login" element={<Login />} />
          <Route 
            path="/menu" 
            element={
              <ProtectedRoute>
                <GameMenu />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/play" 
            element={
              <ProtectedRoute>
                <QuestionPage />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/challenge/:id?" 
            element={
              <ProtectedRoute>
                <Challenge />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/challenge-complete" 
            element={
              <ProtectedRoute>
                <ChallengeComplete />
              </ProtectedRoute>
            } 
          />
        </Routes>
      </div>
    </Router>
  );
}

export default App;