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
type CheckboxValueType = string | number;

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

const SentimentChart: React.FC = () => {
    const [timeRange, setTimeRange] = useState<[dayjs.Dayjs, dayjs.Dayjs]>(() => {
        const end = dayjs();
        const start = end.subtract(30, 'minute');
        return [start, end];
    });

    const [selectedCoins, setSelectedCoins] = useState<string[]>(allCoins);
    const [chartData, setChartData] = useState<SentimentPoint[]>([]);

    const updateChart = (start: Date, end: Date, coins: string[]) => {
        const rangeMinutes = (end.getTime() - start.getTime()) / 60000;
        if (rangeMinutes > 120) {
            message.error('Please select a time range less than or equal to 2 hours.');
            return;
        }
        const data = generateMockData(start, end, 5, coins);
        setChartData(data);
    };

    useEffect(() => {
        updateChart(timeRange[0].toDate(), timeRange[1].toDate(), selectedCoins);
    }, [timeRange, selectedCoins]);

    const handleTimeChange = (
        values: [dayjs.Dayjs | null, dayjs.Dayjs | null] | null
    ) => {
        if (
            values &&
            values[0] !== null &&
            values[1] !== null
        ) {
            setTimeRange([values[0], values[1]]);
        }
    };


    const handleCoinChange = (checkedValues: string[]) => {
        setSelectedCoins(checkedValues as string[]);
    };

    return (
        <div className="chart-container">
            <div className="chart-controls">
                <Title level={3}>Crypto Sentiment Chart</Title>
                <div className="control-group">
                    <label>Date Range:</label>
                    <RangePicker
                        showTime
                        format="MM/DD HH:mm"
                        value={timeRange}
                        onChange={handleTimeChange}
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
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="time" />
                    <YAxis domain={[-1, 1]} tickFormatter={(v) => v.toFixed(1)} />
                    <Tooltip />
                    <Legend />
                    {selectedCoins.map((coin, index) => (
                        <Line
                            key={coin}
                            type="monotone"
                            dataKey={coin}
                            strokeWidth={2}
                            stroke={`hsl(${index * 90}, 70%, 50%)`}
                            dot={{ r: 3 }}
                        />
                    ))}
                </LineChart>
            </ResponsiveContainer>
        </div>
    );
};

export default SentimentChart;
