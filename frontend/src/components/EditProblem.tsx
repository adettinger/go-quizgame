import { useMutation, useQueryClient } from "@tanstack/react-query";
import { ProblemType } from "../types/problem";
import { useToast } from "./Toast/ToastContext";
import { Link, useNavigate, useParams } from "react-router-dom";
import { useProblem } from "../hooks/useProblem";
import { useEffect, useState } from "react";
import { Text } from "@radix-ui/themes";
import type { workingProblem } from "./CreateProblemForm";
import { editProblem } from "../services/problemService";
import { ProblemForm } from "./ProblemForm";

export function EditProblem() {
    const navigate = useNavigate();
    const queryClient = useQueryClient();
    const { showToast } = useToast();
    const { id } = useParams<{ id: string }>();
    const { data: initProblem, isLoading, isError, error } = useProblem(id || '');

    const [formValues, setFormValues] = useState<workingProblem>({
        Type: ProblemType.Text,
        Question: "",
        Choices: [],
        Answer: "",
    });

    useEffect(() => {
        if (!isLoading && !!initProblem) {
            setFormValues(initProblem);
        }
    }, [initProblem]);

    const mutation = useMutation({
        mutationFn: editProblem,

        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['problems'] });
            queryClient.invalidateQueries({ queryKey: ['problem', id] });
            showToast('success', "Success", <>Edited problem <Link to={`/problem/${id}`}>{id}</Link> successfully</>);
            navigate(`/problem/${id}`);
            // TODO: Instead of navigating away, refresh this page
        },

        onError: (error) => {
            console.log('Failed to edit problem', error);
            showToast('error', "Error", "Failed to edit problem");
        },
    });

    const handleSubmit = async (event: any) => {
        event.preventDefault();

        if (!id) {
            showToast('error', "Error", "Cannot save. Missing ID");
            return;
        }

        console.log('Submitting new problem with data:', {
            Question: formValues.Question,
            Answer: formValues.Answer,
        });

        mutation.mutate({
            Id: id,
            Type: formValues.Type,
            Question: formValues.Question.trim(),
            Answer: formValues.Answer.trim(),
            // TODO: Clean all the choices
            Choices: formValues.Type === ProblemType.Text ? [] : formValues.Choices,
        })
    }

    if (isError) {
        return (
            <Text color="red" className="text-center p-4 text-red-600">
                Error loading problem: {error.message}
            </Text>
        );
    }

    if (isLoading) {
        return <Text className="text-center p-4">Loading problem...</Text>;
    }

    return (
        <ProblemForm
            formValues={formValues}
            setFormValues={setFormValues}
            onSubmit={handleSubmit}
            submitDisabled={mutation.isPending}
        />
    );
}