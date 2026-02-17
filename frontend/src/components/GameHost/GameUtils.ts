import type { messageType, WebSocketMessage } from "./GameTypes";

export const radixColors: string[] = [
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
    // 'amber', not visible enough
    'orange',
    // 'brown', not visible enough
    'gold',
    'bronze',
    // 'gray', not clear enough
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

export const parseRawMessage = (event: any): WebSocketMessage => {
    let parsedMsg = JSON.parse(event.data)
    return { ...parsedMsg, timestamp: new Date(parsedMsg.timestamp) }
};