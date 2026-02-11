import { Button, Card, Flex, ScrollArea, Text, TextField } from "@radix-ui/themes";
import { useState } from "react";

export interface chatMessage {
    playerName: string;
    message: string;
    timestamp: Date;
}

export interface chatWindowProps {
    messages: chatMessage[];
    onMessageSend: (message: string) => void;
}

export function ChatWindow(props: chatWindowProps) {
    const [chatLine, setChatLine] = useState('');
    return (
        <Card>
            <Flex direction="column" gap="1">
                <ScrollArea>
                    {props.messages.length == 0 ?
                        <Text color="red">No messages</Text>
                        :
                        <Flex direction="column" gap="1">
                            {
                                props.messages.map((msg) => (
                                    <Flex direction="row" gap="1">
                                        <Text color="red">{msg.playerName}</Text>
                                        <Text>{msg.message}</Text>
                                        <Text color="gray">{msg.timestamp.toLocaleTimeString()}</Text>
                                    </Flex>
                                ))
                            }
                        </Flex>
                    }
                </ScrollArea>
                <Flex direction="row" gap="1">

                    <TextField.Root
                        value={chatLine}
                        onChange={(event) => { setChatLine(event.target.value) }}
                        placeholder="New message"
                    >
                        <TextField.Slot />
                    </TextField.Root>
                    <Button onClick={() => { props.onMessageSend(chatLine); setChatLine(''); }}>Send</Button>
                </Flex>
            </Flex>
        </Card>

    )
}