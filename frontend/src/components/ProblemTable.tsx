import { Link } from 'react-router-dom';
import { Button, DropdownMenu, Flex, IconButton, Table } from "@radix-ui/themes"
import { Cross1Icon, Pencil1Icon, PlusIcon, TrashIcon } from '@radix-ui/react-icons';
import { ProblemType, type Problem } from '../types/problem';
import React from 'react';

export interface ProblemTableProps {
    Problems: Problem[] | undefined,
    IsLoading?: boolean,
    IsError?: boolean,
    Error?: Error | null,
    ShowIds: boolean,
    OnEdit?: (id: string) => void,
    OnDelete?: (id: string) => void,
    OnAdd?: (problem: Problem) => void,
    OnRemove?: (problem: Problem) => void,
    DisableActions?: boolean,
};

export function ProblemTable(props: ProblemTableProps) {
    if (props.IsLoading) {
        return <div className="text-center p-4">Loading problems...</div>;
    }

    if (props.IsError) {
        return (
            <div className="text-center p-4 text-red-600">
                Error loading problems: {props.Error?.message}
            </div>
        );
    }

    if (!props.Problems || props.Problems.length === 0) {
        return <div className="text-center p-4">No problems found.</div>;
    }

    return (
        <Table.Root variant="surface">
            <Table.Header>
                <Table.Row>
                    {props.ShowIds &&
                        <Table.ColumnHeaderCell>ID</Table.ColumnHeaderCell>
                    }
                    <Table.ColumnHeaderCell>Type</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Question</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Choices</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Answer</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Actions</Table.ColumnHeaderCell>
                </Table.Row>
            </Table.Header>

            <Table.Body>
                {props.Problems.map((problem) => (
                    <Table.Row key={problem.Id}>
                        {props.ShowIds &&
                            <Table.RowHeaderCell>
                                <Link to={`/problem/${problem.Id}`}>{problem.Id}</Link>
                            </Table.RowHeaderCell>
                        }
                        <Table.Cell>{problem.Type}</Table.Cell>
                        {/* TODO: Rendering longer questions? */}
                        <Table.Cell>{problem.Question}</Table.Cell>
                        <Table.Cell>
                            {problem.Type === ProblemType.Choice &&
                                <DropdownMenu.Root>
                                    <DropdownMenu.Trigger>
                                        <Button color='gray' variant='soft'>Choices <DropdownMenu.TriggerIcon /></Button>
                                    </DropdownMenu.Trigger>
                                    <DropdownMenu.Content color="gray" variant='soft'>
                                        {problem.Choices.map((choice, index) => (
                                            <DropdownMenu.Item key={index}>{choice}</DropdownMenu.Item>
                                        ))}
                                    </DropdownMenu.Content>
                                </DropdownMenu.Root>
                            }
                        </Table.Cell>
                        <Table.Cell>{problem.Answer}</Table.Cell>
                        <Table.Cell>
                            <Flex gap="2">
                                {props.OnAdd &&
                                    <IconButton
                                        color="blue"
                                        variant="soft"
                                        onClick={() => { props.OnAdd?.(problem) }}
                                        disabled={!!props.DisableActions}
                                    >
                                        <PlusIcon />
                                    </IconButton>
                                }
                                {props.OnRemove &&
                                    <IconButton
                                        color="red"
                                        variant="soft"
                                        onClick={() => { props.OnRemove?.(problem) }}
                                        disabled={!!props.DisableActions}
                                    >
                                        <Cross1Icon />
                                    </IconButton>
                                }
                                {props.OnEdit &&
                                    <IconButton
                                        color="indigo"
                                        variant="soft"
                                        onClick={() => { props.OnEdit?.(problem.Id) }}
                                        disabled={!!props.DisableActions}
                                    >
                                        <Pencil1Icon />
                                    </IconButton>
                                }
                                {props.OnDelete &&
                                    <IconButton
                                        color="red"
                                        variant="soft"
                                        onClick={() => { props.OnDelete?.(problem.Id) }}
                                        disabled={!!props.DisableActions}
                                    >
                                        <TrashIcon />
                                    </IconButton>
                                }
                            </Flex>
                        </Table.Cell>
                    </Table.Row>
                ))}
            </Table.Body>
        </Table.Root>
    );
};