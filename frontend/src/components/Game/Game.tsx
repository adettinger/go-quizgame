import { Button, Flex } from "@radix-ui/themes";
import * as Form from "@radix-ui/react-form";
import { useEffect, useState } from "react";
import { useQuestions } from "../../hooks/useQuestions";
import { useMutation } from "@tanstack/react-query";
import { submitQuiz } from "../../services/problemService";
import type { SubmitQuizResponse } from "../../types/responses";
import { useToast } from "../Toast/ToastContext";


// Note: keeping a map in state is not memory effecient since the map is copied on every update
export function Game() {
    const { data, isLoading, isError, error } = useQuestions();
    const { showToast } = useToast();
    const [answersMap, setAnswersMap] = useState<Map<string, string>>();

    useEffect(() => {
        if (data) {
            const initialAnswers = data.reduce((map, problem) => {
                map.set(problem.Id, '')
                return map;
            }, new Map<string, string>());
            setAnswersMap(initialAnswers);
        }
    }, [data]);

    const updateAnswer = (id: string, answer: string) => {
        const updatedAnswers = new Map(answersMap);
        updatedAnswers.set(id, answer);
        setAnswersMap(updatedAnswers);
    };

    const mutation = useMutation({
        mutationFn: submitQuiz,

        onSuccess: (reponse: SubmitQuizResponse) => {
            // Set scores and answers
            // Show toast
            showToast('success', "Success", "Quiz submitted");
        },

        onError: () => {
            console.log("Request to submit quiz failed")
            showToast('error', "Error", "Failed to submit quiz");
        },
    })

    const handleSubmit = async (event) => {
        event.preventDefault();
        console.log('submitting quiz with answers', answersMap);

        mutation.mutate(
            answersMap ? Array.from(answersMap, ([id, answer]) => ({ Id: id, Answer: answer })) : []
        );
    };

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

    return (
        <Form.Root onSubmit={handleSubmit}>
            <Flex direction={"column"} gap="3" width={"50%"} justify="center">

                {data?.map((problem) => (
                    <Form.Field name={`Question-${problem.Id}`}>
                        <Flex direction="column" gap="3">
                            <Form.Label>{problem.Question}</Form.Label>
                            <Form.Control asChild>
                                <input
                                    type="text"
                                    required
                                    value={answersMap?.get(problem.Id) || ''}
                                    onChange={(event) => updateAnswer(problem.Id, event.target.value)}
                                />
                            </Form.Control>
                        </Flex>
                    </Form.Field>
                ))}
                <Form.Submit asChild>
                    <Button>Submit Quiz</Button>
                </Form.Submit>
            </Flex>
        </Form.Root>
    );
}