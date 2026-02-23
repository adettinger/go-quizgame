import { Link } from 'react-router-dom';
import { useProblems } from '../hooks/useProblems';
import { Button, DropdownMenu, Flex, IconButton, Table } from "@radix-ui/themes"
import { TrashIcon } from '@radix-ui/react-icons';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { deleteProblemById } from '../services/problemService';
import { useToast } from './Toast/ToastContext';
import { ProblemType } from '../types/problem';

export function ProblemList() {
    const queryClient = useQueryClient();
    const { showToast } = useToast();
    const { data, isLoading, isError, error } = useProblems();

    const deleteMutation = useMutation({
        mutationFn: deleteProblemById,
        onSuccess: (id: string) => {
            queryClient.invalidateQueries({ queryKey: ['problems'] });
            queryClient.invalidateQueries({ queryKey: ['problem', id] });
            showToast('success', "Success", `Deleted problem ${id}`);
            console.log(`Deleted problem ${id}`);
        },
        onError: (error) => {
            console.log('Failed to delete problem', error)
        },
    })

    const handleDelete = (id: string) => {
        if (confirm('Are you sure you want to delete this problem?')) {
            deleteMutation.mutate(id);
        }
    }

    if (isLoading) {
        return <div className="text-center p-4">Loading problems...</div>;
    }

    if (isError) {
        return (
            <div className="text-center p-4 text-red-600">
                Error loading problems: {error.message}
            </div>
        );
    }

    if (!data || data.length === 0) {
        return <div className="text-center p-4">No problems found.</div>;
    }

    return (
        <Flex align="center" justify="center">
            <Table.Root variant="surface">
                <Table.Header>
                    <Table.Row>
                        <Table.ColumnHeaderCell>ID</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Type</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Question</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Choices</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Answer</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Actions</Table.ColumnHeaderCell>
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    {data.map((problem, index) => (
                        <Table.Row key={index}>
                            <Table.RowHeaderCell>
                                <Link to={`/problem/${problem.Id}`}>{problem.Id}</Link>
                            </Table.RowHeaderCell>
                            <Table.Cell>{problem.Type}</Table.Cell>
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
                                <IconButton
                                    color="red"
                                    variant="soft"
                                    onClick={() => { handleDelete(problem.Id) }}
                                    disabled={deleteMutation.isPending}
                                >
                                    <TrashIcon />
                                </IconButton>
                            </Table.Cell>
                        </Table.Row>
                    ))}
                </Table.Body>
            </Table.Root>
        </Flex>
    );
}