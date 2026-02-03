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
                        <Table.RowHeaderCell>{problem.Id}</Table.RowHeaderCell>
                        <Table.Cell>{problem.Question}</Table.Cell>
                        <Table.Cell>{problem.Answer}</Table.Cell>
                    </Table.Row>
                ))}
            </Table.Body>
        </Table.Root>
    );

    // return (
    //     <div className="container mx-auto p-4">
    //         <h2 className="text-2xl font-bold mb-4">Problems</h2>
    //         <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
    //             {data.map((problem) => (
    //                 <div
    //                     key={problem.Id}
    //                     className="border rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow"
    //                 >
    //                     <h3 className="font-semibold text-lg mb-2">Question:</h3>
    //                     <p className="mb-4">{problem.Question}</p>
    //                     <h3 className="font-semibold text-lg mb-2">Answer:</h3>
    //                     <p>{problem.Answer}</p>
    //                 </div>
    //             ))}
    //         </div>
    //     </div>
    // );
}