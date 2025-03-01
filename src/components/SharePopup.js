import React, { useRef, useEffect } from 'react';

const SharePopup = ({ challengeId, show, onClose, inviter, score }) => {
  const linkRef = useRef(null);
  
  // Ensure challengeId is valid before constructing link
  const challengeLink = challengeId 
    ? `${window.location.origin}/challenge/${challengeId}`
    : '';
  
  const handleCopyLink = () => {
    if (linkRef.current) {
      linkRef.current.select();
      document.execCommand('copy');
      alert('Link copied to clipboard!');
    }
  };
  
  const handleWhatsAppShare = () => {
    const message = `Join my Globetrotter challenge! Can you beat my score? ${challengeLink}`;
    const whatsappUrl = `https://wa.me/?text=${encodeURIComponent(message)}`;
    window.open(whatsappUrl, '_blank');
  };
  
  useEffect(() => {
    if (show) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'auto';
    }
    
    return () => {
      document.body.style.overflow = 'auto';
    };
  }, [show]);
  
  if (!show) return null;
  
  return (
    <div className="share-overlay">
      <div className="share-container">
        <h2 className="share-title">Challenge Your Friends!</h2>
        
        <div className="share-image" style={{ 
          backgroundColor: '#f5f5f5', 
          height: '200px', 
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          border: '1px solid #ddd'
        }}>
          <h3 style={{ marginBottom: '10px' }}>Globetrotter Challenge</h3>
          <p>From: {inviter}</p>
          {score !== undefined && <p>Score: {score}</p>}
        </div>
        
        <p style={{ margin: '15px 0' }}>Share this link with your friends:</p>
        
        <input
          ref={linkRef}
          type="text"
          className="share-link"
          value={challengeLink}
          readOnly
          onClick={(e) => e.target.select()}
        />
        
        <div className="share-buttons">
          <button className="share-button" onClick={handleWhatsAppShare}>
            Share via WhatsApp
          </button>
          
          <button className="share-button" style={{ backgroundColor: '#2196f3' }} onClick={handleCopyLink}>
            Copy Link
          </button>
        </div>
        
        <button className="close-button" style={{ marginTop: '15px' }} onClick={onClose}>
          Close
        </button>
      </div>
    </div>
  );
};

export default SharePopup;