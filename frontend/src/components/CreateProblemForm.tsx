import { Button } from "@radix-ui/themes";
import * as Form from "@radix-ui/react-form";
import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createProblem } from "../services/problemService";
import './Toast/ToastStyles.scss';
import { useToast } from "./Toast/ToastContext";

export function CreateProblemForm() {
    const queryClient = useQueryClient();
    const { showToast } = useToast();
    const [formValues, setFormValues] = useState({
        Question: "",
        Answer: "",
    });

    const mutation = useMutation({
        mutationFn: createProblem,

        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['problems'] })
            setFormValues({ Question: "", Answer: "" })
            showToast('success', "Success", "Created problem successfully");
        },

        onError: () => {
            console.log("Request to create new problem failed")
            showToast('error', "Error", "Failed to create problem");
        },
    })

    const areAllFieldsValid = () => {
        return formValues.Question.trim() != "" && formValues.Answer.trim() != ""
    }

    const handleSubmit = async (event) => {
        event.preventDefault();

        console.log('Submitting new problem with data:', {
            Question: formValues.Question,
            Answer: formValues.Answer,
        });

        // TODO: Actually send request
        mutation.mutate({
            Question: formValues.Question.trim(),
            Answer: formValues.Answer.trim(),
        })
    }

    return (
        <Form.Root onSubmit={handleSubmit}>
            <Form.Field name="Question">
                <Form.Label>Question: </Form.Label>
                <Form.Message match="valueMissing">
                    Please enter a question
                </Form.Message>
                <Form.Control asChild>
                    <input
                        type="text"
                        required
                        value={formValues.Question}
                        onChange={(event) => { setFormValues({ ...formValues, Question: event.target.value }) }}
                    />
                </Form.Control>
            </Form.Field>

            <Form.Field name="Answer">
                <Form.Label>Answer: </Form.Label>
                <Form.Message match="valueMissing">
                    Please enter a answer
                </Form.Message>
                <Form.Control asChild>
                    <input
                        type="text"
                        required
                        value={formValues.Answer}
                        onChange={(event) => { setFormValues({ ...formValues, Answer: event.target.value }) }}
                    />
                </Form.Control>
            </Form.Field>
            <Form.Submit asChild>
                <Button disabled={!areAllFieldsValid()}>Post a question</Button>
            </Form.Submit>
        </Form.Root>
    );
};