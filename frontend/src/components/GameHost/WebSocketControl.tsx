import { Flex, Button, Table } from "@radix-ui/themes";
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
            addMessage("admin", "Aleady connected!");
            return;
        }

        try {
            setConnectionStatus('Connecting...');
            socketRef.current = new WebSocket('ws://localhost:8080/ws');

            socketRef.current.onopen = () => {
                setIsConnected(true);
                setConnectionStatus('Connected');
                addMessage("admin", "Connection established");
            };

            socketRef.current.onmessage = (event) => {
                addMessage("Received", `${event.data}`);
            };

            socketRef.current.onclose = () => {
                setIsConnected(false);
                setConnectionStatus('Disconnected');
                addMessage("admin", "Connection closed");
            };

            socketRef.current.onerror = (error) => {
                setConnectionStatus('Error');
                addMessage("Error", `${error.type}`);
            };
        } catch (error) {
            setConnectionStatus('Error');
            addMessage("Error", `Connection error: ${error instanceof Error ? error.message : String(error)}`);
        }
    };

    const disconnectWebSocket = () => {
        if (socketRef.current) {
            socketRef.current.close();
            socketRef.current = null;
        } else {
            addMessage("Error", "No active connection to close");
        }
    };

    const sendTestMessage = () => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            const message = `Hello server!`;
            socketRef.current.send(message);
            addMessage("Sent", `${message}`);
        } else {
            addMessage("Error", "Cannot send message: No connection");
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
                        <Table.Root>
                            <Table.Header>
                                <Table.Row>
                                    <Table.ColumnHeaderCell>Type</Table.ColumnHeaderCell>
                                    <Table.ColumnHeaderCell>Message</Table.ColumnHeaderCell>
                                    <Table.ColumnHeaderCell>Time</Table.ColumnHeaderCell>
                                </Table.Row>
                            </Table.Header>
                            <Table.Body>
                                {messages.map((msg) => (
                                    <Table.Row>
                                        <Table.Cell>{msg.type}</Table.Cell>
                                        <Table.Cell>{msg.message}</Table.Cell>
                                        <Table.Cell>{msg.time.toISOString()}</Table.Cell>
                                    </Table.Row>
                                ))}
                            </Table.Body>
                        </Table.Root>
                    )}
                </div>
            </div>
        </Flex>
    );
};