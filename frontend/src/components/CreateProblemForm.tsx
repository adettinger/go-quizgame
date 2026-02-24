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
    })

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
        }
        return true;
    }

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

    const setType = (type: ProblemType) => {
        setFormValues({ ...formValues, Type: type })
    };

    return (
        <Form.Root onSubmit={handleSubmit} onKeyDown={(e) => e.key === 'Enter' && e.preventDefault()}>
            <Flex direction={"column"} gap="3">
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
                        <Button disabled={!areAllFieldsValid()} style={{ alignSelf: 'center' }}>Post a question</Button>
                    </Tooltip>
                </Form.Submit>
            </Flex>
        </Form.Root >
    );
};