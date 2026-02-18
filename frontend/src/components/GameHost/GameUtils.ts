import type { messageType, WebSocketMessage } from "./GameTypes";

export const playerColors: string[] = [
    'tomato',
    // 'red', For system
    'crimson',
    'pink',
    'plum',
    'purple',
    'violet',
    'indigo',
    'blue',
    'cyan',
    'teal',
    'green',
    'grass',
    'lime',
    'yellow',
    // 'amber', too close to yellow
    'orange',
    'brown',
    'gold',
    // 'bronze', to close to brown
    // 'gray', not clear enough
    // 'mauve', not set for badge
    // 'slate', not set for badge
    // 'sage', not set for badge
    // 'olive', not set for badge
    // 'sand', not set for badge
];

export const getRandomColor = (): string => {
    const randomIndex = Math.floor(Math.random() * playerColors.length);
    return playerColors[randomIndex];
};

export const createTextMessage = (type: messageType, content: string): WebSocketMessage => {
    return {
        type: type,
        timestamp: new Date(),
        playerName: '',
        content: { Text: content },
    }
}

export const parseRawMessage = (event: any): WebSocketMessage => {
    let parsedMsg = JSON.parse(event.data)
    return { ...parsedMsg, timestamp: new Date(parsedMsg.timestamp) }
};