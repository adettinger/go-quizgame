import type { messageType, WebSocketMessage } from "./GameTypes";

export const radixColors: string[] = [
    'tomato',
    // 'red',
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
    // 'amber',
    'orange',
    'brown',
    'gold',
    'bronze',
    // 'gray',
    'mauve',
    'slate',
    'sage',
    'olive',
    'sand',
];

export const getRandomColor = (): string => {
    const randomIndex = Math.floor(Math.random() * radixColors.length);
    return radixColors[randomIndex];
};

export const createTextMessage = (type: messageType, content: string): WebSocketMessage => {
    return {
        type: type,
        timestamp: new Date(),
        playerName: '',
        content: { Text: content },
    }
}

export const parseRawMessage = (event): WebSocketMessage => {
    let parsedMsg = JSON.parse(event.data)
    return { ...parsedMsg, timestamp: new Date(parsedMsg.timestamp) }
};