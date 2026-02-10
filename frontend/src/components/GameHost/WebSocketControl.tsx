import { Flex, Button } from "@radix-ui/themes";
import { useEffect, useRef, useState } from "react";

interface message {
    type: string;
    message: string;
    time: Date;
}


export function WebSocketControl() {
    const [isConnected, setIsConnected] = useState(false);
    const [messages, setMessages] = useState<message[]>([]);
    const [connectionStatus, setConnectionStatus] = useState('Disconnected');
    const socketRef = useRef<WebSocket | null>(null);

    const addMessage = (type: string, message: string) => {
        setMessages(prev => [...prev,
        { type: type, message: message, time: new Date() },
        ]);
    };

    // Clean up the WebSocket connection when component unmounts
    useEffect(() => {
        return () => {
            if (socketRef.current) {
                socketRef.current.close();
            }
        };
    }, []);

    const connectWebSocket = () => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            setMessages(prev => [...prev, "Already connected!"]);
            return;
        }

        try {
            setConnectionStatus('Connecting...');
            socketRef.current = new WebSocket('ws://localhost:8080/ws');

            socketRef.current.onopen = () => {
                setIsConnected(true);
                setConnectionStatus('Connected');
                setMessages(prev => [...prev, "Connection established"]);
            };

            socketRef.current.onmessage = (event) => {
                setMessages(prev => [...prev, `Received: ${event.data}`]);
            };

            socketRef.current.onclose = () => {
                setIsConnected(false);
                setConnectionStatus('Disconnected');
                setMessages(prev => [...prev, "Connection closed"]);
            };

            socketRef.current.onerror = (error) => {
                setConnectionStatus('Error');
                setMessages(prev => [...prev, `Error: ${error.type}`]);
            };
        } catch (error) {
            setConnectionStatus('Error');
            setMessages(prev => [...prev, `Connection error: ${error instanceof Error ? error.message : String(error)}`]);
        }
    };

    const disconnectWebSocket = () => {
        if (socketRef.current) {
            socketRef.current.close();
            socketRef.current = null;
        } else {
            setMessages(prev => [...prev, "No active connection to close"]);
        }
    };

    const sendTestMessage = () => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            const message = `Hello server! Time: ${new Date().toLocaleTimeString()}`;
            socketRef.current.send(message);
            setMessages(prev => [...prev, `Sent: ${message}`]);
        } else {
            setMessages(prev => [...prev, "Cannot send message: No connection"]);
        }
    };

    return (
        <Flex className="websocket-control" align="center" justify="center" direction="column" gap="3">
            <h2>WebSocket Connection</h2>
            <div className="status-indicator">
                Status: <span className={`status-${connectionStatus.toLowerCase()}`}>{connectionStatus}</span>
            </div>

            <Flex className="button-container" direction="row" gap="3">
                <Button
                    onClick={connectWebSocket}
                    disabled={isConnected}
                    className={isConnected ? "button-disabled" : "button-connect"}
                >
                    Connect
                </Button>

                <Button
                    onClick={disconnectWebSocket}
                    disabled={!isConnected}
                    className={!isConnected ? "button-disabled" : "button-disconnect"}
                >
                    Disconnect
                </Button>

                <Button
                    onClick={sendTestMessage}
                    disabled={!isConnected}
                    className={!isConnected ? "button-disabled" : "button-send"}
                >
                    Send Test Message
                </Button>
            </Flex>

            <div className="message-log">
                <h3>Message Log</h3>
                <div className="messages">
                    {messages.length === 0 ? (
                        <p className="no-messages">No messages yet</p>
                    ) : (
                        messages.map((msg, index) => (
                            <div key={index} className="message">
                                {msg}
                            </div>
                        ))
                    )}
                </div>
            </div>
        </Flex>
    );
};