import { Button } from "@radix-ui/themes";
import * as Form from "@radix-ui/react-form";
import { useState } from "react";
import { Toast } from "radix-ui"
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createProblem } from "../services/problemService";
import './ToastStyles.scss';

export function CreateProblemForm() {
    const queryClient = useQueryClient();
    // const [toastOpen, setToastOpen] = useState(false);
    const [toastOpen, setToastOpen] = useState(true);
    const [toastType, setToastType] = useState<'success' | 'error'>('success');
    const [toastMessage, setToastMessage] = useState("");
    const [formValues, setFormValues] = useState({
        Question: "",
        Answer: "",
    });

    const mutation = useMutation({
        mutationFn: createProblem,

        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['problems'] })
            // Display success msg
            // Reset inputs
            setFormValues({ Question: "", Answer: "" })
            setToastType('success');
            setToastMessage("Created problem successfully");
            setToastOpen(true);
        },

        onError: () => {
            // Display error
            console.log("Request to create new problem failed")
            setToastType('error');
            setToastMessage("Failed to create problem");
            setToastOpen(true);
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
        <Toast.Provider>
            <Toast.Root
                open={toastOpen}
                onOpenChange={setToastOpen}
                duration={30000}
            // TODO: Stype success vs failure
            >
                <Toast.Title>{toastType === 'success' ? 'Success' : 'Error'}</Toast.Title>
                <Toast.Description>{toastMessage}</Toast.Description>
            </Toast.Root>
            <Toast.Viewport className="ToastViewport" />
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
        </Toast.Provider>
    );
};