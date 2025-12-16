import React from 'react';
import '../styles/StatsCards.css';

function StatsCards({ dailyPeak, weeklyAverage, monitoringRoom /*isActive*/ }) {
  return (
    <div className="row g-4 mb-4">
      {/* Daily Peak Card */}
      <div className="col-md-4">
        <div className="stat-card card h-100">
          <div className="card-body">
            <div className="stat-header">
              <span className="stat-icon">📈</span>
              <span className="stat-label">Daily Peak</span>
            </div>
            <div className="stat-value">{dailyPeak} dB</div>
            <div className="stat-description">Highest level today</div>
          </div>
        </div>
      </div>

      {/* Weekly Average Card */}
      <div className="col-md-4">
        <div className="stat-card card h-100">
          <div className="card-body">
            <div className="stat-header">
              <span className="stat-icon">📊</span>
              <span className="stat-label">Weekly Average</span>
            </div>
            <div className="stat-value">{weeklyAverage} dB</div>
            <div className="stat-description">This week's average</div>
          </div>
        </div>
      </div>

      {/* Monitoring Status Card */}
      <div className="col-md-4">
        <div className="stat-card card h-100">
          <div className="card-body">
            <div className="stat-header">
              <span className="stat-icon">📡</span>
              <span className="stat-label">Monitoring</span>
            </div>
            <div className="stat-value-room">{monitoringRoom}</div>
            {/* <div className="stat-description">
              <span 
                className={`status-indicator ${isActive ? 'active' : 'inactive'}`}
              >
                {isActive ? '● Currently active' : '○ Inactive'}
              </span>
            </div> */}
          </div>
        </div>
      </div>
    </div>
  );
}

export default StatsCards;