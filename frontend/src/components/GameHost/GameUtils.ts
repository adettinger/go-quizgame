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

