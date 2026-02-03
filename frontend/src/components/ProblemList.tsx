import { Link } from 'react-router-dom';
import { useProblems } from '../hooks/useProblems';
import { Table } from "@radix-ui/themes"

export function ProblemList() {
    const { data, isLoading, isError, error } = useProblems();

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
        <Table.Root variant="surface">
            <Table.Header>
                <Table.Row>
                    <Table.ColumnHeaderCell>ID</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Question</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Answer</Table.ColumnHeaderCell>
                </Table.Row>
            </Table.Header>

            <Table.Body>
                {data.map((problem) => (
                    <Table.Row>
                        <Table.RowHeaderCell>
                            <Link to={`/problem/${problem.Id}`}>{problem.Id}</Link>
                        </Table.RowHeaderCell>
                        <Table.Cell>{problem.Question}</Table.Cell>
                        <Table.Cell>{problem.Answer}</Table.Cell>
                    </Table.Row>
                ))}
            </Table.Body>
        </Table.Root>
    );
}