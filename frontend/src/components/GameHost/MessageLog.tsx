import { Card, Flex, Table, Text } from "@radix-ui/themes";
import type { WebSocketMessage } from "./WebSocketControl";

export function MessageLog({ messages }: { messages: WebSocketMessage[] }) {
    return (
        <Card className="message-log" style={{ minWidth: '590px', maxWidth: '80%' }} >
            <Flex direction="column" align="center" gap="3">
                <Text align="center" weight="bold">Message Log</Text>
                <div className="messages">
                    {messages.length === 0 ? (
                        <p className="no-messages">No messages yet</p>
                    ) : (
                        <Table.Root variant="surface">
                            <Table.Header>
                                <Table.Row>
                                    <Table.ColumnHeaderCell>Type</Table.ColumnHeaderCell>
                                    <Table.ColumnHeaderCell>PlayerName</Table.ColumnHeaderCell>
                                    <Table.ColumnHeaderCell>Content</Table.ColumnHeaderCell>
                                    <Table.ColumnHeaderCell>Time</Table.ColumnHeaderCell>
                                </Table.Row>
                            </Table.Header>
                            <Table.Body>
                                {messages.map((msg) => (
                                    <Table.Row>
                                        <Table.Cell>{msg.type.charAt(0).toUpperCase() + msg.type.slice(1)}</Table.Cell>
                                        <Table.Cell>{msg.playerName}</Table.Cell>
                                        <Table.Cell>{JSON.stringify(msg.content)}</Table.Cell>
                                        <Table.Cell>{msg.timestamp.toISOString()}</Table.Cell>
                                    </Table.Row>
                                ))}
                            </Table.Body>
                        </Table.Root>
                    )}
                </div>
            </Flex>
        </Card>
    );
};