import { Flex, Button, TextField, Text, Card } from "@radix-ui/themes";
import { useEffect, useRef, useState } from "react";
import { ChatWindow, type chatMessage } from "../ChatWindow/ChatWindow";
import { MessageLog } from "./MessageLog";
import { getRandomColor, radixColors } from "./GameUtils";

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

export interface Player {
    name: string;
    color: string;
}



export function WebSocketControl() {
    const [playerName, setPlayerName] = useState('');
    const [serverError, setServerError] = useState('');
    const [isConnected, setIsConnected] = useState(false);
    const [playerList, setPlayerList] = useState<Player[]>([]);
    const [messages, setMessages] = useState<WebSocketMessage[]>([]);
    const [connectionStatus, setConnectionStatus] = useState('Disconnected');
    const socketRef = useRef<WebSocket | null>(null);
    const [chatMessages, setChatMessages] = useState<chatMessage[]>([])

    const availableColors = [...radixColors];

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

    // Create Players
    const createPlayerToAdd = (name: string): Player => {
        let color: string;
        if (availableColors.length > 0) {
            let colorIndex = Math.floor(Math.random() * availableColors.length);
            color = availableColors[colorIndex];
            // Remove the used color to avoid duplicates
            availableColors.splice(colorIndex, 1);
        } else {
            // Fallback if we run out of unique colors
            color = getRandomColor();
        }
        return {
            name,
            color
        }
    };

    const createPlayers = (names: string[]): Player[] => {
        // Create a copy of colors to track used colors (to avoid duplicates if needed)
        const players = names.map(name => {
            return createPlayerToAdd(name)
        });
        return players;
    };

    // Clean up the WebSocket connection when component unmounts
    useEffect(() => {
        return () => {
            if (socketRef.current) {
                socketRef.current.close();
            }
        };
    }, []);

    const handleMessage = (event) => {
        let msg = parseRawMessage(event);
        switch (msg.type) {
            case messageType.Chat:
                if ('Text' in msg.content) {
                    let msgToAdd = msg.content.Text;
                    console.log(playerList);
                    setChatMessages(prev => [...prev, {
                        playerName: msg.playerName,
                        color: playerList.find(player => player.name === msg.playerName)?.color || "gray",
                        message: msgToAdd,
                        timestamp: msg.timestamp,
                    }]);
                }
                break;
            case messageType.Join:
                setPlayerList(prev => {
                    if (prev.some(player => player.name === msg.playerName)) {
                        return prev;
                    }
                    return [...prev, createPlayerToAdd(msg.playerName)]
                });
                if ('Text' in msg.content) {
                    let msgToAdd = `${msg.playerName} ${msg.content.Text}`
                    setChatMessages(prev => [...prev, { playerName: "System", color: "red", message: msgToAdd, timestamp: msg.timestamp }]);
                }
                break;
            case messageType.Leave:
                setPlayerList(prev => prev.filter(player => player.name !== msg.playerName));
                if ('Text' in msg.content) {
                    let msgToAdd = `${msg.playerName} ${msg.content.Text}`
                    setChatMessages(prev => [...prev, { playerName: "System", color: "red", message: msgToAdd, timestamp: msg.timestamp }]);
                }
                break;
            case messageType.PlayerList:
                if ('Names' in msg.content) {
                    // create player with random color
                    setPlayerList(createPlayers(msg.content?.Names));
                }
                break;
            case messageType.Error:
                if ('Text' in msg.content) {
                    setServerError(msg.content.Text);
                }
                break;
        }
        addMessage(msg);
    };

    useEffect(() => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            socketRef.current.onmessage = handleMessage;
        }
    }, [playerList])

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

            socketRef.current.onmessage = handleMessage;

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

            <Flex gap="3" maxWidth={"50%"} wrap={"wrap"}>
                {isConnected && playerList.map((player) => (
                    <Card key={player.name} style={{ backgroundColor: player.color }}>{player.name}</Card>
                ))}
            </Flex>

            {isConnected &&
                <ChatWindow onMessageSend={sendChatMessage} messages={chatMessages} />
            }

            <MessageLog messages={messages} />
        </Flex>
    );
};