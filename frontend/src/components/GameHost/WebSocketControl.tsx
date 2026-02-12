import { Flex, Button, Table, TextField, Text } from "@radix-ui/themes";
import { useEffect, useRef, useState } from "react";
import { ChatWindow, type chatMessage } from "./ChatWindow";
import { MessageLog } from "./MessageLog";

enum messageType {
    Admin = "admin",
    Sent = "sent",
    Chat = "chat",
    Join = "join",
    Leave = "leave",
    GameUpdate = "game_update",
    Error = "error",
}


export interface WebSocketMessage {
    type: messageType;
    timestamp: Date;
    playerName: string;
    content: string;
}


export function WebSocketControl() {
    const [playerName, setPlayerName] = useState('');
    const [serverError, setServerError] = useState('');
    const [isConnected, setIsConnected] = useState(false);
    const [messages, setMessages] = useState<WebSocketMessage[]>([]);
    const [connectionStatus, setConnectionStatus] = useState('Disconnected');
    const socketRef = useRef<WebSocket | null>(null);
    const [chatMessages, setChatMessages] = useState<chatMessage[]>([])

    const addMessage = (message: WebSocketMessage) => {
        if (message.type != messageType.Error && serverError != '') {
            setServerError('');
        }
        setMessages(prev => [...prev, message]);
    };

    const createMessage = (type: messageType, content: string): WebSocketMessage => {
        return {
            type: type,
            timestamp: new Date(),
            playerName: '',
            content: content,
        }
    }

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
            addMessage(createMessage(messageType.Admin, "Aleady connected!"));
            return;
        }

        try {
            setConnectionStatus('Connecting...');
            socketRef.current = new WebSocket(`ws://localhost:8080/liveGame/player/${playerName}`);

            socketRef.current.onopen = () => {
                setIsConnected(true);
                setConnectionStatus('Connected');
                addMessage(createMessage(messageType.Admin, "Connection established"));
            };

            socketRef.current.onmessage = (event) => {
                let msg = parseRawMessage(event);
                switch (msg.type) {
                    case messageType.Chat:
                        setChatMessages(prev => [...prev, { playerName: msg.playerName, message: msg.content, timestamp: msg.timestamp }]);
                        break;
                    case messageType.Join:
                    case messageType.Leave:
                        setChatMessages(prev => [...prev, { playerName: "System", message: msg.playerName + " " + msg.content, timestamp: msg.timestamp }]);
                        break;
                    case messageType.Error:
                        setServerError(msg.content);
                        break;
                }
                addMessage(msg);
            };

            socketRef.current.onclose = () => {
                setIsConnected(false);
                setConnectionStatus('Disconnected');
                addMessage(createMessage(messageType.Admin, "Connection closed"));
            };

            socketRef.current.onerror = (error) => {
                setConnectionStatus('Error');
                addMessage(createMessage(messageType.Error, `${error.type}`));
            };
        } catch (error) {
            setConnectionStatus('Error');
            addMessage(createMessage(messageType.Error, `Connection error: ${error instanceof Error ? error.message : String(error)}`));
        }
    };

    const disconnectWebSocket = () => {
        if (socketRef.current) {
            socketRef.current.close();
            socketRef.current = null;
        } else {
            addMessage(createMessage(messageType.Error, "No active connection to close"));
        }
    };

    const sendChatMessage = (message: string) => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            socketRef.current.send(JSON.stringify(createMessage(messageType.Chat, message)));
        } else {
            addMessage(createMessage(messageType.Error, "Cannot send message: No connection"));
        }
    };

    return (
        <Flex className="websocket-control" align="center" justify="center" direction="column" gap="3">
            <h2>Player View</h2>

            <Flex direction="row" gap="3">
                <TextField.Root value={playerName} onChange={(event) => { setPlayerName(event.target.value) }} placeholder="Enter player name">
                    <TextField.Slot />
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
            {serverError !== '' &&
                <Text color="red">Error: {serverError}</Text>
            }
            <Flex direction="row" gap="3">
                <Button
                    onClick={disconnectWebSocket}
                    disabled={!isConnected}
                    className={!isConnected ? "button-disabled" : "button-disconnect"}
                >
                    Disconnect
                </Button>
            </Flex>

            {isConnected &&
                <ChatWindow onMessageSend={sendChatMessage} messages={chatMessages} />
            }

            <MessageLog messages={messages} />
        </Flex>
    );
};