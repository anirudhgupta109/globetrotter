import React from 'react';
import { useNavigate } from 'react-router-dom';

const GameMenu = () => {
  const navigate = useNavigate();
  const username = localStorage.getItem('username');

  const handleStartGame = () => {
    navigate('/play');
  };

  const handleStartChallenge = () => {
    navigate('/challenge');
  };

  const handleLogout = () => {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('username');
    navigate('/login');
  };

  return (
    <div className="menu-container">
      <h1 className="menu-title">Globetrotter</h1>
      
      <p style={{ marginBottom: '20px' }}>Welcome, {username}!</p>
      
      <button className="menu-button" onClick={handleStartGame}>
        Start Game
      </button>
      
      <button className="menu-button challenge" onClick={handleStartChallenge}>
        Challenge a Friend
      </button>
      
      <button className="menu-button" style={{ backgroundColor: '#f44336' }} onClick={handleLogout}>
        Logout
      </button>
    </div>
  );
};

export default GameMenu;