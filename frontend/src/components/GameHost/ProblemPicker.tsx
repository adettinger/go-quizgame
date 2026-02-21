import { Flex, IconButton, Table, Text } from "@radix-ui/themes";
import type { Problem } from "../../types/problem";
import { Cross1Icon, PlusIcon } from "@radix-ui/react-icons";

export interface ProblemPickerProps {
    selectedQuestions: Problem[];
    availableQuestions: Problem[];
    selectQuestion: (questionToAdd: Problem) => void;
    deselectQuestion: (questionToRemove: Problem) => void;
    isLoading: boolean;
    isError: boolean;
    error: Error | null;
};

export function ProblemPicker({ selectedQuestions, availableQuestions, selectQuestion, deselectQuestion, isLoading, isError, error }: ProblemPickerProps) {
    if (isLoading) {
        return <div className="text-center p-4">Loading problems...</div>;
    }

    if (isError && !!error) {
        return (
            <div className="text-center p-4 text-red-600">
                Error loading problems: {error.message}
            </div>
        );
    }

    if (availableQuestions.length === 0 && selectedQuestions.length === 0) {
        return <div className="text-center p-4">No problems found.</div>;
    }

    return (
        <Flex direction="column" align="center" justify="center">
            <h3>Selected Questions</h3>
            {selectedQuestions.length === 0 ?
                <Text>No Questions Selected</Text>
                :
                <Table.Root variant="surface">
                    <Table.Header>
                        <Table.Row>
                            <Table.ColumnHeaderCell>ID</Table.ColumnHeaderCell>
                            <Table.ColumnHeaderCell>Question</Table.ColumnHeaderCell>
                            <Table.ColumnHeaderCell>Answer</Table.ColumnHeaderCell>
                            <Table.ColumnHeaderCell>Actions</Table.ColumnHeaderCell>
                        </Table.Row>
                    </Table.Header>

                    <Table.Body>
                        {selectedQuestions.map((question) => (
                            <Table.Row key={question.Id}>
                                <Table.Cell>{question.Id}</Table.Cell>
                                <Table.Cell>{question.Question}</Table.Cell>
                                <Table.Cell>{question.Answer}</Table.Cell>
                                <Table.Cell>
                                    <IconButton
                                        color="red"
                                        variant="soft"
                                        onClick={() => { deselectQuestion(question) }}
                                    >
                                        <Cross1Icon />
                                    </IconButton>
                                </Table.Cell>
                            </Table.Row>
                        ))}
                    </Table.Body>
                </Table.Root>
            }

            <h3>Available Questions</h3>
            {availableQuestions.length === 0 ?
                <Text>No available questions remaining</Text>
                :
                <Table.Root variant="surface">
                    <Table.Header>
                        <Table.Row>
                            <Table.ColumnHeaderCell>ID</Table.ColumnHeaderCell>
                            <Table.ColumnHeaderCell>Question</Table.ColumnHeaderCell>
                            <Table.ColumnHeaderCell>Answer</Table.ColumnHeaderCell>
                            <Table.ColumnHeaderCell>Actions</Table.ColumnHeaderCell>
                        </Table.Row>
                    </Table.Header>

                    <Table.Body>
                        {availableQuestions.map((question) => (
                            <Table.Row key={question.Id}>
                                <Table.Cell>{question.Id}</Table.Cell>
                                <Table.Cell>{question.Question}</Table.Cell>
                                <Table.Cell>{question.Answer}</Table.Cell>
                                <Table.Cell>
                                    <IconButton
                                        color="blue"
                                        variant="soft"
                                        onClick={() => { selectQuestion(question) }}
                                    >
                                        <PlusIcon />
                                    </IconButton>
                                </Table.Cell>
                            </Table.Row>
                        ))}
                    </Table.Body>
                </Table.Root>
            }
        </Flex>
    );
};