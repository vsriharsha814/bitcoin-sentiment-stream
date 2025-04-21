import React, { useState } from 'react';
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
import {
    DatePicker,
    message,
    Checkbox,
    Typography,
    Spin,
    Button,
} from 'antd';
import dayjs from 'dayjs';
import './SentimentChart.css';

const { Title } = Typography;

const allCoins = [
    'Bitcoin', 'Ethereum', 'Tether', 'XRP', 'BNB',
    'Solana', 'USD Coin', 'TRON', 'Dogecoin', 'Cardano'
];

const coinColors = [
    '#e6194b', '#3cb44b', '#ffe119', '#0082c8', '#f58231',
    '#911eb4', '#46f0f0', '#f032e6', '#d2f53c', '#fabebe'
];

type SentimentPoint = {
    time: string;
    sentiment: Record<string, number>;
    title: Record<string, string>;
};

const generateMockData = (
    start: Date,
    end: Date,
    intervalMinutes: number,
    selectedCoins: string[]
): SentimentPoint[] => {
    const data: SentimentPoint[] = [];
    const current = new Date(start);

    const mockTitles = [
        "Strong investor interest pushes price higher",
        "Bearish trend due to market uncertainty",
        "Major exchange lists the token",
        "Whale activity detected",
        "Analyst upgrades forecast",
        "Market correction in progress",
        "Protocol upgrade announced",
        "Partnership news drives optimism",
        "Legal clarity improves sentiment",
        "Community votes on governance proposal"
    ];

    while (current <= end) {
        const point: SentimentPoint = {
            time: dayjs(current).format('DD-MMMM-YYYY HH:mm'),
            sentiment: {},
            title: {}
        };

        selectedCoins.forEach((coin) => {
            const sentimentScore = parseFloat((Math.random() * 2 - 1).toFixed(2));
            const randomTitle = mockTitles[Math.floor(Math.random() * mockTitles.length)];
            point.sentiment[coin] = sentimentScore;
            point.title[coin] = randomTitle;
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

    console.log('payload:', payload);

    return new Promise((resolve) =>
        setTimeout(() => {
            resolve(generateMockData(start, end, 5, coins));
        }, 1000)
    );
};

const SentimentChart: React.FC = () => {
    const [timeRange, setTimeRange] = useState<[dayjs.Dayjs, dayjs.Dayjs]>(() => {
        const end = dayjs();
        const start = end.subtract(30, 'minute');
        return [start, end];
    });

    const [selectedCoins, setSelectedCoins] = useState<string[]>(allCoins);
    const [chartData, setChartData] = useState<SentimentPoint[]>([]);
    const [loading, setLoading] = useState(false);

    const fetchAndUpdateChart = async () => {
        setLoading(true);
        try {
            const data = await fetchSentimentData(
                timeRange[0].toDate(),
                timeRange[1].toDate(),
                selectedCoins
            );
            setChartData(data);
            message.success('Sentiment data loaded');
        } catch (error) {
            message.error('Failed to fetch sentiment data.');
        } finally {
            setLoading(false);
        }
    };

    const handleCoinChange = (checkedValues: string[]) => {
        setSelectedCoins(checkedValues);
    };

    return (
        <div className="chart-container">
            <Spin spinning={loading} tip="Loading sentiment data...">
                <div className="chart-controls">
                    <Title level={3}>Crypto Sentiment Chart</Title>

                    <div className="date-submit-row">
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

                        <div className="submit-button-wrapper">
                            <Button type="primary" onClick={fetchAndUpdateChart}>
                                Submit
                            </Button>
                        </div>
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
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis
                            dataKey="time"
                            tickFormatter={(v) =>
                                dayjs(v, 'DD-MMMM-YYYY HH:mm').format('DD-MMMM HH:mm')
                            }
                            interval="preserveStartEnd"
                        />
                        <YAxis domain={[-1, 1]} tickFormatter={(v) => v.toFixed(1)} />
                        <Tooltip
                            wrapperStyle={{ pointerEvents: 'auto', maxHeight: 300, overflowY: 'auto' }}
                            content={({ active, payload, label }) => {
                                if (!active || !payload || payload.length === 0) return null;

                                return (
                                    <div className="custom-tooltip">
                                        <div className="custom-tooltip-header">{label}</div>
                                        <div className="custom-tooltip-body">
                                            {payload.map((entry) => {
                                                const coin = (entry.dataKey as string).split('.')[1];
                                                const value = entry.value;
                                                const fullData = entry.payload as any;
                                                const reason = fullData.title?.[coin];

                                                return (
                                                    <div key={coin} className="custom-tooltip-item">
                                                        <strong>{coin}: {Number(value).toFixed(2)}</strong>
                                                        {reason && <small>{reason}</small>}
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    </div>
                                );
                            }}
                        />
                        <Legend
                            formatter={(value) => {
                                const coin = value.split('.')[1];
                                return <span style={{ color: '#333' }}>{coin}</span>;
                            }}
                        />
                        {selectedCoins.map((coin, index) => (
                            <Line
                                key={coin}
                                type="monotone"
                                dataKey={`sentiment.${coin}`}
                                strokeWidth={2}
                                stroke={coinColors[index % coinColors.length]}
                                dot={{ r: 3 }}
                            />
                        ))}
                    </LineChart>
                </ResponsiveContainer>
            </Spin>
        </div>
    );
};

export default SentimentChart;
