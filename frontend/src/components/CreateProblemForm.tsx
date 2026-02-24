import { Button, DropdownMenu, Flex, Text, Tooltip } from "@radix-ui/themes";
import * as Form from "@radix-ui/react-form";
import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createProblem } from "../services/problemService";
import './Toast/ToastStyles.scss';
import { useToast } from "./Toast/ToastContext";
import { Link } from "react-router-dom";
import { getEnumKeyByValue, MAX_PROBLEM_CHOICES, ProblemType, type Problem } from "../types/problem";
import { EditChoicesTable } from "./EditChoices";
import { ProblemForm } from "./ProblemForm";

export interface workingProblem {
    Type: ProblemType;
    Question: string;
    Choices: string[];
    Answer: string;
}

export function CreateProblemForm() {
    const queryClient = useQueryClient();
    const { showToast } = useToast();
    const [formValues, setFormValues] = useState<workingProblem>({
        Type: ProblemType.Text,
        Question: "",
        Choices: [],
        Answer: "",
    });

    const mutation = useMutation({
        mutationFn: createProblem,

        onSuccess: (problem: Problem) => {
            queryClient.invalidateQueries({ queryKey: ['problems'] })
            setFormValues({ Type: ProblemType.Text, Question: "", Choices: [], Answer: "" })
            showToast('success', "Success", <>Created problem <Link to={`/problem/${problem.Id}`}>{problem.Id}</Link> successfully</>);
        },

        onError: () => {
            console.log("Request to create new problem failed")
            showToast('error', "Error", "Failed to create problem");
        },
    });

    const handleSubmit = async (event: any) => {
        event.preventDefault();

        console.log('Submitting new problem with data:', {
            Question: formValues.Question,
            Answer: formValues.Answer,
        });

        mutation.mutate({
            Type: formValues.Type,
            Question: formValues.Question.trim(),
            Answer: formValues.Answer.trim(),
            // TODO: Clean all the choices
            Choices: formValues.Type === ProblemType.Text ? [] : formValues.Choices,
        })
    }

    return (
        <ProblemForm
            formValues={formValues}
            setFormValues={setFormValues}
            onSubmit={handleSubmit}
            submitDisabled={mutation.isPending}
        />
    );
};