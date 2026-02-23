import { Cross1Icon, PlusIcon } from "@radix-ui/react-icons";
import { Flex, IconButton, Table, TextField } from "@radix-ui/themes";
import { MAX_PROBLEM_CHOICES } from "../types/problem";
import { useState } from "react";


export function EditChoicesTable({ choices, setChoices }: { choices: string[], setChoices: (choices: string[]) => void }) {
    const [choiceToAdd, setChoiceToAdd] = useState('');

    const addChoice = (input: string) => {
        setChoices([...choices, input]);
        setChoiceToAdd('');
    };

    const removeChoice = (input: string) => {
        setChoices(choices.filter((choice) => choice !== input));
    };

    return (
        <Flex justify="center">
            <Table.Root variant="surface" size="1" style={{ width: '50%' }} >
                <Table.Header>
                    <Table.Row>
                        <Table.ColumnHeaderCell justify="center">Choices</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell width="1%" justify="center">Actions</Table.ColumnHeaderCell>
                    </Table.Row>
                </Table.Header>
                <Table.Body>
                    {/* Show choices */}
                    {choices.map((choice, index) => (
                        <Table.Row key={index}>
                            <Table.Cell justify="center">{choice}</Table.Cell>
                            <Table.Cell width="1%" justify="center">
                                <IconButton
                                    color="red"
                                    variant="soft"
                                    onClick={() => { removeChoice(choice) }}
                                >
                                    <Cross1Icon />
                                </IconButton>
                            </Table.Cell>
                        </Table.Row>
                    ))}
                    {/* Add choice option */}
                    {choices.length < MAX_PROBLEM_CHOICES &&
                        <Table.Row>
                            <Table.Cell justify="center">
                                <TextField.Root
                                    placeholder="Enter choice text"
                                    value={choiceToAdd}
                                    onChange={(event) => setChoiceToAdd(event.target.value)}
                                >

                                </TextField.Root>
                            </Table.Cell>
                            <Table.Cell width="1%" justify="center">
                                <IconButton
                                    color="blue"
                                    variant="soft"
                                    onClick={() => { addChoice(choiceToAdd) }}
                                    disabled={choiceToAdd.trim() === ""}
                                >
                                    <PlusIcon />
                                </IconButton>
                            </Table.Cell>
                        </Table.Row>
                    }
                </Table.Body>
            </Table.Root>
        </Flex>
    );
}