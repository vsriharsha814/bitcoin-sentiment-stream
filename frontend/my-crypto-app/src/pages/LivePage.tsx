import React, { useEffect, useRef, useState } from 'react';
import {
    LineChart, Line, XAxis, YAxis, Tooltip, Legend, CartesianGrid, ResponsiveContainer,
} from 'recharts';
import { Checkbox, Typography, message } from 'antd';
import dayjs from 'dayjs';
import './SentimentChart.css';

const { Title } = Typography;

const SOCKET_URL = 'ws://localhost:8080'; // Use your local or real backend
const allCoins = ['Bitcoin', 'Ethereum', 'Tether', 'XRP', 'BNB', 'Solana', 'USD Coin', 'TRON', 'Dogecoin', 'Cardano'];
const coinColors = [
    '#e6194b', '#3cb44b', '#ffe119', '#0082c8', '#f58231',
    '#911eb4', '#46f0f0', '#f032e6', '#d2f53c', '#fabebe'
];

type SentimentPoint = {
    time: string;
    [coin: string]: number | string;
};

const LivePage: React.FC = () => {
    const [selectedCoins, setSelectedCoins] = useState<string[]>(allCoins.slice(0, 5));
    const [data, setData] = useState<SentimentPoint[]>([]);
    const socketRef = useRef<WebSocket | null>(null);

    // Connect WebSocket
    useEffect(() => {
        const ws = new WebSocket(SOCKET_URL);
        socketRef.current = ws;

        ws.onopen = () => {
            console.log('WebSocket connected');
            ws.send(JSON.stringify({ coins: selectedCoins }));
        };

        ws.onmessage = (event) => {
            const incoming = JSON.parse(event.data);
            const formatted: SentimentPoint = {
                time: dayjs(incoming.time).format('DD-MMMM HH:mm'),
                ...incoming.sentiment,
            };
            setData((prev) => [...prev.slice(-35), formatted]); // Keep latest ~35 records
        };

        ws.onerror = (err) => {
            console.error('WebSocket error:', err);
            message.error('WebSocket connection failed.');
        };

        ws.onclose = () => {
            console.log('WebSocket disconnected');
        };

        return () => {
            ws.close();
        };
    }, []);

    // Send new coin subscription when selection changes
    useEffect(() => {
        if (socketRef.current?.readyState === WebSocket.OPEN) {
            socketRef.current.send(JSON.stringify({ coins: selectedCoins }));
        }
    }, [selectedCoins]);

    const handleCoinChange = (values: string[]) => {
        setSelectedCoins(values);
    };

    return (
        <div className="chart-container">
            <Title level={3}>Live Crypto Sentiment</Title>
            <div className="control-group">
                <label>Select Coins:</label>
                <Checkbox.Group
                    options={allCoins}
                    value={selectedCoins}
                    onChange={handleCoinChange}
                    style={{ marginBottom: 16 }}
                />
            </div>

            <ResponsiveContainer width="100%" height={400}>
                <LineChart data={data}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="time" />
                    <YAxis domain={[-1, 1]} />
                    <Tooltip />
                    <Legend />
                    {selectedCoins.map((coin, index) => (
                        <Line
                            key={coin}
                            type="monotone"
                            dataKey={coin}
                            stroke={coinColors[index % coinColors.length]}
                            strokeWidth={2}
                            dot={false}
                        />
                    ))}
                </LineChart>
            </ResponsiveContainer>
        </div>
    );
};

export default LivePage;
