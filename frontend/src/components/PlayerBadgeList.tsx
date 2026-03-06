import { Badge, Text, Flex } from "@radix-ui/themes";
import { type Player } from "./GameHost/GameTypes";

export function PlayerBadgeList({ players }: { players: Player[] }) {
    if (players.length === 0) {
        return <></>
    }

    return (
        <Flex direction={"column"} maxWidth={"50%"} justify={"center"} align={"center"} gap="3">
            <Text>Players:</Text>
            <Flex direction={"row"} gap="3" wrap={"wrap"}>
                {players.map((player) => (
                    <Badge key={player.name} color={player.color || "gray" as any} size="3">{player.name}</Badge>
                ))}
            </Flex>
        </Flex>
    );
};