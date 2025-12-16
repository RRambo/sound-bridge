import React from 'react';
import '../styles/NoiseMeter.css';

function NoiseMeter({ currentLevel, threshold, roomName }) {
  // Calculate percentage for the circular progress
  const percentage = Math.min((currentLevel / 120) * 100, 100);

  // Determine status based on level vs threshold
  const getStatus = () => {
    if (currentLevel < threshold * 0.7) {
      return { label: 'Quiet', color: '#5F9EA0', bgColor: '#E0F2F1' };
    } else if (currentLevel < threshold) {
      return { label: 'Moderate', color: '#FFA726', bgColor: '#FFF3E0' };
    } else {
      return { label: 'Loud', color: '#EF5350', bgColor: '#FFEBEE' };
    }
  };

  const status = getStatus();

  // SVG circle parameters
  const size = 280;
  const strokeWidth = 20;
  const radius = (size - strokeWidth) / 2;
  const circumference = 2 * Math.PI * radius;
  const strokeDashoffset = circumference - (percentage / 100) * circumference;

  return (
    <div className="noise-meter-card card">
      <div className="card-body">
        {/* Header */}
        <div className="d-flex justify-content-between align-items-center mb-4">
          <h5 className="card-title mb-0">Current Noise Level</h5>
          <span className="room-badge">{roomName}</span>
        </div>

        {/* Circular Meter */}
        <div className="meter-container">
          <svg width={size} height={size} className="circular-meter">
            {/* Background circle */}
            <circle
              cx={size / 2}
              cy={size / 2}
              r={radius}
              fill="none"
              stroke="#e0e0e0"
              strokeWidth={strokeWidth}
            />

            {/* Progress circle */}
            <circle
              cx={size / 2}
              cy={size / 2}
              r={radius}
              fill="none"
              stroke={status.color}
              strokeWidth={strokeWidth}
              strokeDasharray={circumference}
              strokeDashoffset={strokeDashoffset}
              strokeLinecap="round"
              transform={`rotate(-90 ${size / 2} ${size / 2})`}
              style={{ transition: 'stroke-dashoffset 0.5s ease' }}
            />

            {/* Threshold line */}
            <line
              x1={size / 2 + (radius - strokeWidth / 2) * Math.cos((threshold / 120) * 2 * Math.PI - Math.PI / 2)}
              y1={size / 2 + (radius - strokeWidth / 2) * Math.sin((threshold / 120) * 2 * Math.PI - Math.PI / 2)}
              x2={size / 2 + (radius + strokeWidth / 2) * Math.cos((threshold / 120) * 2 * Math.PI - Math.PI / 2)}
              y2={size / 2 + (radius + strokeWidth / 2) * Math.sin((threshold / 120) * 2 * Math.PI - Math.PI / 2)}
              stroke="#FF0000"
              strokeWidth="5"
            />
          </svg>

          {/* Center content */}
          <div className="meter-content">
            <div className="sound-icon">🔊</div>
            <div className="noise-level">{currentLevel}</div>
            <div className="noise-unit">dB</div>
          </div>

          {/* Status badge below meter */}
          <div className="status-badge-container">
            <span
              className="status-badge"
              style={{
                backgroundColor: status.bgColor,
                color: status.color
              }}
            >
              {status.label}
            </span>
          </div>
        </div>

        {/* Info Row 
        <div className="info-row">
          <div className="info-item">
            <div className="info-label">Threshold</div>
            <div className="info-value">{threshold} dB</div>
          </div>
          <div className="info-divider"></div>
          <div className="info-item">
            <div className="info-label">Status</div>
            <div className="info-value">{status.label}</div>
          </div>
          <div className="info-divider"></div>
          <div className="info-item">
            <div className="info-label">Level</div>
            <div className="info-value">{currentLevel} dB</div>
          </div>
        </div>
        */}
      </div>
    </div>
  );
}

export default NoiseMeter;