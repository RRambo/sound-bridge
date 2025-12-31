# Sound Bridge
This repository only contains a **build version** of the [Noise Meter Server](https://github.com/RRambo/Noise-meter-server/) 'Sound Bridge', deployed in a cloud service.

This web service is currently hosted in render and anyone can test it here [sound-bridge.onrender.com](https://sound-bridge.onrender.com/)

> The render instance used here is free and will spin down with inactivity, which can delay requests by 50+ seconds. 

Data can be posted to 'https://sound-bridge.onrender.com/api/data'
with the following headers  
```
Basic Authentication:
username: kids_noisemeter_admin
password: passwordkids

Content-Type: application/json
```
And a **request body** that looks like this:
```json
{
    "sound_level": 67.45,
    "IsPeriodic": true
}
```
## Data Model
The **IsPeriodic** field determines whether the data is for the charts or the current noise level.<br> `true`=chart, `false`=current noise level  
Only the sound_level is required (a float between 0 and 150 `xx.xx`)
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

## Images of the UI (with data)
![Screenshot of the User Interface with data including the daily summary.](<images/UI_Full.png>)

![Screenshot of the User Interface with data including the weekly summary and quiet time.](<images/UI_WeeklySummary.png>)

## License

Educational project for Intelligent Devices course.

---

**Last Updated:** 31/12/2025
