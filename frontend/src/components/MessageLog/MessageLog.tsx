import { Button, Card, Flex, ScrollArea, Table, Text } from "@radix-ui/themes";
import './MessageLog.scss'
import type { WebSocketMessage } from "../GameHost/GameTypes";
import { Dropdown } from "../Dropdown";

export function MessageLog({ messages }: { messages: WebSocketMessage[] }) {

    return (
        <Card className="messageLog">
            <Flex direction="column" align="center" gap="3">

                {/* <Flex align="center" justify="between" width="100%">
                    <div style={{ width: '24px' }} /> {/* Spacer to balance the layout
                    <Text align="center" weight="bold">Message Log</Text>
                    <Button onClick={() => setIsVisible(!isVisible)} variant="outline" color="gray">
                        {isVisible ?
                            <CaretUpIcon fontWeight={"bold"} />
                            :
                            <CaretDownIcon fontWeight={"bold"} />
                        }
                    </Button>
                </Flex> */}
                <Dropdown title="Message Log">

                    <ScrollArea className="messages" scrollbars="vertical">
                        {messages.length === 0 ? (
                            <Flex align="center" justify="center">
                                <p className="no-messages">No messages yet</p>
                            </Flex>
                        ) : (
                            <Table.Root variant="surface">
                                <Table.Header>
                                    <Table.Row>
                                        <Table.ColumnHeaderCell>Type</Table.ColumnHeaderCell>
                                        <Table.ColumnHeaderCell >Player Name</Table.ColumnHeaderCell>
                                        <Table.ColumnHeaderCell>Content</Table.ColumnHeaderCell>
                                        <Table.ColumnHeaderCell>Time</Table.ColumnHeaderCell>
                                    </Table.Row>
                                </Table.Header>
                                <Table.Body>
                                    {messages.map((msg, index) => (
                                        <Table.Row key={index} >
                                            <Table.Cell>{msg.type.charAt(0).toUpperCase() + msg.type.slice(1)}</Table.Cell>
                                            <Table.Cell>{msg.playerName}</Table.Cell>
                                            <Table.Cell className="content">{JSON.stringify(msg.content)}</Table.Cell>
                                            <Table.Cell>{msg.timestamp.toISOString()}</Table.Cell>
                                        </Table.Row>
                                    ))}
                                </Table.Body>
                            </Table.Root>
                        )}
                    </ScrollArea>
                </Dropdown>
            </Flex>
        </Card>
    );
};