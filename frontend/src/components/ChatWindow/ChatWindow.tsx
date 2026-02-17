import { Button, Card, Flex, ScrollArea, Text, TextField } from "@radix-ui/themes";
import { useState } from "react";
import './ChatWindow.scss';

export interface chatMessage {
    playerName: string;
    color: string;
    message: string;
    timestamp: Date;
}

export interface chatWindowProps {
    messages: chatMessage[];
    onMessageSend: (message: string) => void;
}

export function ChatWindow(props: chatWindowProps) {
    const [chatLine, setChatLine] = useState('');

    const handleSubmit = () => {
        props.onMessageSend(chatLine);
        setChatLine('');
    };

    return (
        <Card className="fullCard">
            <Flex direction="column" gap="3">
                <Card className="messagesCard">
                    <ScrollArea type="auto" scrollbars="both" style={{ height: 180 }}>
                        {props.messages.length == 0 ?
                            <Text color="red">No messages</Text>
                            :
                            <Flex direction="column" gap="1">
                                {
                                    props.messages.map((msg, index) => (
                                        <Flex direction="row" gap="1" align={"center"} key={index}>
                                            <Text className="playerName" color={msg.color || "gray" as any}>{msg.playerName}</Text>
                                            <Text className="message">{msg.message}</Text>
                                            <Text className="time" color="gray" size="1">{msg.timestamp.toLocaleTimeString()}</Text>
                                        </Flex>
                                    ))
                                }
                            </Flex>
                        }
                    </ScrollArea>
                </Card>
                <Flex direction="row" gap="1">
                    <TextField.Root
                        value={chatLine}
                        onChange={(event) => { setChatLine(event.target.value) }}
                        onKeyDown={(event) => {
                            if (event.key === 'Enter' && chatLine.trim() !== '') {
                                handleSubmit();
                            }
                        }}
                        placeholder="New message"
                        style={{ flex: 1 }}
                    >
                        <TextField.Slot />
                    </TextField.Root>
                    <Button
                        onClick={handleSubmit}
                        disabled={chatLine == ''}
                    >
                        Send
                    </Button>
                </Flex>
            </Flex>
        </Card>

    )
}