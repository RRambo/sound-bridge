import React, { useRef, useEffect } from 'react';
import '../styles/SettingsPanel.css';

function SettingsPanel({
  locations,
  chosenLocation,
  onLocationAdd,
  newLocationName,
  setNewLocationName,
  onLocationChange,
  onLocationDelete,
  isDropdownOpen,
  setIsDropdownOpen,
  threshold,
  onThresholdChange,
  onThresholdUpdate
}) {
  const dropdownRef = useRef(null);
  // Close dropdown when clicking outside of it
  useEffect(() => {
    function handleClickOutside(event) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsDropdownOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [setIsDropdownOpen]);

  return (
    <div className="settings-panel card">
      <div className="card-body">
        {/* Header */}
        <div className="settings-header mb-3">
          <h5 className="card-title">
            <span className="settings-icon">⚙️</span> Settings
          </h5>
          <p className="card-subtitle text-muted">
            Configure monitoring preferences
          </p>
        </div>

        {/* Add location form */}
        <div className="mb-4">
          <label className="form-label fw-semibold">Add New Room</label>
          <form onSubmit={onLocationAdd}>
            <input
              type="text"
              className="form-control d-inline-block"
              style={{ margin: '0 10px 10px 0' }}
              placeholder="New location name"
              value={newLocationName}
              onChange={(e) => setNewLocationName(e.target.value)}
              required
            />
            <button type="submit" className="btn btn-outline-primary w-100">
              + Add
            </button>
          </form>
        </div>

        {/* Listening Room Selector */}
        <div className="mb-4">
          <label className="form-label fw-semibold">Listening Room</label>
          <div className="" style={{ position: 'relative', margin: '0 auto' }} ref={dropdownRef}>
            <button
              className="form-select"
              onClick={() => setIsDropdownOpen(prev => !prev)}
              aria-expanded={isDropdownOpen}
              aria-haspopup="listbox"
            >
              {chosenLocation ? `${chosenLocation.name}` : 'Select location'}
            </button>

            {isDropdownOpen && (
              <ul
                className="list-group position-absolute w-100"
                style={{ zIndex: 50, maxHeight: '220px', overflow: 'auto' }}
                role="listbox"
              >
                {locations.length === 0 ? (
                  <li>No locations available</li>
                ) : (
                  locations.map(location => (
                    <li
                      key={location.id}
                      className="list-group-item d-flex justify-content-between align-items-center"
                      style={{ cursor: 'pointer' }}
                      onClick={() => {
                        onLocationChange(location.id)
                        setIsDropdownOpen(false);
                      }}
                    >
                      <span>{location.name}</span>
                      <button
                        className="btn btn-sm btn-outline-danger"
                        onClick={(e) => {
                          e.stopPropagation();
                          e.preventDefault();
                          onLocationDelete(location.id, location.name, e)
                        }}
                        title={`Delete ${location.name}`}
                      >
                        X
                      </button>
                    </li>
                  ))
                )}
              </ul>
            )}
          </div>
        </div>

        {/* Noise Threshold Slider */}
        <div className="mb-3">
          <div className="d-flex justify-content-between align-items-center mb-2">
            <label className="form-label fw-semibold mb-0">Noise Threshold</label>
            <span className="badge bg-light text-dark">{threshold} dB</span>
          </div>

          <input
            type="range"
            className="form-range"
            min="0"
            max="120"
            step="5"
            value={threshold}
            onChange={(e) => onThresholdChange(parseInt(e.target.value))}
            onMouseUp={(e) => onThresholdUpdate(parseInt(e.target.value))}
          />

          <p className="text-muted small mt-2 mb-0">
            You'll be notified when noise exceeds this level
          </p>
        </div>
      </div>
    </div>
  );
}

export default SettingsPanel;