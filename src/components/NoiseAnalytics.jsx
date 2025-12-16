import React, { useState, useEffect } from 'react';
import { AreaChart, Area, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, LineChart, Line } from 'recharts';
import '../styles/NoiseAnalytics.css';
import ChartEmptyState from './ChartEmptyState';

function NoiseAnalytics({ roomName, allLocations, getDailySummary, setDailyPeak, getWeeklySummary, setWeeklyAverage }) {
    const [activeTab, setActiveTab] = useState('firstRender');
    const [selectedDate, setSelectedDate] = useState(new Date());
    const [selectedRoom, setSelectedRoom] = useState(roomName);
    const [weekOffset, setWeekOffset] = useState(0);
    const [chartData, setChartData] = useState([]);
    const [quietTimeData, setQuietTimeData] = useState([]);

    // Update selected room when roomName changes
    useEffect(() => {
        setSelectedRoom(roomName);
    }, [roomName]);

    // Generate chart data
    useEffect(() => {
        generateChartData();
    }, [activeTab, selectedDate, selectedRoom, weekOffset]);

    // Get which day of week from date
    const getDayOfWeek = (date) => {
        const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
        return days[date.getDay()];
    };

    // Get which day from chosen day of week
    const getDateForDay = (dayName) => {
        const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
        const targetDayIndex = days.indexOf(dayName);

        if (targetDayIndex === -1) return selectedDate;

        const today = new Date();
        const currentDayIndex = today.getDay();
        const diff = targetDayIndex - currentDayIndex;

        const newDate = new Date(today);
        newDate.setDate(today.getDate() + diff);

        return newDate;
    };

    // Which day of week is selected date
    const currentDayOfWeek = getDayOfWeek(selectedDate);

    // Handle date change
    const handleDateChange = (newDate) => {
        if (newDate === "") {
            return;
        } else {
            setSelectedDate(newDate);
        }
    };

    // Handle week selection change
    const handleDayChange = (dayName) => {
        const newDate = getDateForDay(dayName);
        setSelectedDate(newDate);
    };

    const generateChartData = async () => {
        // kinda laggy but works :/
        if (activeTab === 'firstRender') {
            // get weekly data of the last 5 weeks from a specific room.
            const weeklyData = await getWeeklySummary(selectedRoom) || []

            // Filter data for hours between 8 and 17
            const filteredWeeklyData = weeklyData.filter((d) => {
                const hour = new Date(d.measure_time).getHours();
                return hour >= 6 && hour <= 17;
            });

            // variables for gettign data for the weekly average statCard
            const now = new Date();
            const dayOfWeek = now.getDay(); // 0 = Sunday, 1 = Monday, ..., 6 = Saturday
            const monday = new Date(now);
            monday.setDate(now.getDate() - ((dayOfWeek + 6) % 7)); // Get Monday of this week
            monday.setHours(0, 0, 0, 0); // Set to midnight

            // Get data for the weekly average
            const averageData = filteredWeeklyData.filter((d) => {
                const measureDate = new Date(d.measure_time);
                return measureDate >= monday && measureDate <= now;
            });

            if (averageData.length > 0) {
                // Calculate and set weekly average statCard
                const levels = averageData.map(d => d.sound_level);
                const average = Math.round(levels.reduce((a, b) => a + b, 0) / levels.length);
                setWeeklyAverage(average);
            }
            setActiveTab('daily');
        }

        if (activeTab === 'daily' || activeTab === 'firstRender') {
            // get data from a specific room with measure_time during a specific day
            // (specific day = today on first page render)
            const dailyData = await getDailySummary(selectedRoom, selectedDate) || []
            // Check if there are data within the hours of 8 and 17
            const hasDataInWindow = dailyData.some(e => {
                const h = new Date(e.measure_time).getHours();
                return h >= 6 && h <= 17;
            });

            // check if selectedDate is today
            let isToday = false;
            const currentTime = new Date();

            if (selectedDate.getFullYear() === currentTime.getFullYear() &&
                selectedDate.getMonth() === currentTime.getMonth() &&
                selectedDate.getDate() === currentTime.getDate()
            ) isToday = true;

            // if there is no data between 8 and 17, exit with no chartData
            if (!hasDataInWindow) {
                setChartData([])
                if (isToday) {
                    setDailyPeak('-');
                    console.warn(`No data for daily peak statCard`)
                }
                return;
            }
            // if requested date is today, 
            // get the peak for the daily peak StatCard as well
            let highestPeak = 0;

            // calculates the average and looks for the peaks
            // Works while we wait for the hourly data implementation

            // set daily data (6:00-17:00)

            // Original implementation with fixed hours 8-17
            //const hours = Array.from({ length: 10 }, (_, i) => 8 + i);

            // Now dynamically removes hours with no data from the beginning and end of the chart
            // Find min and max hour in the data (between 6 and 17)
            const hourValues = dailyData
                .map(e => {
                    const match = e.measure_time.match(/T(\d{1,2}):/);
                    return match ? Number(match[1]) : null;
                })
                .filter(h => h !== null && h >= 6 && h <= 17);

            const minHour = Math.min(...hourValues);
            const maxHour = Math.max(...hourValues);

            // Generate a continuous range of hours
            const hours = [];
            for (let h = minHour; h <= maxHour; h++) {
                hours.push(h);
            }
            
            // buckets for data within each hour
            const buckets = Object.fromEntries(hours.map(h => [h, { count: 0, sum: 0, max: null }]));
            const data = [];
            // fill the buckets with data from the server
            dailyData.forEach(e => {
                // ignore timezone offset
                const match = e.measure_time.match(/T(\d{1,2}):/);
                const hour = match ? Number(match[1]) : NaN;

                // Only add the data within 8 - 17
                if (hours.includes(hour)) {
                    const b = buckets[hour];

                    b.count += 1;
                    b.sum += e.sound_level;
                    b.max = b.max === null ? e.sound_level : Math.max(b.max, e.sound_level);
                }
            });
            // Calculate the average and peak to the data array
            hours.forEach(h => {
                const { count, sum, max } = buckets[h];
                data.push({
                    time: `${h}:00`,
                    avgLevel: count > 0 ? Math.round(sum / count) : null,
                    peakLevel: max === null ? null : max,
                });

                // Look for the highest peak for the daily peak StatCard
                if (isToday && max > highestPeak) {
                    highestPeak = max;
                }
            });
            /*// for checking the data in the buckets
            data.forEach(e =>
                console.log(`time: ${e.time}, average: ${e.avgLevel}, peak: ${e.peakLevel}`)
            )*/
            setDailyPeak(highestPeak)
            setChartData(data);
        } else {
            // get weekly data of the last 5 weeks from a specific room.
            const weeklyData = await getWeeklySummary(selectedRoom) || []

            // Filter data for hours between 8 and 17
            const filteredWeeklyData = weeklyData.filter((d) => {
                const hour = new Date(d.measure_time).getHours();
                return hour >= 6 && hour <= 17;
            });

            if (filteredWeeklyData.length === 0) {
                setChartData([])
                return;
            }

            // Buckets for generating weekly bar chart data
            const days = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday'];
            const dayBuckets = Object.fromEntries(days.map(day => [day, []]));

            // variables for gettign data for the weekly average statCard
            const now = new Date();
            const dayOfWeek = now.getDay(); // 0 = Sunday, 1 = Monday, ..., 6 = Saturday
            const monday = new Date(now);
            monday.setDate(now.getDate() - ((dayOfWeek + 6) % 7)); // Get Monday of this week
            monday.setHours(0, 0, 0, 0); // Set to midnight

            // variables for generating chart data
            const chartMonday = new Date(now);
            chartMonday.setDate(now.getDate() - dayOfWeek + (dayOfWeek === 0 ? -6 : 1) + (weekOffset * 7));
            chartMonday.setHours(0, 0, 0, 0);

            const chartFriday = new Date(chartMonday);
            chartFriday.setDate(chartMonday.getDate() + 4);
            chartFriday.setHours(23, 59, 59, 999);

            // Get data for the current week
            // And group data by weekday accounting for weekOffset
            const averageData = filteredWeeklyData.filter((d) => {
                const measureDate = new Date(d.measure_time);
                // for the chart data
                if (measureDate >= chartMonday && measureDate <= chartFriday) {
                    const weekdayIndex = measureDate.getDay();
                    if (weekdayIndex >= 1 && weekdayIndex <= 5) {
                        const weekday = days[weekdayIndex - 1];
                        dayBuckets[weekday].push(d.sound_level);
                    }
                }
                // for the weekly average statCard
                return measureDate >= monday && measureDate <= now;
            });

            if (averageData.length > 0) {
                // Calculate and set weekly average statCard
                const levels = averageData.map(d => d.sound_level);
                const average = Math.round(levels.reduce((a, b) => a + b, 0) / levels.length);
                setWeeklyAverage(average);
            }

            const isWeekEmpty = Object.values(dayBuckets).every(levels => levels.length === 0);

            if (isWeekEmpty) {
                setChartData([])
                return;
            }

            // Continue to generate chart data
            // Calculate average and peak for each day
            const data = days.map(day => {
                const levels = dayBuckets[day];
                const avgNoise = levels.length > 0 ? Math.round(levels.reduce((a, b) => a + b, 0) / levels.length) : null;
                const peakNoise = levels.length > 0 ? Math.max(...levels) : null;
                return {
                    day: day,
                    avgNoise,
                    peakNoise
                };
            });
            setChartData(data);

            // Generate quiet time duration data
            // Calculates the quiet time based on the amount of measurements below 70 % of the threshold
            // Currently assumes that measurements are taken every 10 minutes
            //const QUIET_THRESHOLD = chosenThreshold * 0.7; // If threshold = 75, QUIET_THRESHOLD = 52.5
            const MEASUREMENT_INTERVAL_MINUTES = 10; // If each measurement is every 10 minutes

            const quietData = days.map(day => {
                const levels = dayBuckets[day];
                // Count how many measurements are below the quiet threshold (60 dB)
                const quietCount = levels.filter(level => level < 60).length;
                // Calculate total quiet minutes
                const duration = Math.round(quietCount * MEASUREMENT_INTERVAL_MINUTES);
                return {
                    day: day,
                    duration: duration
                };
            });
            setQuietTimeData(quietData);
        }
    };

    const getWeekDateRange = () => {
        const today = new Date();
        const currentDay = today.getDay(); // 0 = Sunday, 1 = Monday, etc.
        const monday = new Date(today);
        monday.setDate(today.getDate() - currentDay + (currentDay === 0 ? -6 : 1) + (weekOffset * 7));

        const friday = new Date(monday);
        friday.setDate(monday.getDate() + 4);

        const formatDate = (date) => {
            const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
            return `${months[date.getMonth()]} ${date.getDate()}`;
        };

        return `${formatDate(monday)} - ${formatDate(friday)}, ${monday.getFullYear()}`;
    };

    const getWeekLabel = () => {
        if (weekOffset === 0) return 'Current Week';
        if (weekOffset === -1) return 'Last Week';
        return `${Math.abs(weekOffset)} Weeks Ago`;
    };

    const canGoBack = weekOffset > -4;
    const canGoForward = weekOffset < 0;

    return (
        <div className="noise-analytics-card card">
            <div className="card-body">
                {/* Header */}
                <div className="analytics-header">
                    <div>
                        <h5 className="card-title mb-1">
                            <span className="icon">📊</span> Noise Analysis
                        </h5>
                        <p className="card-subtitle text-muted">
                            Track noise patterns to plan activities during quieter periods
                        </p>
                    </div>
                    <div className="room-selector">
                        <label className="form-label mb-1">Room</label>
                        <select
                            className="form-select form-select-sm"
                            value={selectedRoom}
                            onChange={(e) => setSelectedRoom(e.target.value)}
                        >
                            {allLocations && allLocations.length > 0 ? (
                                allLocations.map(loc => (
                                    <option key={loc.id} value={loc.name}>{loc.name}</option>
                                ))
                            ) : (
                                <option>{roomName}</option>
                            )}
                        </select>
                    </div>
                </div>

                {/* Tab Buttons */}
                <div className="tab-buttons">
                    <button
                        className={`tab-btn ${activeTab === 'daily' || activeTab === 'firstRender' ? 'active' : ''}`}
                        onClick={() => setActiveTab('daily')}
                    >
                        Daily Analysis
                    </button>
                    <button
                        className={`tab-btn ${activeTab === 'weekly' ? 'active' : ''}`}
                        onClick={() => setActiveTab('weekly')}
                    >
                        Weekly Analysis
                    </button>
                </div>

                {/* Selectors */}
                {activeTab === 'daily' || activeTab === 'firstRender' ? (
                    <div className="selectors">
                        <div className="selector-item">
                            <label>Date</label>
                            <input
                                type="date"
                                className="form-control form-control-sm"
                                value={selectedDate.toISOString().split('T')[0]}
                                onChange={(e) => {
                                    // prevent error when clearing date
                                    e.target.value === "" ? handleDateChange(new Date()) : handleDateChange(new Date(e.target.value))
                                }}
                            />
                        </div>
                        <div className="selector-item">
                            <label>Choose Day from This Week</label>
                            <select
                                className="form-select form-select-sm"
                                value={currentDayOfWeek}
                                onChange={(e) => handleDayChange(e.target.value)}
                            >
                                {['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday'].map(day => (
                                    <option key={day} value={day}>{day}</option>
                                ))}
                            </select>
                        </div>
                    </div>
                ) : (
                    <div className="week-selector">
                        <div className="week-label">This Week</div>
                        <div className="week-navigation">
                            <button
                                className="nav-btn"
                                onClick={() => setWeekOffset(prev => prev - 1)}
                                disabled={!canGoBack}
                            >
                                &#8249;
                            </button>
                            <span className="week-text">{getWeekLabel()}</span>
                            <button
                                className="nav-btn"
                                onClick={() => setWeekOffset(prev => prev + 1)}
                                disabled={!canGoForward}
                            >
                                &#8250;
                            </button>
                        </div>
                        <div className="week-range">{getWeekDateRange()}</div>
                    </div>
                )}

                {/* Chart */}
                <div className="chart-container">
                    {activeTab === 'daily' || activeTab === 'firstRender' ? (
                        chartData.length === 0 ? (
                            // Overkill of an empty chart state XD
                            <ChartEmptyState
                                chosenLabel={selectedDate.toLocaleDateString()}
                                mode="daily"
                                onRetry={() => handleDateChange(new Date(selectedDate))}
                            />
                        ) : (
                            <ResponsiveContainer width="100%" height={300}>
                                <AreaChart data={chartData} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                                    <defs>
                                        <linearGradient id="colorAvg" x1="0" y1="0" x2="0" y2="1">
                                            <stop offset="5%" stopColor="#5F9EA0" stopOpacity={0.3} />
                                            <stop offset="95%" stopColor="#5F9EA0" stopOpacity={0.05} />
                                        </linearGradient>
                                    </defs>
                                    <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
                                    <XAxis
                                        dataKey="time"
                                        stroke="#666"
                                        style={{ fontSize: '0.85rem' }}
                                    />
                                    <YAxis
                                        stroke="#666"
                                        style={{ fontSize: '0.85rem' }}
                                        domain={[0, 100]}
                                        label={{ value: 'Decibels (dB)', angle: -90, position: 'insideLeft', style: { fontSize: '0.85rem' } }}
                                    />
                                    <Tooltip
                                        contentStyle={{
                                            backgroundColor: 'white',
                                            border: '1px solid #ccc',
                                            borderRadius: '8px',
                                            padding: '10px'
                                        }}
                                    />
                                    <Area
                                        type="monotone"
                                        dataKey="avgLevel"
                                        stroke="#5F9EA0"
                                        strokeWidth={2}
                                        fill="url(#colorAvg)"
                                        name="Average Level"
                                    />
                                    <Area
                                        type="monotone"
                                        dataKey="peakLevel"
                                        stroke="#5F9EA0"
                                        strokeWidth={2}
                                        strokeDasharray="5 5"
                                        fill="none"
                                        name="Peak Level"
                                    />
                                </AreaChart>
                            </ResponsiveContainer>
                        )
                    ) : (
                        chartData.length === 0 ? (
                            // Overkill of an empty chart state XD
                            <ChartEmptyState
                                chosenLabel={getWeekDateRange()}
                                mode={"weekly"}
                                onRetry={() => generateChartData()}
                            />
                        ) : (
                            <>
                                <ResponsiveContainer width="100%" height={300}>
                                    <BarChart data={chartData} margin={{ top: 20, right: 30, left: 0, bottom: 5 }}>
                                        <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
                                        <XAxis
                                            dataKey="day"
                                            stroke="#666"
                                            style={{ fontSize: '0.85rem' }}
                                        />
                                        <YAxis
                                            stroke="#666"
                                            style={{ fontSize: '0.85rem' }}
                                            domain={[0, 100]}
                                            label={{ value: 'Decibels (dB)', angle: -90, position: 'insideLeft', style: { fontSize: '0.85rem' } }}
                                        />
                                        <Tooltip
                                            contentStyle={{
                                                backgroundColor: 'white',
                                                border: '1px solid #ccc',
                                                borderRadius: '8px',
                                                padding: '10px'
                                            }}
                                        />
                                        <Legend
                                            wrapperStyle={{ paddingTop: '10px' }}
                                            iconType="rect"
                                        />
                                        <Bar dataKey="avgNoise" fill="#81C9CC" name="Average Noise" />
                                        <Bar dataKey="peakNoise" fill="#5F9EA0" name="Peak Noise" />
                                    </BarChart>
                                </ResponsiveContainer>

                                {/* Quiet Time Duration Chart */}
                                <div className="quiet-time-section">
                                    <h6 className="section-title">
                                        <span className="icon">🔇</span> Quiet Time Duration (minutes)
                                    </h6>
                                    <ResponsiveContainer width="100%" height={200}>
                                        <LineChart data={quietTimeData} margin={{ top: 10, right: 30, left: 0, bottom: 5 }}>
                                            <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
                                            <XAxis
                                                dataKey="day"
                                                stroke="#666"
                                                style={{ fontSize: '0.85rem' }}
                                            />
                                            <YAxis
                                                stroke="#666"
                                                style={{ fontSize: '0.85rem' }}
                                                domain={[0, 300]}
                                            />
                                            <Tooltip
                                                contentStyle={{
                                                    backgroundColor: 'white',
                                                    border: '1px solid #ccc',
                                                    borderRadius: '8px',
                                                    padding: '10px'
                                                }}
                                            />
                                            <Line
                                                type="monotone"
                                                dataKey="duration"
                                                stroke="#5F9EA0"
                                                strokeWidth={2}
                                                dot={{ fill: '#5F9EA0', r: 4 }}
                                            />
                                        </LineChart>
                                    </ResponsiveContainer>
                                </div>
                            </>
                        )
                    )}
                </div>

                {/* Legend for Daily */}
                {activeTab === 'daily' && (
                    <div className="chart-legend">
                        <div className="legend-item">
                            <span className="legend-line solid"></span>
                            <span>Average Level</span>
                        </div>
                        <div className="legend-item">
                            <span className="legend-line dashed"></span>
                            <span>Peak Level</span>
                        </div>
                    </div>
                )}
                {activeTab === 'firstRender' && (
                    <div className="chart-legend">
                        <div className="legend-item">
                            <span className="legend-line solid"></span>
                            <span>Average Level</span>
                        </div>
                        <div className="legend-item">
                            <span className="legend-line dashed"></span>
                            <span>Peak Level</span>
                        </div>
                    </div>
                )}

                {/* Planning Tip */}
                <div className="planning-tip">
                    <strong>💡 Planning Tip:</strong> {activeTab === 'daily'
                        ? 'Schedule quiet activities (story time, nap time) during naturally quieter periods and active play during peak energy times.'
                        : 'The day showing the longest quiet periods and lowest average noise usually is ideal for introducing new concepts or activities requiring focus. The day with higher energy levels is perfect for group activities and celebrations.'}
                    {/* put template tips here */}
                </div>
            </div>
        </div>
    );
}

export default NoiseAnalytics;