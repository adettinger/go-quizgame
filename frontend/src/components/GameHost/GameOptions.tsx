import { Button, Flex, Text, TextField, Tooltip } from "@radix-ui/themes";
import * as Form from "@radix-ui/react-form";
import { useEffect, useRef, useState } from "react";
import { ProblemPicker } from "./ProblemPicker";
import { QuestionMarkCircledIcon } from "@radix-ui/react-icons";
import { ConnectionStatus, GameStatus, messageType, QuestionStatus, type Player, type WebSocketMessage } from "./GameTypes";
import { ChatWindow, type chatMessage } from "../ChatWindow/ChatWindow";
import { createTextMessage, getRandomColor, parseRawMessage, playerColors } from "./GameUtils";
import { MessageLog } from "../MessageLog/MessageLog";
import { PlayerBadgeList } from "../PlayerBadgeList";
import { GameHostForm } from "./GameHostForm";

interface gameOptions {
    timeLimit: number,
    questionIds: string[],
}

export function GameOptions() {
    const [playerList, setPlayerList] = useState<Player[]>([]);
    const [messages, setMessages] = useState<WebSocketMessage[]>([]);
    const [chatMessages, setChatMessages] = useState<chatMessage[]>([]);
    const [gameStatus, setGameStatus] = useState<GameStatus>(GameStatus.NotStarted);
    const [questionStatus, setQuestionStatus] = useState<QuestionStatus>(QuestionStatus.NotStarted);



    const socketRef = useRef<WebSocket | null>(null);

    const [serverError, setServerError] = useState('');
    const [connectionStatus, setConnectionStatus] = useState<ConnectionStatus>(ConnectionStatus.Disconnected);

    const availableColors = [...playerColors];


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

    const handleMessage = (event: any) => {
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
        return () => {
            // Clean up to avoid memory leaks
            if (socketRef.current) {
                socketRef.current.onmessage = null;
            }
        };
    }, [playerList])

    const connectWebSocket = async (timeLimit: number, questionIds: string[]) => {
        if (isSocketConnected()) {
            addMessage(createTextMessage(messageType.Admin, "Aleady connected!"));
            return;
        }
        try {
            const wsUrl = `ws://localhost:8080/liveGame/host?timeLimit=${timeLimit}&questionIds=${questionIds.map(id => encodeURIComponent(id)).join(",")}`;
            const httpUrl = wsUrl.replace('ws:', 'http:');

            const response = await fetch(httpUrl, {
                method: 'GET',
            }).catch(error => {
                console.error("Fetch error: ", error);
                throw new Error("Network error when checking connection availability");
            });

            // If the response is not successful, handle different status codes
            if (!response.ok) {
                const statusCode = response.status;
                let errorMessage = `Server rejected connection (${statusCode})`;

                // Differentiate between different status codes
                if (statusCode === 400) {
                    setServerError('Invalid game options');
                    errorMessage = "Invalid game options (400 Bad Request)";
                } else if (statusCode === 403) {
                    setServerError('Access denied');
                    errorMessage = "Access denied (403 Forbidden)";
                } else if (statusCode >= 500) {
                    setServerError('Server error occurred');
                    errorMessage = "Server error occurred (500 Internal Server Error)";
                }

                throw new Error(errorMessage);
            }

            setConnectionStatus(ConnectionStatus.Connecting);
            socketRef.current = new WebSocket(wsUrl);

            socketRef.current.onopen = () => {
                setConnectionStatus(ConnectionStatus.Connected);
                addMessage(createTextMessage(messageType.Admin, "Connection established"));
                setServerError('');
            };

            socketRef.current.onmessage = handleMessage;

            socketRef.current.onclose = () => {
                setConnectionStatus(ConnectionStatus.Disconnected);
                setChatMessages([]);
                setPlayerList([]); //Not necessary because will be overwritten if rejoin, but render looks cleaner
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
    }

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
            addMessage(createTextMessage(messageType.Error, `Cannot send message: Connection not open. Connection status: ${connectionStatus}`));
        }
    };

    const startGame = () => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            socketRef.current.send(JSON.stringify(createTextMessage(messageType.StartGame, "")));
        } else {
            addMessage(createTextMessage(messageType.Error, `Cannot start game: Connection not open. Connection status: ${connectionStatus}`));
        }
    };

    return (
        <Flex direction="column" justify={"center"} align={"center"} gap="3">
            <div className="status-indicator">
                WebSocket Status: <span>{connectionStatus}</span>
            </div>
            <>
                {serverError !== '' &&
                    <Text color="red">Error: {serverError}</Text>
                }
            </>
            {
                isSocketConnected() ?
                    <>
                        <PlayerBadgeList players={playerList} />
                        <ChatWindow onMessageSend={sendChatMessage} messages={chatMessages} />
                    </>
                    :
                    <GameHostForm
                        onSubmit={(timeLimit: number, questionIds: string[]) => connectWebSocket(timeLimit, questionIds)}
                    />
            }
            <MessageLog messages={messages} />
        </Flex>
    );
};