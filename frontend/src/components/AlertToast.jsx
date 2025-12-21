import React, { useEffect } from 'react';
import '../styles/AlertToast.css';

function AlertToast({ alert, onClose }) {
  // Auto close after 5 seconds
  useEffect(() => {
    if (alert) {
      const timer = setTimeout(() => {
        onClose();
      }, 5000);
      return () => clearTimeout(timer);
    }
  }, [alert, onClose]);

  if (!alert) return null;

  return (
    <div className="alert-toast-container">
      <div className="alert-toast alert-toast-danger">
        <div className="alert-toast-header">
          <span className="alert-icon">🔔</span>
          <strong>Noise Alert!</strong>
          <button 
            className="alert-close-btn" 
            onClick={onClose}
            aria-label="Close"
          >
            ×
          </button>
        </div>
        <div className="alert-toast-body">
          <div className="alert-detail">
            <span className="alert-label">Room:</span>
            <span className="alert-value">{alert.roomName}</span>
          </div>
          <div className="alert-detail">
            <span className="alert-label">Level:</span>
            <span className="alert-value alert-highlight">{alert.noiseLevel} dB</span>
          </div>
          <div className="alert-detail">
            <span className="alert-label">Time:</span>
            <span className="alert-value">{alert.time}</span>
          </div>
        </div>
      </div>
    </div>
  );
}

export default AlertToast;