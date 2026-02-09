import { Button, Flex, Text } from "@radix-ui/themes";
import * as Form from "@radix-ui/react-form";
import { useEffect, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { submitQuiz } from "../../services/problemService";
import type { SubmitQuizResponse } from "../../types/responses";
import { useToast } from "../Toast/ToastContext";
import type { Question } from "../../types/question";
import { useStartQuiz } from "../../hooks/useStartQuiz";
import { CountdownTimer } from "./timer";

interface quizItems {
    Id: string;
    Question: string;
    Guess: string;
    Answer: string;
    Correct: boolean | undefined;
}

export function Game() {
    const { data, isLoading, isError, error } = useStartQuiz();
    const { showToast } = useToast();
    const [sessionID, setSessionID] = useState('');
    const [timeout, setTimeout] = useState<Date>();
    const [score, setScore] = useState(-1);
    const [quizItems, setQuizItems] = useState<quizItems[]>([]);
    // const [answersMap, setAnswersMap] = useState<Map<string, string>>();
    // const [quizResponse, setQuizResponse] = useState<SubmitQuizResponse>();

    useEffect(() => {
        if (data) {
            setSessionID(data?.SessionId);
            setTimeout(data?.Timeout);
            let initialQuizItems: quizItems[] = [];
            data.Questions.forEach((question: Question) => {
                initialQuizItems.push({
                    Id: question.Id,
                    Question: question.Question,
                    Guess: '',
                    Answer: '',
                    Correct: undefined,
                })
            })
            setQuizItems(initialQuizItems);
        }
    }, [data]);

    const updateGuess = (id: string, answer: string) => {
        setQuizItems((quizItems) =>
            quizItems.map(item =>
                item.Id === id
                    ? { ...item, Guess: answer }
                    : item
            )
        );
    };

    const updateQuizItemsFromResponse = (response: SubmitQuizResponse) => {
        setScore(response.Score);
        const updatedItems = quizItems.map(item => {
            // Find matching item from response
            const matchingResponse = response.Answers.find((element) =>
                element.Id === item.Id  // Use === for comparison and RETURN the result
            );

            // Update Answer and correct
            if (matchingResponse) {
                return {
                    ...item,
                    Answer: matchingResponse.Answer,
                    Correct: matchingResponse.Correct
                };
            } else {
                console.log("Failed to find matching response for id:", item.Id);
                return item;
            }
        });
        setQuizItems(updatedItems);

    };

    const mutation = useMutation({
        mutationFn: submitQuiz,

        onSuccess: (response: SubmitQuizResponse) => {
            // Set scores and answers
            showToast('success', "Success", "Quiz submitted");
            // Update quizItems and score
            updateQuizItemsFromResponse(response);
        },

        onError: () => {
            console.log("Request to submit quiz failed")
            showToast('error', "Error", "Failed to submit quiz");
        },
    })

    const createQuizSubmission = () => {
        return {
            SessionId: sessionID,
            Questions: quizItems ? Array.from(quizItems, (quizItem) => ({ Id: quizItem.Id, Answer: quizItem.Guess })) : [],
        };
    }

    const handleSubmit = async (event) => {
        event.preventDefault();
        console.log('submitting quiz with answers', quizItems);

        mutation.mutate(
            createQuizSubmission()
        );
    };

    const handleExpire = async () => {
        console.log("Time's up!");

        mutation.mutate(
            createQuizSubmission()
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
        <>
            {timeout !== undefined && score < 0 &&
                <CountdownTimer
                    deadline={timeout !== undefined ? timeout : ''}
                    onExpire={handleExpire}
                />
            }
            {score >= 0 &&
                <Text>Score: {score}</Text>
            }
            <Form.Root onSubmit={handleSubmit}>
                <Flex direction={"column"} gap="3" width={"50%"} justify="center">

                    {quizItems?.map((problem) => (
                        <Form.Field name={`Question-${problem.Id}`}>
                            <Flex direction="column" gap="3">
                                <Form.Label>{problem.Question}</Form.Label>
                                <Form.Control asChild>
                                    <input
                                        type="text"
                                        // TODO: Find better style when input is disabled
                                        style={{
                                            color: !!problem.Correct ? (problem.Correct ? "green" : "red") : ""
                                        }}
                                        required
                                        value={problem.Guess}
                                        onChange={(event) => updateGuess(problem.Id, event.target.value)}
                                        disabled={score >= 0}
                                    />
                                </Form.Control>
                                {problem.Correct === false &&
                                    <Text color="red">{problem.Answer}</Text>
                                }
                            </Flex>
                        </Form.Field>
                    ))}
                    <Form.Submit asChild>
                        <Button disabled={score >= 0}>Submit Quiz</Button>
                    </Form.Submit>
                </Flex>
            </Form.Root>
        </>
    );
}