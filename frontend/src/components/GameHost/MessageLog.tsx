import { Table } from "@radix-ui/themes";
import type { WebSocketMessage } from "./WebSocketControl";

export function MessageLog({ messages }: { messages: WebSocketMessage[] }) {
    return (
        <div className="message-log">
            <h3>Message Log</h3>
            <div className="messages">
                {messages.length === 0 ? (
                    <p className="no-messages">No messages yet</p>
                ) : (
                    <Table.Root>
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
                                    <Table.Cell>{msg.content}</Table.Cell>
                                    <Table.Cell>{msg.timestamp.toISOString()}</Table.Cell>
                                </Table.Row>
                            ))}
                        </Table.Body>
                    </Table.Root>
                )}
            </div>
        </div>
    );
};