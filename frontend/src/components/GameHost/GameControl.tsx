import { Flex, Text, Card } from "@radix-ui/themes";
import { useEffect, useRef, useState } from "react";
import { ChatWindow, type chatMessage } from "../ChatWindow/ChatWindow";
import { MessageLog } from "../MessageLog/MessageLog";
import { createTextMessage, getRandomColor, parseRawMessage, radixColors } from "./GameUtils";
import { ConnectionStatus, messageType, type Player, type WebSocketMessage } from "./GameTypes";
import { PlayerNameForm } from "../PlayerNameForm";

export function GameControl() {
    const [serverError, setServerError] = useState('');
    const [playerList, setPlayerList] = useState<Player[]>([]);
    const [messages, setMessages] = useState<WebSocketMessage[]>([]);
    const [connectionStatus, setConnectionStatus] = useState<ConnectionStatus>(ConnectionStatus.Disconnected); //For display purposes only!
    const socketRef = useRef<WebSocket | null>(null);
    const [chatMessages, setChatMessages] = useState<chatMessage[]>([])

    const availableColors = [...radixColors];

    const addMessage = (message: WebSocketMessage) => {
        setMessages(prev => [...prev, message]);
    };

    const createPlayerToAdd = (name: string): Player => {
        let color: string;
        if (availableColors.length > 0) {
            let colorIndex = Math.floor(Math.random() * availableColors.length);
            color = availableColors[colorIndex];
            availableColors.splice(colorIndex, 1);
        } else {
            color = getRandomColor();
        }
        return {
            name,
            color
        }
    };

    const createPlayers = (names: string[]): Player[] => {
        const players = names.map(name => {
            return createPlayerToAdd(name)
        });
        return players;
    };

    const isSocketConnected = (): boolean => {
        return !!socketRef.current && socketRef.current.readyState === WebSocket.OPEN
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
        if (msg.type != messageType.Error && serverError != '') {
            setServerError('');
        }
        switch (msg.type) {
            case messageType.Chat:
                if ('Text' in msg.content) {
                    let msgToAdd = msg.content.Text;
                    setChatMessages(prev => [...prev, {
                        playerName: msg.playerName,
                        color: msg.playerName === "System" ? "red" : playerList.find(player => player.name === msg.playerName)?.color || "gray",
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

    // socketRef.current.onmessage must be reset to use updated playerList state
    useEffect(() => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            socketRef.current.onmessage = handleMessage;
        }
    }, [playerList])

    const connectWebSocket = (name: string) => {
        if (isSocketConnected()) {
            addMessage(createTextMessage(messageType.Admin, "Aleady connected!"));
            return;
        }

        try {
            setConnectionStatus(ConnectionStatus.Connecting);
            socketRef.current = new WebSocket(`ws://localhost:8080/liveGame/player/${name.trim()}`);

            socketRef.current.onopen = () => {
                setConnectionStatus(ConnectionStatus.Connected);
                addMessage(createTextMessage(messageType.Admin, "Connection established"));
            };

            socketRef.current.onmessage = handleMessage;

            socketRef.current.onclose = () => {
                setConnectionStatus(ConnectionStatus.Disconnected);
                setChatMessages([]);
                addMessage(createTextMessage(messageType.Admin, "Connection closed"));
            };

            socketRef.current.onerror = (error) => {
                setConnectionStatus(ConnectionStatus.Error);
                addMessage(createTextMessage(messageType.Error, `${error.type}`));
            };
        } catch (error) {
            setConnectionStatus(ConnectionStatus.Error);
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
        <Flex align="center" justify="center" direction="column" gap="3">
            <h2>Player View</h2>

            <PlayerNameForm
                isConnected={isSocketConnected()}
                onSubmit={(name: string) => { connectWebSocket(name) }}
                onQuit={disconnectWebSocket}
            />

            <div className="status-indicator">
                WebSocket Status: <span>{connectionStatus}</span>
            </div>

            {serverError !== '' &&
                <Text color="red">Error: {serverError}</Text>
            }

            {isSocketConnected() &&
                <>
                    <Flex gap="3" maxWidth={"50%"} wrap={"wrap"}>
                        {playerList.map((player) => (
                            <Card key={player.name} style={{ backgroundColor: player.color }}>{player.name}</Card>
                        ))}
                    </Flex>

                    <ChatWindow onMessageSend={sendChatMessage} messages={chatMessages} />
                </>
            }

            <MessageLog messages={messages} />
        </Flex>
    );
};