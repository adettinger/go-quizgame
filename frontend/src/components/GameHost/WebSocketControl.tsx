import { Flex, Button, TextField, Text } from "@radix-ui/themes";
import { useEffect, useRef, useState } from "react";
import { ChatWindow, type chatMessage } from "../ChatWindow/ChatWindow";
import { MessageLog } from "./MessageLog";

enum messageType {
    Admin = "admin",
    Sent = "sent",
    Chat = "chat",
    Join = "join",
    Leave = "leave",
    GameUpdate = "game_update",
    Error = "error",
    PlayerList = "player_list"
}

interface MessageTextContent {
    Text: string;
}

interface MessagePlayerListContent {
    Names: string[];
}


export interface WebSocketMessage {
    type: messageType;
    timestamp: Date;
    playerName: string;
    content: MessageTextContent | MessagePlayerListContent; //Set to union of possible message content types that are defined in backend
}


export function WebSocketControl() {
    const [playerName, setPlayerName] = useState('');
    const [serverError, setServerError] = useState('');
    const [isConnected, setIsConnected] = useState(false);
    const [playerList, setPlayerList] = useState<string[]>([]);
    const [messages, setMessages] = useState<WebSocketMessage[]>([]);
    const [connectionStatus, setConnectionStatus] = useState('Disconnected');
    const socketRef = useRef<WebSocket | null>(null);
    const [chatMessages, setChatMessages] = useState<chatMessage[]>([])

    const addMessage = (message: WebSocketMessage) => {
        if (message.type != messageType.Error && serverError != '') { //redundant
            setServerError('');
        }
        setMessages(prev => [...prev, message]);
    };

    const createTextMessage = (type: messageType, content: string): WebSocketMessage => {
        return {
            type: type,
            timestamp: new Date(),
            playerName: '',
            content: { Text: content },
        }
    }

    const parseRawMessage = (event): WebSocketMessage => {
        let parsedMsg = JSON.parse(event.data)
        return { ...parsedMsg, timestamp: new Date(parsedMsg.timestamp) }
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
            addMessage(createTextMessage(messageType.Admin, "Aleady connected!"));
            return;
        }

        try {
            setConnectionStatus('Connecting...');
            socketRef.current = new WebSocket(`ws://localhost:8080/liveGame/player/${playerName.trim()}`);

            socketRef.current.onopen = () => {
                setIsConnected(true);
                setConnectionStatus('Connected');
                addMessage(createTextMessage(messageType.Admin, "Connection established"));
            };

            socketRef.current.onmessage = (event) => {
                let msg = parseRawMessage(event);
                switch (msg.type) {
                    case messageType.Chat:
                        if ('Text' in msg.content) {
                            let msgToAdd = msg.content.Text;
                            setChatMessages(prev => [...prev, { playerName: msg.playerName, message: msgToAdd, timestamp: msg.timestamp }]);
                        }
                        break;
                    case messageType.Join:
                    case messageType.Leave:
                        if ('Text' in msg.content) {
                            let msgToAdd = `${msg.playerName} ${msg.content.Text}`
                            setChatMessages(prev => [...prev, { playerName: "System", message: msgToAdd, timestamp: msg.timestamp }]);
                        }
                        break;
                    case messageType.PlayerList:
                        console.log('received player list', msg.content);
                        break;
                    case messageType.Error:
                        if ('Text' in msg.content) {
                            setServerError(msg.content.Text);
                        }
                        break;
                }
                addMessage(msg);
            };

            socketRef.current.onclose = () => {
                setIsConnected(false);
                setConnectionStatus('Disconnected');
                addMessage(createTextMessage(messageType.Admin, "Connection closed"));
            };

            socketRef.current.onerror = (error) => {
                setConnectionStatus('Error');
                addMessage(createTextMessage(messageType.Error, `${error.type}`));
            };
        } catch (error) {
            setConnectionStatus('Error');
            addMessage(createTextMessage(messageType.Error, `Connection error: ${error instanceof Error ? error.message : String(error)}`));
        }
    };

    const disconnectWebSocket = () => {
        if (socketRef.current) {
            socketRef.current.close();
            socketRef.current = null;
        } else {
            addMessage(createTextMessage(messageType.Error, "No active connection to close"));
        }
    };

    const sendChatMessage = (message: string) => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            socketRef.current.send(JSON.stringify(createTextMessage(messageType.Chat, message)));
        } else {
            addMessage(createTextMessage(messageType.Error, "Cannot send message: No connection"));
        }
    };

    return (
        <Flex className="websocket-control" align="center" justify="center" direction="column" gap="3">
            <h2>Player View</h2>

            <Flex direction="row" gap="3">
                <TextField.Root
                    value={playerName}
                    onChange={(event) => { setPlayerName(event.target.value) }}
                    placeholder="Enter player name"
                    disabled={isConnected}
                    onKeyDown={(event) => {
                        if (event.key === 'Enter' && playerName.trim() !== '') {
                            connectWebSocket();
                        }
                    }}
                >
                    <TextField.Slot />
                </TextField.Root>
                <Button
                    onClick={connectWebSocket}
                    disabled={isConnected || playerName.trim() === ""}
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