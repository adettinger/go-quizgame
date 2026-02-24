import { Button, DataList, DropdownMenu, Flex } from "@radix-ui/themes";
import { useNavigate, useParams } from "react-router-dom";
import { useProblem } from "../../hooks/useProblem";
import { Pencil1Icon, TrashIcon } from "@radix-ui/react-icons";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { deleteProblemById } from "../../services/problemService";
import { useToast } from "../Toast/ToastContext";
import { ProblemType } from "../../types/problem";


export function ViewProblem() {
    const navigate = useNavigate();
    const queryClient = useQueryClient();
    const { showToast } = useToast();
    const { id } = useParams<{ id: string }>();
    const { data: problem, isLoading, isError, error } = useProblem(id || '');

    const deleteMutation = useMutation({
        mutationFn: deleteProblemById,
        onSuccess: (id: string) => {
            queryClient.invalidateQueries({ queryKey: ['problems'] });
            queryClient.invalidateQueries({ queryKey: ['problem', id] });
            showToast('success', "Success", `Deleted problem ${id}`);
            console.log(`Deleted problem ${id}`);
            navigate("/problems")
        },
        onError: (error) => {
            console.log('Failed to delete problem', error)
            showToast('error', "Error", "Failed to delete problem");
        },
    })

    const handleDelete = (id: string) => {
        if (confirm('Are you sure you want to delete this problem?')) {
            deleteMutation.mutate(id);
        }
    }

    const handleEdit = (id: string) => {
        navigate(`/problem/edit/${id}`)
    }

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
        return <div className="text-center p-4">Problem not found.</div>;
    }

    return (
        <Flex direction={"column"} gap={"3"}>
            <DataList.Root>
                <DataList.Item>
                    <DataList.Label>ID</DataList.Label>
                    <DataList.Value>{problem.Id}</DataList.Value>
                </DataList.Item>
                <DataList.Item>
                    <DataList.Label>Type</DataList.Label>
                    <DataList.Value>{problem.Type}</DataList.Value>
                </DataList.Item>
                <DataList.Item>
                    <DataList.Label>Question</DataList.Label>
                    <DataList.Value>{problem.Question}</DataList.Value>
                </DataList.Item>

                {problem.Type === ProblemType.Choice &&
                    <DataList.Item>
                        <DataList.Label>Choices</DataList.Label>
                        <DataList.Value>
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
                        </DataList.Value>
                    </DataList.Item>
                }
                <DataList.Item>
                    <DataList.Label>Answer</DataList.Label>
                    <DataList.Value>{problem.Answer}</DataList.Value>
                </DataList.Item>
            </DataList.Root>
            <Button
                color="indigo"
                onClick={() => { handleEdit(problem.Id) }}
            >
                Edit
                <Pencil1Icon />
            </Button>
            <Button
                color="red"
                onClick={() => { handleDelete(problem.Id) }}
            >
                Delete
                <TrashIcon />
            </Button>
        </Flex>
    );
};