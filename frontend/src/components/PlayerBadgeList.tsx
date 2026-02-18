import { Badge, Flex } from "@radix-ui/themes";
import { type Player } from "./GameHost/GameTypes";

export function PlayerBadgeList({ players }: { players: Player[] }) {
    return (
        <Flex gap="3" maxWidth={"50%"} wrap={"wrap"}>
            {players.map((player) => (
                <Badge key={player.name} color={player.color || "gray" as any} size="3">{player.name}</Badge>
            ))}
        </Flex>
    );
};