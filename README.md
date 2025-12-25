# Sound Bridge
Build version of the noise meter server 'Sound Bridge', deployed in a cloud service.

You can test the service through this link: https://sound-bridge.onrender.com/

Data can be posted to 'https://sound-bridge.onrender.com/api/data'
with the following headers
<br>**Basic Authentication**: 
<br>
username: kids_noisemeter_admin
<br>
password: passwordkids
<br>**Content-Type:** application/json

And a **data model** that looks like this:
```json
{
    "sound_level": 67.45,
    "IsPeriodic": true
}
```
## Data Model
The IsPeriodic field determines whether the data is for the charts or the current noise level at the top of the UI (true=chart, false=current noise level).
<br>Only the sound_level is required (a float between 0 and 150)
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

**Last Updated:** 25/12/2025
