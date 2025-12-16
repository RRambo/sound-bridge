import React from 'react';
import '../styles/AlertHistory.css';

function AlertHistory({ alerts }) {
  if (alerts.length === 0) {
    return (
      <div className="alert-history-card card">
        <div className="card-body">
          <div className="alert-history-header">
            <h5 className="card-title">
              <span className="history-icon">📋</span> Today's Alerts
            </h5>
            <span className="alert-count-badge">{alerts.length}</span>
          </div>
          <p className="text-muted text-center py-4">No alerts today</p>
        </div>
      </div>
    );
  }

  return (
    <div className="alert-history-card card">
      <div className="card-body">
        <div className="alert-history-header">
          <h5 className="card-title">
            <span className="history-icon">📋</span> Today's Alerts
          </h5>
          <span className="alert-count-badge">{alerts.length}</span>
        </div>

        <div className="alert-history-list">
          {alerts.map((alert, index) => (
            <div key={index} className="alert-history-item">
              <div className="alert-history-time">
                <span className="time-icon">🕐</span>
                {alert.time}
              </div>
              <div className="alert-history-content">
                <div className="alert-history-room">
                  <span className="room-icon">📍</span>
                  {alert.roomName}
                </div>
                <div className="alert-history-level">
                  <span className="level-badge">{alert.noiseLevel} dB</span>
                  <span className="exceeded-text">
                    (Threshold: {alert.threshold} dB)
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default AlertHistory;