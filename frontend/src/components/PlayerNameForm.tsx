import { Button, Flex, Text, TextField, Tooltip } from "@radix-ui/themes";
import { useState } from "react";


export interface PlayerNameFormProps {
    isConnected: boolean;
    onSubmit: (name: string) => void;
    onQuit: () => void;
}

export function PlayerNameForm({ isConnected, onSubmit, onQuit }: PlayerNameFormProps) {
    const [playerName, setPlayerName] = useState('');

    const MaxNameLength = 20;

    const isNameValid = (name: string) => {
        if (playerName.trim() === "" || playerName.trim() !== name || name.replace(/[^a-zA-Z0-9\s]/g, '') !== name || playerName.length > MaxNameLength) {
            return false;
        }
        return true;
    };

    return (
        <Flex direction="row" gap="3">
            <TextField.Root
                value={playerName}
                onChange={(event) => { setPlayerName(event.target.value.replace(/[^a-zA-Z0-9\s]/g, '')) }}
                placeholder="Enter player name"
                disabled={isConnected}
                onKeyDown={(event) => {
                    if (event.key === 'Enter' && playerName.trim() !== '') {
                        onSubmit(playerName);
                    }
                }}
            >
                <TextField.Slot />
            </TextField.Root>
            {isConnected ?
                <Button
                    onClick={onQuit}
                >
                    Quit Game
                </Button>
                :
                <Tooltip content={
                    isNameValid(playerName)
                        ? "Click to join the game"
                        :
                        <Flex direction="column">
                            <Text>Player name must contain only letters, numbers, and spaces, and cannot be empty</Text>
                            <Text>{`Cannot be longer than ${MaxNameLength} characters`}</Text>
                            <Text>Cannot begin or end with space</Text>
                        </Flex>
                }>
                    <Button
                        onClick={() => onSubmit(playerName)}
                        disabled={!isNameValid(playerName)}
                    >
                        Join Game
                    </Button>
                </Tooltip>

            }
        </Flex>

    );
};