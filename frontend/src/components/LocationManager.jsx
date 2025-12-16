import React, { useState, useEffect } from 'react';
import { locationAPI } from '../services/api';

// This component is not in use now, but I'll leave it here for a while, in case.

function LocationManager() {
  // State management
  const [locations, setLocations] = useState([]);
  const [newLocationName, setNewLocationName] = useState('');
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [chosenLocation, setChosenLocation] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Load locations when component mounts
  useEffect(() => {
    loadLocations();
  }, []);

  // Fetch locations from backend
  const loadLocations = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await locationAPI.getAll();
      const locationData = response.data.locations || [];
      setLocations(locationData);

      // Find the chosen location
      const chosen = locationData.find(loc => loc.chosen);
      setChosenLocation(chosen);
    } catch (err) {
      console.error('Error loading locations:', err);
      setError('Failed to load locations. Make sure the server is running.');
    } finally {
      setLoading(false);
    }
  };

  // Add a new location
  const handleAddLocation = async (e) => {
    e.preventDefault();
    if (!newLocationName.trim()) return;

    try {
      await locationAPI.create(newLocationName.trim());
      setNewLocationName('');
      await loadLocations();
    } catch (err) {
      console.error('Error adding location:', err);
      alert('Failed to add location.');
    }
  };

  // Change chosen location
  const handleChooseLocation = async (locationId) => {
    try {
      await locationAPI.setChosen(locationId);
      setIsDropdownOpen(false);
      await loadLocations();
    } catch (err) {
      console.error('Error updating location:', err);
      alert('Failed to update location.');
    }
  };

  // Delete a location
  const handleDeleteLocation = async (id, name, e) => {
    e.stopPropagation();

    if (!window.confirm(`Are you sure you want to delete location (${name})?`)) {
      return;
    }

    try {
      await locationAPI.delete(id);
      await loadLocations();
    } catch (err) {
      console.error('Error deleting location:', err);
      alert('Failed to delete location.');
    }
  };

  // Render loading state
  if (loading) {
    return <div className="text-center">Loading locations...</div>;
  }

  // Render error state
  if (error) {
    return <div className="alert alert-danger">{error}</div>;
  }

  return (
    <div className="location-manager">
      <h1>Kindergarten Noise Meter</h1>

      {/* Add location form */}
      <div className="mb-3">
        <form onSubmit={handleAddLocation}>
          <input
            type="text"
            className="form-control d-inline-block"
            style={{ width: '200px', marginRight: '10px' }}
            placeholder="New location name"
            value={newLocationName}
            onChange={(e) => setNewLocationName(e.target.value)}
            required
          />
          <button type="submit" className="btn btn-primary">
            Add
          </button>
        </form>
      </div>

      {/* Custom dropdown */}
      <div className="dropdown" style={{ position: 'relative', width: '200px', margin: '0 auto' }}>
        <p>Chosen Location:</p>
        <button
          className="btn btn-secondary w-100"
          onClick={() => setIsDropdownOpen(!isDropdownOpen)}
        >
          {chosenLocation ? `${chosenLocation.name} ▾` : 'Select location ▾'}
        </button>

        {isDropdownOpen && (
          <ul
            className="list-group position-absolute w-100"
            style={{ zIndex: 50, maxHeight: '220px', overflow: 'auto' }}
          >
            {locations.map(location => (
              <li
                key={location.id}
                className="list-group-item d-flex justify-content-between align-items-center"
                style={{ cursor: 'pointer' }}
                onClick={() => handleChooseLocation(location.id)}
              >
                <span>{location.name}</span>
                <button
                  className="btn btn-sm btn-danger"
                  onClick={(e) => handleDeleteLocation(location.id, location.name, e)}
                  title={`Delete ${location.name}`}
                >
                  X
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>

      {/*<select
        className="form-select"
        value={chosenLocation?.id || ''}
        onChange={(e) => onLocationChange(parseInt(e.target.value))}
      >
        {locations.length === 0 ? (
          <option value="">No locations available</option>
        ) : (
          locations.map(location => (
            <option key={location.id} value={location.id}>
              {location.name}
            </option>
          ))
        )}
      </select>*/}
    </div>
  );
}

export default LocationManager;