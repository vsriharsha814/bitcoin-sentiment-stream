import React, { useState, useEffect } from 'react';
import {
    LineChart,
    Line,
    XAxis,
    YAxis,
    Tooltip,
    Legend,
    CartesianGrid,
    ResponsiveContainer,
} from 'recharts';
import { DatePicker, message, Checkbox, Typography } from 'antd';
import dayjs from 'dayjs';
import './SentimentChart.css';

const { RangePicker } = DatePicker;
const { Title } = Typography;

const allCoins = ['Bitcoin', 'Ethereum', 'Solana', 'Cardano'];

type SentimentPoint = {
    time: string;
    timestamp: number;
    [coin: string]: number | string;
};

const generateMockData = (
    start: Date,
    end: Date,
    intervalMinutes: number,
    selectedCoins: string[]
): SentimentPoint[] => {
    const data: SentimentPoint[] = [];
    const current = new Date(start);

    while (current <= end) {
        const point: SentimentPoint = {
            time: dayjs(current).format('MM/DD HH:mm'),
            timestamp: current.getTime(),
        };

        selectedCoins.forEach((coin) => {
            point[coin] = parseFloat((Math.random() * 2 - 1).toFixed(2));
        });

        data.push(point);
        current.setMinutes(current.getMinutes() + intervalMinutes);
    }

    return data;
};

const fetchSentimentData = async (
    start: Date,
    end: Date,
    coins: string[]
): Promise<SentimentPoint[]> => {
    const rangeMinutes = (end.getTime() - start.getTime()) / 60000;

    if (rangeMinutes > 180) {
        message.error('Please select a time range less than or equal to 2 hours.');
        return [];
    }
    const payload = {
        startTime: start.toISOString(),
        endTime: end.toISOString(),
        coins,
    };
    // TODO: Replace this with a real API call when backend is ready
    /*


    try {
      const response = await fetch('/api/sentiment', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('API fetch error:', error);
      message.error('Failed to fetch sentiment data.');
      return [];
    }
    */

    console.log( 'payload:', payload);
    return generateMockData(start, end, 5, coins);
};

const SentimentChart: React.FC = () => {
    const [timeRange, setTimeRange] = useState<[dayjs.Dayjs, dayjs.Dayjs]>(() => {
        const end = dayjs();
        const start = end.subtract(30, 'minute');
        return [start, end];
    });

    const [selectedCoins, setSelectedCoins] = useState<string[]>(allCoins);
    const [chartData, setChartData] = useState<SentimentPoint[]>([]);

    useEffect(() => {
        const fetchData = async () => {
            const data = await fetchSentimentData(
                timeRange[0].toDate(),
                timeRange[1].toDate(),
                selectedCoins
            );
            setChartData(data);
        };
        fetchData();
    }, [timeRange, selectedCoins]);

    const handleTimeChange = (
        values: [dayjs.Dayjs | null, dayjs.Dayjs | null] | null
    ) => {
        if (values && values[0] && values[1]) {
            const diffMinutes = values[1].diff(values[0], 'minute');
            if (diffMinutes > 120) {
                message.error('Please select a time range less than or equal to 2 hours.');
                return;
            }

            setTimeRange([values[0], values[1]]);
        }
    };


    const handleCoinChange = (checkedValues: string[]) => {
        setSelectedCoins(checkedValues);
    };

    return (
        <div className="chart-container">
            <div className="chart-controls">
                <Title level={3}>Crypto Sentiment Chart</Title>
                <div className="control-group">
                    <label>Start Time:</label>
                    <DatePicker
                        showTime
                        value={timeRange[0]}
                        format="MM/DD HH:mm"
                        onChange={(value) => {
                            if (value) {
                                const maxEnd = value.add(3, 'hour');
                                let newEnd = timeRange[1];

                                if (!newEnd || newEnd.isBefore(value) || newEnd.isAfter(maxEnd)) {
                                    newEnd = maxEnd;
                                }

                                setTimeRange([value, newEnd]);
                            }
                        }}
                    />
                </div>

                <div className="control-group">
                    <label>End Time (max 3 hours after start):</label>
                    <DatePicker
                        showTime
                        value={timeRange[1]}
                        format="MM/DD HH:mm"
                        disabledDate={(current) => {
                            const start = timeRange[0];
                            return current.isBefore(start) || current.isAfter(start.add(3, 'hour'));
                        }}
                        disabledTime={(date) => {
                            const start = timeRange[0];
                            const max = start.add(3, 'hour');
                            if (!date) return {};

                            if (!date.isSame(max, 'day')) return {};

                            return {
                                disabledHours: () => {
                                    const hours: number[] = [];
                                    for (let h = max.hour() + 1; h < 24; h++) {
                                        hours.push(h);
                                    }
                                    return hours;
                                },
                                disabledMinutes: (selectedHour) => {
                                    if (selectedHour === max.hour()) {
                                        const mins: number[] = [];
                                        for (let m = max.minute() + 1; m < 60; m++) {
                                            mins.push(m);
                                        }
                                        return mins;
                                    }
                                    return [];
                                },
                            };
                        }}
                        onChange={(value) => {
                            if (value) {
                                const duration = value.diff(timeRange[0], 'minute');
                                if (duration > 180) {
                                    message.error('Range must be 3 hours or less.');
                                } else {
                                    setTimeRange([timeRange[0], value]);
                                }
                            }
                        }}
                    />
                </div>


                <div className="control-group">
                    <label>Select Coins:</label>
                    <Checkbox.Group
                        options={allCoins}
                        value={selectedCoins}
                        onChange={handleCoinChange}
                    />
                </div>
            </div>

            <ResponsiveContainer width="100%" height={400}>
                <LineChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3"/>
                    <XAxis dataKey="time"/>
                    <YAxis domain={[-1, 1]} tickFormatter={(v) => v.toFixed(1)}/>
                    <Tooltip/>
                    <Legend/>
                    {selectedCoins.map((coin, index) => (
                        <Line
                            key={coin}
                            type="monotone"
                            dataKey={coin}
                            strokeWidth={2}
                            stroke={`hsl(${index * 90}, 70%, 50%)`}
                            dot={{r: 3}}
                        />
                    ))}
                </LineChart>
            </ResponsiveContainer>
        </div>
    );
};

export default SentimentChart;
