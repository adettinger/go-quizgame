import { DataList } from "@radix-ui/themes";
import { useParams } from "react-router-dom";
import { useProblem } from "../../hooks/useProblem";


export function ViewProblem() {
    const { id } = useParams<{ id: string }>();
    const { data: problem, isLoading, isError, error } = useProblem(id || '');

    if (isLoading) {
        return <div className="text-center p-4">Loading problem...</div>;
    }

    if (isError) {
        return (
            <div className="text-center p-4 text-red-600">
                Error loading problem: {error.message}
            </div>
        );
    }

    if (!problem) {
        return <div className="text-center p-4">Problem not found.</div>; /* TODO: I dont think this is hit since we throw error*/
    }

    return (
        <DataList.Root>
            <DataList.Item>
                <DataList.Label>ID</DataList.Label>
                <DataList.Value>{problem.Id}</DataList.Value>
            </DataList.Item>
            <DataList.Item>
                <DataList.Label>Question</DataList.Label>
                <DataList.Value>{problem.Question}</DataList.Value>
            </DataList.Item>
            <DataList.Item>
                <DataList.Label>Answer</DataList.Label>
                <DataList.Value>{problem.Answer}</DataList.Value>
            </DataList.Item>
        </DataList.Root>
    );
};