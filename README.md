# Sound Bridge
Build version of the noise meter server 'Sound Bridge', deployed in a cloud service.
<br>You can test the service through this link: https://sound-bridge.onrender.com/
<br>Data can be posted to 'https://sound-bridge.onrender.com/api/data'
with the following headers
<br>**Basic Authentication**: 
<br>
username: kids_noisemeter_admin
<br>
password: passwordkids
<br>**Content-Type:** application/json
<br>
And a **data model** that looks like this:
{
    "sound_level": 67.45
}
## Data Model
Only the sound_level is required (a float between 0 and 150)
```json
{
  "id": 1,
  "device_id": "arduino_001",
  "room_name": "PlayRoom_A",
  "sound_level": 78.5,
  "threshold": 70.0,
  "measure_time": "2024-10-27T09:30:00Z",
  "is_alert": true,
  "description": "Morning playtime noise peak",
  "IsPeriodic": true
}
```

## License

Educational project for Intelligent Devices course.

---

**Last Updated:** 15/12/2025
