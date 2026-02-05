import { Button, Flex } from "@radix-ui/themes";
import { useProblems } from "../../hooks/useProblems";
import * as Form from "@radix-ui/react-form";
import { useEffect, useState } from "react";

export function Game() {
    const { data, isLoading, isError, error } = useProblems();
    const [answersMap, setAnswersMap] = useState<Record<string, string>>({});

    useEffect(() => {
        if (data) {
            const initialAnswers = data.reduce((map, problem) => {
                map[problem.Id] = '';
                return map;
            }, {} as Record<string, string>);
            setAnswersMap(initialAnswers);
        }
    }, [data]);

    const updateAnswer = (id: string, answer: string) => {
        setAnswersMap(prev => ({
            ...prev,
            [id]: answer
        }));
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        console.log('submitting quiz with answers', answersMap);
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
                                    value={answersMap[problem.Id] || ''}
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