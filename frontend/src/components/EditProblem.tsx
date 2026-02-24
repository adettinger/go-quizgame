import { useMutation, useQueryClient } from "@tanstack/react-query";
import { getEnumKeyByValue, MAX_PROBLEM_CHOICES, ProblemType } from "../types/problem";
import { useToast } from "./Toast/ToastContext";
import { Link, useNavigate, useParams } from "react-router-dom";
import { useProblem } from "../hooks/useProblem";
import { useEffect, useState } from "react";
import * as Form from "@radix-ui/react-form";
import { Button, DropdownMenu, Flex, Text, Tooltip } from "@radix-ui/themes";
import type { workingProblem } from "./CreateProblemForm";
import { EditChoicesTable } from "./EditChoices";
import { editProblem } from "../services/problemService";

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

    const areAllFieldsValid = () => {
        // TODO: Only allow certain chars in fields

        // Question and answer are not empty
        if (formValues.Question.trim() === "" && formValues.Answer.trim() === "") {
            return false
        }
        // Type type has no choices
        if (formValues.Type === ProblemType.Choice) {
            if (formValues.Choices.length < 2 || formValues.Choices.length > MAX_PROBLEM_CHOICES) {
                return false
            }
            // Choice cannot be empty string
            if (formValues.Choices.some(choice => choice === "")) {
                return false
            }
            // Answer must be one of the choices
            if (!formValues.Choices.some(choice => choice.toLowerCase() === formValues.Answer.toLowerCase())) {
                return false
            }
            // No duplicate choices
            if (new Set(formValues.Choices).size !== formValues.Choices.length) {
                return false
            }
        }

        return true;
    }

    if (isError) {
        return (
            <Text className="text-center p-4 text-red-600">
                Error loading problem: {error.message}
            </Text>
        );
    }

    if (isLoading) {
        return <Text className="text-center p-4">Loading problem...</Text>;
    }

    const setType = (type: ProblemType) => {
        setFormValues({ ...formValues, Type: type })
    };

    return (
        <Form.Root onSubmit={handleSubmit}>
            <Flex direction="column" gap="3">
                <Form.Field name="Type">
                    <Form.Label>Type: </Form.Label>
                    <DropdownMenu.Root>
                        <DropdownMenu.Trigger>
                            <Button color='gray' variant='soft'>{getEnumKeyByValue(ProblemType, formValues.Type)}<DropdownMenu.TriggerIcon /></Button>
                        </DropdownMenu.Trigger>
                        <DropdownMenu.Content color="gray" variant='soft'>
                            {Object.keys(ProblemType).map((type, index) => (
                                <DropdownMenu.Item key={index} onClick={() => setType(ProblemType[type as keyof typeof ProblemType])}>{type}</DropdownMenu.Item>
                            ))}
                        </DropdownMenu.Content>
                    </DropdownMenu.Root>

                </Form.Field>

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

                {formValues.Type === ProblemType.Choice &&
                    <EditChoicesTable
                        choices={formValues.Choices}
                        setChoices={(choices: string[]) => setFormValues({ ...formValues, Choices: choices })}
                    />
                }

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

                <Form.Submit asChild >
                    <Tooltip content={
                        formValues.Type === ProblemType.Text
                            ? "Fill in all fields"
                            :
                            <Flex direction="column">
                                <Text>Must have at least 2 choices</Text>
                                <Text>Answer must be one of the choices</Text>
                                <Text>{`Can have at most ${MAX_PROBLEM_CHOICES} choices`}</Text>
                            </Flex>
                    }>
                        <Button disabled={!areAllFieldsValid() || mutation.isPending} style={{ alignSelf: 'center' }}>Edit question</Button>
                    </Tooltip>
                </Form.Submit>
            </Flex>
        </Form.Root>
    );
}