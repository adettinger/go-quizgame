export enum messageType {
    Admin = "admin",
    Sent = "sent",
    Chat = "chat",
    Join = "join",
    Leave = "leave",
    GameUpdate = "game_update",
    Error = "error",
    PlayerList = "player_list"
}

export interface MessageTextContent {
    Text: string;
}

export interface MessagePlayerListContent {
    Names: string[];
}


export interface WebSocketMessage {
    type: messageType;
    timestamp: Date;
    playerName: string;
    content: MessageTextContent | MessagePlayerListContent; //Set to union of possible message content types that are defined in backend
}

export enum ConnectionStatus {
    Disconnected = "Disconnected",
    Connecting = "Connecting...",
    Connected = "Connected",
    Error = "Error",
}

export interface Player {
    name: string;
    color: string;
}