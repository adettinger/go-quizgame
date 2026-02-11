import { Flex, Button, Table, TextField } from "@radix-ui/themes";
import { useEffect, useRef, useState } from "react";

enum protocolType {
    Admin = "Admin",
    Sent = "Sent",
    Received = "Received",
    Error = "Error",
}

enum messageType {
    MessageTypeAdmin = "admin",
    MessageTypeSent = "sent",
    MessageTypeChat = "chat",
    MessageTypeJoin = "join",
    MessageTypeLeave = "leave",
    MessageTypeGameUpdate = "game_update",
    MessageTypeError = "error",
}


interface WebSocketMessage {
    type: messageType;
    timestamp: Date;
    playerName: string;
    content: string;
}


export function WebSocketControl() {
    const [playerName, setPlayerName] = useState('');
    const [isConnected, setIsConnected] = useState(false);
    const [messages, setMessages] = useState<WebSocketMessage[]>([]);
    const [connectionStatus, setConnectionStatus] = useState('Disconnected');
    const socketRef = useRef<WebSocket | null>(null);

    const addMessage = (message: WebSocketMessage) => {
        setMessages(prev => [...prev, message]);
    };

    const createAdminMessage = (content: string): WebSocketMessage => {
        return {
            type: messageType.MessageTypeAdmin,
            timestamp: new Date(),
            playerName: '',
            content: content,

        }
    };

    const createErrorMessage = (content: string): WebSocketMessage => {
        return {
            type: messageType.MessageTypeError,
            timestamp: new Date(),
            playerName: '',
            content: content,

        }
    };

    const parseRawMessage = (event): WebSocketMessage => {
        let rawMessage = JSON.parse(event.data)
        return { ...rawMessage, timestamp: new Date(rawMessage.timestamp) }
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
            addMessage(createAdminMessage("Aleady connected!"));
            return;
        }

        try {
            setConnectionStatus('Connecting...');
            socketRef.current = new WebSocket(`ws://localhost:8080/liveGame/player/${playerName}`);

            socketRef.current.onopen = () => {
                setIsConnected(true);
                setConnectionStatus('Connected');
                addMessage(createAdminMessage("Connection established"));
            };

            socketRef.current.onmessage = (event) => {
                addMessage(parseRawMessage(event));
            };

            socketRef.current.onclose = () => {
                setIsConnected(false);
                setConnectionStatus('Disconnected');
                addMessage(createAdminMessage("Connection closed"));
            };

            socketRef.current.onerror = (error) => {
                setConnectionStatus('Error');
                addMessage(createErrorMessage(`${error.type}`));
            };
        } catch (error) {
            setConnectionStatus('Error');
            addMessage(createErrorMessage(`Connection error: ${error instanceof Error ? error.message : String(error)}`));
        }
    };

    const disconnectWebSocket = () => {
        if (socketRef.current) {
            socketRef.current.close();
            socketRef.current = null;
        } else {
            addMessage(createErrorMessage("No active connection to close"));
        }
    };

    // const sendTestMessage = () => {
    //     if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
    //         const message = `Hello server!`;
    //         socketRef.current.send(message);
    //         addMessage(protocolType.Sent, `${message}`);
    //     } else {
    //         addMessage(protocolType.Error, "Cannot send message: No connection");
    //     }
    // };

    return (
        <Flex className="websocket-control" align="center" justify="center" direction="column" gap="3">
            <h2>Player View</h2>

            <Flex direction="row" gap="3">
                <TextField.Root value={playerName} onChange={(event) => { setPlayerName(event.target.value) }} placeholder="Enter player name">
                    <TextField.Slot></TextField.Slot>
                </TextField.Root>
                <Button
                    onClick={connectWebSocket}
                    disabled={isConnected || playerName === ""}
                    className={isConnected ? "button-disabled" : "button-connect"}
                >
                    Join Game
                </Button>
            </Flex>
            <div className="status-indicator">
                WebSocket Status: <span className={`status-${connectionStatus.toLowerCase()}`}>{connectionStatus}</span>
            </div>
            <Flex direction="row" gap="3">
                <Button
                    onClick={disconnectWebSocket}
                    disabled={!isConnected}
                    className={!isConnected ? "button-disabled" : "button-disconnect"}
                >
                    Disconnect
                </Button>

                {/* <Button
                    onClick={sendTestMessage}
                    disabled={!isConnected}
                    className={!isConnected ? "button-disabled" : "button-send"}
                >
                    Send Test Message
                </Button> */}
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
                                    <Table.ColumnHeaderCell>PlayerName</Table.ColumnHeaderCell>
                                    <Table.ColumnHeaderCell>Content</Table.ColumnHeaderCell>
                                    <Table.ColumnHeaderCell>Time</Table.ColumnHeaderCell>
                                </Table.Row>
                            </Table.Header>
                            <Table.Body>
                                {messages.map((msg) => (
                                    <Table.Row>
                                        <Table.Cell>{msg.type.charAt(0).toUpperCase() + msg.type.slice(1)}</Table.Cell>
                                        <Table.Cell>{msg.playerName}</Table.Cell>
                                        <Table.Cell>{msg.content}</Table.Cell>
                                        <Table.Cell>{msg.timestamp.toISOString()}</Table.Cell>
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