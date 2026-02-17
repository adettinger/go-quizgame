import { Button, Flex, TextField } from "@radix-ui/themes";
import { useState } from "react";


export interface PlayerNameFormProps {
    isConnected: boolean;
    onSubmit: (name: string) => void;
    onQuit: () => void;
}

export function PlayerNameForm({ isConnected, onSubmit, onQuit }: PlayerNameFormProps) {
    const [playerName, setPlayerName] = useState('');

    return (
        <Flex direction="row" gap="3">
            <TextField.Root
                value={playerName}
                onChange={(event) => { setPlayerName(event.target.value) }}
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
                <Button
                    onClick={() => onSubmit(playerName)}
                    disabled={playerName.trim() === ""}
                >
                    Join Game
                </Button>
            }
        </Flex>

    );
};